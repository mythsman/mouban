FROM golang:1.17 as build

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /srv

ADD . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o main main.go

FROM scratch as prod

WORKDIR /srv

COPY --from=build /srv/main /srv
COPY --from=build /srv/application.yml.sample /srv/application.yml
COPY --from=build /srv/cookie.txt.sample /srv/cookie.txt

EXPOSE 8080

# 启动服务
CMD ["./main"]
