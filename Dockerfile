FROM alpine

WORKDIR /srv
ENV TZ=Asia/Shanghai

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add --no-cache tzdata curl

COPY main /srv
COPY build /srv/build

EXPOSE 8080

CMD ["./main"]
