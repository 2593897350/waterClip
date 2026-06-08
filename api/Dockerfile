FROM golang:1.22-bookworm

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod ./
RUN go mod download || true

COPY . .
RUN go build -o /usr/local/bin/waterclip-api ./cmd/server

EXPOSE 8080

CMD ["/usr/local/bin/waterclip-api"]
