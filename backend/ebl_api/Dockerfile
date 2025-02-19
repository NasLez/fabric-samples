FROM golang:1.23.5-alpine3.21
LABEL authors="wxl"
EXPOSE 8888

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 编译 Go 代码
RUN go build -o main .

# 设置容器启动时运行的命令
CMD ["./main"]