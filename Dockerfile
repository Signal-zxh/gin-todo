# 构建阶段
FROM golang:1.26.3 AS build
WORKDIR /src

# 先复制依赖文件并下载（利用 Docker 缓存）
COPY go.mod go.sum ./
RUN go mod download

# 再复制源码并编译
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /bin/server ./

# 运行阶段（用最小镜像）
FROM alpine:latest AS final

WORKDIR /app

# 从构建阶段复制编译好的程序
COPY --from=build /bin/server /app/server
COPY templates /app/templates

EXPOSE 8080
ENTRYPOINT ["/app/server"]
