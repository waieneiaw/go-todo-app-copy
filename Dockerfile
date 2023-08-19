# デプロイ用コンテナに含めるバイナリを作成するコンテナ
FROM golang:1.21.0-bookworm as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# -----------------------------------------------------------------------------

# デプロイ用のコンテナ
FROM debian:bookworm-slim as deploy

RUN apt update

COPY --from=deploy-builder /app/app .

CMD ["./app"]

# -----------------------------------------------------------------------------

# ローカル環境で利用するホットリロード環境
# `error obtaining VCS status: exit status 128`が発生するため`alpine`を利用する
FROM golang:1.21.0-alpine as dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

CMD ["air"]
