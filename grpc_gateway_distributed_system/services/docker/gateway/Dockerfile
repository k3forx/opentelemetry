FROM golang:1.24.1-bullseye

# ヘルスチェック用にcurlをインストール
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# airをインストール（ホットリロード用）
RUN go install github.com/air-verse/air@latest

RUN mkdir /app
WORKDIR /app

# Copy the entire project
COPY . .

RUN go mod download && \
    go clean -modcache && \
    rm -rf /tmp/* /var/tmp/* /root/.cache

WORKDIR /app/services/gateway

# Set Go environment variables to prevent tmp directory creation
RUN mkdir -p /tmp/go-cache /tmp/go-tmp
ENV GOCACHE=/tmp/go-cache
ENV GOTMPDIR=/tmp/go-tmp

# airのtmpディレクトリを作成
RUN mkdir -p tmp

CMD ["air", "-c", ".air.toml"] 
