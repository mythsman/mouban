# mouban

本服务作为 [hexo-douban](https://github.com/mythsman/hexo-douban) 项目的后台数据获取服务，用于根据用户的豆瓣ID，获取用户在豆瓣的 书 影 音 游 中的标注信息，方便用户快速提取。

[![dockeri.co](https://dockerico.blankenship.io/image/mythsman/mouban)](https://hub.docker.com/r/mythsman/mouban)

## 流程

简要的处理流程如下：

1. 用户输入个人豆瓣ID。
2. 访问读书首页获取用户头像、域名等信息。
3. 访问个人 rss 页面获得用户最新更新时间用于去重。
4. 访问用户书、影、游的首页获取总数等信息。
5. 滚动访问用户书、影、游的评论页获取评论信息、条目概览。
6. 访问条目详情页获取详细信息，并自动发现其他推荐的用户和条目。

## 部署

可以使用 docker-compose 进行快速部署，环境变量与 application.yml.sample 中的配置对应：

```yaml
    mouban:
      image: mythsman/mouban
      container_name: "mouban"
      restart: always
      expose:
        - "8080"
      environment:
        - GIN_MODE=release
        - agent__enable=true
        - agent__flow__discover=false
        - agent__discover__level=1
        - agent__item__concurrency=5
        - agent__item__max=10000
        - http__timeout=30000
        - http__retry_max=20
        - http__interval__user=5000
        - http__interval__item=2000
        - http__auth=11111:ABCDEFG,http://user:pass@ip:port;
        - server__cors=https://yourdomain.com
        - server__limit=30m
        - datasource__host=host for mysql
        - datasource__username=user name for mysql
        - datasource__password=passwd for mysql
```

其中最重要的是 http__auth 参数，用于配置登陆态的用户信息和走的http代理，格式为 `<doubanUid>:<bid>,http://<user>:<password>@<proxyIp>:<proxyPort>;`
，可以配置多个。

bid需要在cookie中查看：

![bid.png](image/img.png)

## 接口

```
# 将 {your_douban_id} 改为你的豆瓣数字ID

# 用户录入/更新

http://localhost:8080/guest/check_user?id={your_douban_id}

# 查询用户的读书评论

http://localhost:8080/guest/user_book?id={your_douban_id}&action=wish

http://localhost:8080/guest/user_book?id={your_douban_id}&action=do

http://localhost:8080/guest/user_book?id={your_douban_id}&action=collect

# 查询用户的电影评论

http://localhost:8080/guest/user_movie?id={your_douban_id}&action=wish

http://localhost:8080/guest/user_movie?id={your_douban_id}&action=do

http://localhost:8080/guest/user_movie?id={your_douban_id}&action=collect

# 查询用户的游戏评论

http://localhost:8080/guest/user_game?id={your_douban_id}&action=wish

http://localhost:8080/guest/user_game?id={your_douban_id}&action=do

http://localhost:8080/guest/user_game?id={your_douban_id}&action=collect

# 查询用户的音乐评论

http://localhost:8080/guest/user_song?id={your_douban_id}&action=wish

http://localhost:8080/guest/user_song?id={your_douban_id}&action=do

http://localhost:8080/guest/user_song?id={your_douban_id}&action=collect
```
