import json
import pymysql
import asyncio
import aiohttp
from tqdm import tqdm
from urllib.parse import urlparse, urlunparse
import urllib3
import time
from collections import deque
import hashlib
import mimetypes
import boto3

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


def load_config(path="config.json"):
    with open(path, "r", encoding="utf-8") as f:
        return json.load(f)


def main():
    cfg = load_config()

    host = cfg.get("host", "127.0.0.1")
    port = int(cfg.get("port", 3306))
    user = cfg["username"]
    password = cfg["password"]
    database = cfg["database"]
    force_ip = cfg.get("force_ip", "127.0.0.1")

    s3 = boto3.client(
        "s3",
        endpoint_url=cfg.get("s3_endpoint"),
        region_name=cfg.get("s3_region"),
        aws_access_key_id=cfg.get("s3_access_key"),
        aws_secret_access_key=cfg.get("s3_secret_key"),
    )

    conn = pymysql.connect(
        host=host,
        port=port,
        user=user,
        password=password,
        database=database,
        charset="utf8mb4",
        cursorclass=pymysql.cursors.DictCursor,
    )

    with conn.cursor() as cursor:
        cursor.execute("SELECT COUNT(*) AS cnt FROM storage")
        total = cursor.fetchone()["cnt"]

    batch_size = 2000
    last_id = 0

    concurrency = 500
    q = asyncio.Queue(maxsize=5000)
    semaphore = asyncio.Semaphore(concurrency)

    stats = {"200": 0, "non200": 0, "error": 0}
    timestamps = deque()
    window = 5

    def calc_md5(data):
        h = hashlib.md5()
        h.update(data)
        return h.hexdigest()

    def ext_from_content_type(ct):
        if not ct:
            return "jpg"
        ext = mimetypes.guess_extension(ct.split(";")[0])
        if ext:
            return ext.lstrip(".")
        return "jpg"

    async def repair_404(session, row):
        rid = row["id"]

        db = pymysql.connect(
            host=host,
            port=port,
            user=user,
            password=password,
            database=database,
            charset="utf8mb4",
            cursorclass=pymysql.cursors.DictCursor,
        )

        with db.cursor() as c:
            c.execute("SELECT source FROM storage WHERE id=%s", (rid,))
            r = c.fetchone()

        db.close()

        if not r or not r.get("source"):
            return

        source = r["source"]

        try:
            async with session.get(
                source,
                headers={
                    "Referer": "https://www.douban.com",
    	            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
                },
                
                allow_redirects=True,
            ) as resp:
                if resp.status != 200:
                    return
                data = await resp.read()
                ct = resp.headers.get("content-type")
        except Exception:
            return

        md5 = calc_md5(data)
        ext = ext_from_content_type(ct)
        filename = f"{md5}.{ext}"

        try:
            s3.put_object(
                Bucket="douban",
                Key=filename,
                Body=data,
                ContentType=ct or "image/jpeg",
            )
        except Exception:
            return

        new_url = f"{cfg['s3_endpoint']}/douban/{filename}"

        db = pymysql.connect(
            host=host,
            port=port,
            user=user,
            password=password,
            database=database,
            charset="utf8mb4",
            cursorclass=pymysql.cursors.DictCursor,
        )

        with db.cursor() as c:
            c.execute(
                "UPDATE storage SET target=%s, md5=%s WHERE id=%s",
                (new_url, md5, rid),
            )
            db.commit()

        db.close()

    async def check_row(session, row):
        url = row["target"]

        async with semaphore:
            parsed = urlparse(url)

            forced_url = url
            headers = None

            if force_ip and parsed.hostname:
                new_netloc = force_ip
                if parsed.port:
                    new_netloc += f":{parsed.port}"

                forced_url = urlunparse(parsed._replace(netloc=new_netloc))
                headers = {"Host": parsed.hostname}

            try:
                async with session.head(forced_url, headers=headers, allow_redirects=True) as resp:
                    status = resp.status

                if status >= 400:
                    async with session.get(forced_url, headers=headers, allow_redirects=True) as resp:
                        status = resp.status

            except Exception:
                return row, "error"

        return row, status

    async def producer():
        nonlocal last_id
        with conn.cursor() as cursor:
            while True:
                cursor.execute(
                    "SELECT id,target FROM storage WHERE id > %s ORDER BY id ASC LIMIT %s",
                    (last_id, batch_size),
                )
                rows = cursor.fetchall()

                if not rows:
                    break

                last_id = rows[-1]["id"]

                for r in rows:
                    await q.put(r)

        for _ in range(concurrency):
            await q.put(None)

    async def worker(session, pbar):
        while True:
            row = await q.get()
            if row is None:
                break

            r, status = await check_row(session, row)

            if status == 200:
                stats["200"] += 1
            elif isinstance(status, int):
                stats["non200"] += 1
                if status == 404:
                    await repair_404(session, r)
            else:
                stats["error"] += 1

            pbar.update(1)

            now = time.time()
            timestamps.append(now)

            while timestamps and now - timestamps[0] > window:
                timestamps.popleft()

            recent_qps = len(timestamps) / window if timestamps else 0

            processed = stats["200"] + stats["non200"] + stats["error"]
            remaining = total - processed
            eta = int(remaining / recent_qps) if recent_qps else 0

            pbar.set_postfix(
                {
                    "200": stats["200"],
                    "bad": stats["non200"],
                    "err": stats["error"],
                    "qps": int(recent_qps),
                    "eta": f"{eta}s",
                }
            )

    async def runner():
        timeout = aiohttp.ClientTimeout(total=6)
        connector = aiohttp.TCPConnector(limit=concurrency, ssl=False)

        async with aiohttp.ClientSession(timeout=timeout, connector=connector) as session:
            with tqdm(total=total, desc="scan+repair", mininterval=0.5) as pbar:
                prod = asyncio.create_task(producer())
                workers = [asyncio.create_task(worker(session, pbar)) for _ in range(concurrency)]
                await prod
                await asyncio.gather(*workers)

    asyncio.run(runner())


if __name__ == "__main__":
    main()
