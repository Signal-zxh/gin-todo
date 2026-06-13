# 构建基础阶段
FROM golang:1.26.3 AS build-base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

# 编译阶段（先测试再编译）
FROM build-base AS build
RUN go test -v ./... && CGO_ENABLED=0 go build -o /bin/server ./

# 运行阶段
FROM alpine:latest AS final
WORKDIR /app
COPY --from=build /bin/server /app/server
COPY templates /app/templates
EXPOSE 8080
ENTRYPOINT ["/app/server"]
