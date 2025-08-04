FROM golang:1.24-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN LD_LIBRARY_PATH=/app/pkg/vosk-linux-x86_64-0.3.45 CGO_CPPFLAGS="-I /app/pkg/vosk-linux-x86_64-0.3.45" CGO_LDFLAGS="-L /app/pkg/vosk-linux-x86_64-0.3.45" go build -o /app/copium-bot cmd/copium-bot/main.go

FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    openssl \
    ffmpeg && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY  --from=builder /app/copium-bot /app/copium-bot

COPY pkg /app/pkg

ENV LD_LIBRARY_PATH /app/pkg/vosk-linux-x86_64-0.3.45

ENTRYPOINT ["/app/copium-bot"]