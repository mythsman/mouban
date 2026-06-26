# mouban

> 截至 2025年06月，已汇总有效数据（去除了不存在的条目、被封禁条目后）：图书 447w 本、电影 48w 部、音乐 120w 首、游戏 3.8w 部。

本服务作为 [hexo-douban](https://github.com/mythsman/hexo-douban) 项目的后台数据获取服务，用于根据用户的豆瓣ID，获取用户在豆瓣的书、影、音、游中的标注信息，方便用户快速提取。

## 流程

简要的处理流程如下：

1. 用户输入个人豆瓣ID。
2. 访问读书首页获取用户头像、域名等信息。
3. 访问个人 rss 页面获得用户最新更新时间用于去重。
4. 访问用户书、影、音、游的首页获取总数等信息。
5. 滚动访问用户书、影、音、游的评论页获取评论信息、条目概览。
6. 访问条目详情页获取详细信息，并自动发现其他推荐的用户和条目。
7. 每天定时更新书、影、音、游的首页，获取最新条目。

## 部署

可以使用 docker-compose 进行快速部署，环境变量与 `.env.sample` 中的配置对应（全部采用 `MOUBAN_` 前缀 + 大写命名）：

```yaml
    mouban:
      image: mythsman/mouban
      container_name: "mouban"
      restart: always
      expose:
        - "8080"
      environment:
        - MOUBAN_AGENT_ENABLE=true
        - MOUBAN_AGENT_DISCOVER_LEVEL=1
        - MOUBAN_AGENT_ITEM_CONCURRENCY=5
        - MOUBAN_AGENT_ITEM_MAX=10000
        - MOUBAN_CRAWL_ENABLE=true
        - MOUBAN_STORAGE_ENABLE=true
        - MOUBAN_HTTP_TIMEOUT=30000
        - MOUBAN_HTTP_RETRY_MAX=20
        - MOUBAN_HTTP_INTERVAL_USER=5000
        - MOUBAN_HTTP_INTERVAL_ITEM=2000
        - MOUBAN_HTTP_AUTH=11111:ABCDEFG,http://user:pass@ip:port;
        - MOUBAN_USER_RECHECK_INTERVAL=30m
        - MOUBAN_SERVER_MODE=release
        - MOUBAN_DATASOURCE_HOST=host for mysql
        - MOUBAN_DATASOURCE_USERNAME=user name for mysql
        - MOUBAN_DATASOURCE_PASSWORD=passwd for mysql
```

其中最重要的是 `MOUBAN_HTTP_AUTH` 参数，用于配置登陆态的用户信息和走的http代理，格式为 `<dbcl2>,http://<user>:<password>@<proxyIp>:<proxyPort>;`
，可以配置多个。需要注意的是，豆瓣对于未登录的账号有概率会投毒（[例子](https://movie.douban.com/subject/4881682/)），所以这里采用登陆态账号来处理。

dbcl2需要在cookie中查看：
![img.png](docs/images/img.png)

## 接口

以下为服务暴露的主要接口。线上已经部署了一套公共服务，域名使用 `https://mouban.mythsman.com/` 。

### 常用接口

#### 页面入口（后端渲染）

`https://mouban.mythsman.com/`

页面说明：

- 这是项目默认的前端页面入口（SSR），直接访问即可使用。
- 支持输入 Douban UID / domain / 用户名（精确匹配）查询候选用户。
- 选中用户后可查看该用户当前收录的书、影、游、音全部信息（wish/do/collect）。

#### 用户解析（通过ID/domain/name精确匹配）

`https://mouban.mythsman.com/guest/resolve_user?q={id_or_domain_or_name}`

#### 用户录入/更新

`https://mouban.mythsman.com/guest/check_user?id={your_douban_id}`

```json
{
  "result": {
    "id": 1000001,
    "domain": "ahbei",
    "name": "阿北",
    "thumbnail": "https://img1.doubanio.com/icon/u1000001-30.jpg",
    "book_wish": 81,
    "book_do": 61,
    "book_collect": 115,
    "game_wish": 1,
    "game_do": 0,
    "game_collect": 0,
    "movie_wish": 77,
    "movie_do": 17,
    "movie_collect": 218,
    "song_wish": 23,
    "song_do": 21,
    "song_collect": 24,
    "sync_at": 1667232000,
    "check_at": 1679646797,
    "publish_at": 1570409179
  },
  "success": true
}
```

其中：

* publish_at 表示用户最近一次更新的时间戳。
* check_at 表示最近一次**检测**用户是否有更新的时间戳。
* sync_at 表示最近一次**同步**用户信息的时间戳。

#### 查询用户的读书评论

`https://mouban.mythsman.com/guest/user_book?id={your_douban_id}&action=wish`

`https://mouban.mythsman.com/guest/user_book?id={your_douban_id}&action=do`

`https://mouban.mythsman.com/guest/user_book?id={your_douban_id}&action=collect`

#### 查询用户的电影评论

`https://mouban.mythsman.com/guest/user_movie?id={your_douban_id}&action=wish`

`https://mouban.mythsman.com/guest/user_movie?id={your_douban_id}&action=do`

`https://mouban.mythsman.com/guest/user_movie?id={your_douban_id}&action=collect`

#### 查询用户的游戏评论

`https://mouban.mythsman.com/guest/user_game?id={your_douban_id}&action=wish`

`https://mouban.mythsman.com/guest/user_game?id={your_douban_id}&action=do`

`https://mouban.mythsman.com/guest/user_game?id={your_douban_id}&action=collect`

#### 查询用户的音乐评论

`https://mouban.mythsman.com/guest/user_song?id={your_douban_id}&action=wish`

`https://mouban.mythsman.com/guest/user_song?id={your_douban_id}&action=do`

`https://mouban.mythsman.com/guest/user_song?id={your_douban_id}&action=collect`

### 后台接口

#### 强制更新条目

目前条目下载好后，后续不会进行更新，如有更新需要，目前暂时需要手动强制更新一下。

item_type 取: 1-book 2-movie 3-game 4-song

`https://mouban.mythsman.com/admin/refresh_item?type={item_type}&id={item_douban_id}`

#### 强制更新用户

目前用户的评论信息更新下载好后，后续只会进行增量更新。如果对老的条目进行评论修改、删除等操作是不会同步更新的。

如有更新的要求，目前暂时需要手动强制更新一下。（谨慎使用，会对系统造成较大压力）

`https://mouban.mythsman.com/admin/refresh_user?id={douban_uid}`
