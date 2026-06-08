FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download || true

COPY . .
RUN go build -o /usr/local/bin/waterclip-api ./cmd/server

EXPOSE 8080

CMD ["/usr/local/bin/waterclip-api"]
