FROM golang:1.22-bookworm

WORKDIR /app

ARG GOPROXY=https://proxy.golang.org,direct
ARG GOSUMDB=sum.golang.org

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN go env -w GOPROXY="$GOPROXY" \
    && go env -w GOSUMDB="$GOSUMDB"

COPY go.mod ./
RUN go mod download || true

COPY . .
RUN go build -o /usr/local/bin/waterclip-api ./cmd/server

EXPOSE 8080

CMD ["/usr/local/bin/waterclip-api"]
