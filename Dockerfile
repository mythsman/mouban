FROM alpine

WORKDIR /srv
ENV TZ=Asia/Shanghai

RUN apk add --no-cache tzdata

COPY main /srv
COPY application.yml.sample /srv/application.yml
COPY cookie.txt.sample /srv/cookie.txt

EXPOSE 8080

# 启动服务
CMD ["./main"]
