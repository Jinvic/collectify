# Stage 1：构建前端
FROM node:18-alpine AS frontend-builder
WORKDIR /app
# 复制前端依赖文件
COPY web/package.json web/pnpm-lock.yaml ./web/
# 安装pnpm及前端依赖
RUN npm install -g pnpm && \
    cd web && \
    pnpm install --frozen-lockfile
# 复制前端源码并构建
COPY web/ ./web/
RUN cd web && pnpm run build

# Stage 2：构建后端
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app
# 安装Git（Go模块需要）
RUN apk add --no-cache git
# 复制Go模块文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download
# 复制源码并构建
COPY . .
COPY --from=frontend-builder /app/web/build ./web/build
RUN GOOS=linux go build -a -ldflags="-w -s" -o collectify .

# Stage 3：生成最终镜像
FROM alpine:latest
WORKDIR /app/data
# 安装CA证书
RUN apk --no-cache add ca-certificates
# 复制构建好的文件
COPY --from=backend-builder /app/collectify /app/collectify
# 设置文件权限
RUN chmod +x /app/collectify
# 启动程序
EXPOSE 8080
CMD ["/app/collectify"]