import json
import pymysql
import asyncio
import aiohttp
import hashlib
import mimetypes
import boto3
import urllib3
import time
from urllib.parse import urlparse, urlunparse
from collections import deque
from tqdm import tqdm

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


def load_config():
    with open("config.json", "r", encoding="utf-8") as f:
        return json.load(f)


def md5(data):
    h = hashlib.md5()
    h.update(data)
    return h.hexdigest()


def ext_from_type(ct):
    if not ct:
        return "jpg"
    ext = mimetypes.guess_extension(ct.split(";")[0])
    if ext:
        return ext.lstrip(".")
    return "jpg"


def main():
    cfg = load_config()

    db = pymysql.connect(
        host=cfg["host"],
        port=int(cfg["port"]),
        user=cfg["username"],
        password=cfg["password"],
        database=cfg["database"],
        charset="utf8mb4",
        cursorclass=pymysql.cursors.DictCursor,
    )

    with db.cursor() as c:
        c.execute("SELECT COUNT(*) AS c FROM storage")
        total = c.fetchone()["c"]

    scan_q = asyncio.Queue(10000)
    repair_q = asyncio.Queue(2000)

    stats = {"scan": 0, "checked": 0, "bad": 0, "fixed": 0, "err": 0}

    timestamps = deque()
    window = 5

    s3 = boto3.client(
        "s3",
        endpoint_url=cfg["s3_endpoint"],
        region_name=cfg["s3_region"],
        aws_access_key_id=cfg["s3_access_key"],
        aws_secret_access_key=cfg["s3_secret_key"],
    )

    async def scan_worker():
        last = 0
        batch = 2000
        with db.cursor() as c:
            while True:
                c.execute(
                    "SELECT id,target FROM storage WHERE id>%s ORDER BY id LIMIT %s",
                    (last, batch),
                )
                rows = c.fetchall()
                if not rows:
                    break
                last = rows[-1]["id"]
                for r in rows:
                    await scan_q.put(r)
                    stats["scan"] += 1

        for _ in range(300):
            await scan_q.put(None)

    async def head_worker(session):
        force_ip = cfg.get("force_ip")
        while True:
            row = await scan_q.get()
            if row is None:
                await repair_q.put(None)
                break

            url = row["target"]
            parsed = urlparse(url)

            forced = url
            headers = None

            if force_ip and parsed.hostname:
                netloc = force_ip
                if parsed.port:
                    netloc += f":{parsed.port}"
                forced = urlunparse(parsed._replace(netloc=netloc))
                headers = {"Host": parsed.hostname}

            try:
                async with session.head(forced, headers=headers, allow_redirects=True) as r:
                    status = r.status

                if status >= 400:
                    async with session.get(forced, headers=headers, allow_redirects=True) as r:
                        status = r.status

            except Exception as e:
                stats["err"] += 1
                print(f"[ERR][HEAD] id={row['id']} url={url} err={e}")
                continue

            stats["checked"] += 1

            if status == 404:
                stats["bad"] += 1
                print(f"[BAD][404] id={row['id']} url={url}")
                await repair_q.put(row)

    async def repair_worker(session):
        while True:
            row = await repair_q.get()
            if row is None:
                break

            rid = row["id"]

            db2 = pymysql.connect(
                host=cfg["host"],
                port=int(cfg["port"]),
                user=cfg["username"],
                password=cfg["password"],
                database=cfg["database"],
                charset="utf8mb4",
                cursorclass=pymysql.cursors.DictCursor,
            )

            with db2.cursor() as c:
                c.execute("SELECT source FROM storage WHERE id=%s", (rid,))
                r = c.fetchone()

            db2.close()

            if not r or not r.get("source"):
                print(f"[ERR][SOURCE] id={rid} source_missing")
                continue

            try:
                async with session.get(
                    r["source"],
                    headers={"Referer": "https://www.douban.com"},
                    allow_redirects=True,
                ) as resp:
                    if resp.status != 200:
                        continue
                    data = await resp.read()
                    ct = resp.headers.get("content-type")
            except Exception:
                print(f"[ERR][DOWNLOAD] id={rid} source={r.get('source')}")
                continue

            name = f"{md5(data)}.{ext_from_type(ct)}"

            try:
                s3.put_object(
                    Bucket="douban",
                    Key=name,
                    Body=data,
                    ContentType=ct or "image/jpeg",
                )
            except Exception:
                print(f"[ERR][UPLOAD] id={rid} file={name}")
                continue

            new_url = f"{cfg['s3_endpoint']}/douban/{name}"

            db3 = pymysql.connect(
                host=cfg["host"],
                port=int(cfg["port"]),
                user=cfg["username"],
                password=cfg["password"],
                database=cfg["database"],
                charset="utf8mb4",
                cursorclass=pymysql.cursors.DictCursor,
            )

            with db3.cursor() as c:
                c.execute(
                    "UPDATE storage SET target=%s, md5=%s WHERE id=%s",
                    (new_url, name.split(".")[0], rid),
                )
                db3.commit()

            db3.close()

            stats["fixed"] += 1
            print(f"[FIXED] id={rid} -> {new_url}")

    async def runner():
        timeout = aiohttp.ClientTimeout(total=6)
        conn = aiohttp.TCPConnector(limit=300, ssl=False)

        async with aiohttp.ClientSession(timeout=timeout, connector=conn) as session:
            with tqdm(total=total) as pbar:

                scan = asyncio.create_task(scan_worker())
                heads = [asyncio.create_task(head_worker(session)) for _ in range(300)]
                repairs = [asyncio.create_task(repair_worker(session)) for _ in range(20)]

                async def progress():
                    last = 0
                    while True:
                        await asyncio.sleep(1)
                        cur = stats["checked"]
                        pbar.update(cur - last)
                        last = cur

                        pbar.set_postfix(
                            scan=stats["scan"],
                            checked=stats["checked"],
                            bad=stats["bad"],
                            fixed=stats["fixed"],
                            err=stats["err"],
                        )

                prog = asyncio.create_task(progress())

                await scan
                await asyncio.gather(*heads)
                await asyncio.gather(*repairs)
                prog.cancel()

    asyncio.run(runner())


if __name__ == "__main__":
    main()
