FROM node:20-alpine

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY web/package.json package.json
COPY pnpm-lock.yaml pnpm-lock.yaml
RUN npm install -g pnpm && pnpm install --frozen-lockfile

COPY web .
RUN pnpm build

EXPOSE 3000

CMD ["pnpm", "start", "--hostname", "0.0.0.0", "--port", "3000"]
