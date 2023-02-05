FROM alpine

WORKDIR /srv
ENV TZ=Asia/Shanghai

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add --no-cache tzdata

COPY main /srv
COPY application.yml.sample /srv/application.yml

EXPOSE 8080

# 启动服务
CMD ["./main"]
