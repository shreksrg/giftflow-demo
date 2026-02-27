# Stage 1: Builder
# 建议使用具体版本号而非隐式 latest，保证构建可复现
FROM public.ecr.aws/docker/library/golang:1.26-alpine AS builder

WORKDIR /app

# 安装 git (仅当依赖中包含 git 引用时需要)
RUN apk add --no-cache git ca-certificates

# 【关键优化】先复制依赖管理文件
COPY go.mod go.sum ./

# 【关键优化】下载依赖
# 如果 go.mod 和 go.sum 没有变化，这一步将直接使用缓存，无需重新下载
RUN go mod download

# 复制源代码
# 此时才复制代码，确保代码变动不会导致上面的依赖下载层失效
COPY . .

# 静态编译优化
# CGO_ENABLED=0 确保二进制完全静态链接，可在 alpine 甚至 scratch 中运行
# -s -w 去除符号表和调试信息，减小体积
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /usr/local/bin/giftflow-backend .

# Stage 2: Runner
# 固定 Alpine 版本，避免 "latest" 变动带来的不确定性
FROM public.ecr.aws/docker/library/alpine:3.23

# 安装运行时必要工具
# ca-certificates: 用于 HTTPS 请求
# curl: 用于健康检查 (如果不需要健康检查可移除以进一步减小体积)
RUN apk add --no-cache ca-certificates curl

WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件
COPY --from=builder /usr/local/bin/giftflow-backend /usr/local/bin/giftflow-backend

# 创建非 root 用户运行（安全最佳实践，可选）
# RUN adduser -D -g '' appuser
# USER appuser

EXPOSE 8080

# 直接执行二进制文件
CMD ["/usr/local/bin/giftflow-backend"]
