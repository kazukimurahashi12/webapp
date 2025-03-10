FROM golang:latest

# ホストのファイルをコンテナにコピー
COPY server-app/ /server-app

# .env ファイルをコンテナ内にコピー
COPY server-app/build/app/.env /server-app/build/app/.env

# 作業ディレクトリを指定
WORKDIR /server-app

# ポートを開放
EXPOSE 8080

# ビルドコマンドや実行コマンド
RUN go build -o main .

# コンテナが起動した際に実行するコマンドやアプリケーションを指定
CMD ["./main"]