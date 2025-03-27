# ビルドステージ
FROM golang:latest AS builder

# 作業ディレクトリ
WORKDIR /server-app

# ホストのファイルをコンテナにコピー
COPY server-app/ .

# アプリケーションビルド
RUN go build -o main .

# 実行ステージ
FROM ubuntu:latest AS runner

# ランタイム依存関係
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# 作業ディレクトリを作成
WORKDIR /app

# ビルドステージからビルド済みのバイナリをコピー
COPY --from=builder /server-app/main .

# .env ファイルコピー
COPY server-app/build/app/.env ./build/app/.env

# ポート開放
EXPOSE 8080

# 実行コマンド
CMD ["./main"]