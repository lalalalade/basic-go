# 基础镜像
FROM ubuntu:22.04

# 把编译后的打包进来这个镜像，放进工作目录 /app
COPY webook /app/webook
WORKDIR /app

# 执行命令
ENTRYPOINT ["/app/webook"]