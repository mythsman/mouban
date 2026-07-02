# syntax=docker/dockerfile:1.7

FROM node:24-alpine AS frontend-builder
WORKDIR /src/web

COPY web/package.json web/package-lock.json* ./
RUN --mount=type=cache,target=/root/.npm \
    if [ -f package-lock.json ]; then npm ci --no-audit --no-fund; else npm install --no-audit --no-fund; fi

COPY web/ ./
RUN npm run build


FROM golang:1.26-alpine AS backend-builder
WORKDIR /src

ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/swaggo/swag/cmd/swag@v1.16.6 && \
    /go/bin/swag init --outputTypes json,yaml -g main.go -o docs && \
    mkdir -p build/swagger && \
    cp docs/swagger.json build/swagger/swagger.json && \
    cp tools/swagger/index.html build/swagger/index.html && \
    go build -ldflags="-s -w" -o /out/main .


FROM alpine:3.22

WORKDIR /srv
ENV TZ=Asia/Shanghai

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache tzdata curl

COPY --from=backend-builder /out/main /srv/main
COPY --from=frontend-builder /src/build /srv/build
COPY --from=backend-builder /src/build/swagger /srv/build/swagger

EXPOSE 8080

CMD ["./main"]
