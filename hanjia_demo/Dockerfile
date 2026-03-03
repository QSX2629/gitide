# 构建Go应用镜像
FROM golang:1.23.0-alpine
WORKDIR /app

# 复制代码+下载依赖+编译
COPY . .

# 设置GOPROXY环境变量，使用国内代理，并指定main目录进行构建
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download && \
    go build -o app ./main

# 暴露应用端口
EXPOSE 8080

# 启动应用
CMD ["./app"]