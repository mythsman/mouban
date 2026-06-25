FROM alpine

WORKDIR /srv
ENV TZ=Asia/Shanghai

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add --no-cache tzdata

COPY main /srv
COPY templates /srv/templates
COPY static /srv/static

EXPOSE 8080

# 启动服务
CMD ["./main"]
