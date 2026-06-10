FROM node:20-bookworm-slim

WORKDIR /app

ARG NPM_REGISTRY=https://registry.npmjs.org/

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY web/package.json package.json
COPY pnpm-lock.yaml pnpm-lock.yaml
RUN npm config set registry "$NPM_REGISTRY" \
    && npm install -g pnpm \
    && pnpm config set registry "$NPM_REGISTRY" \
    && pnpm install --frozen-lockfile

COPY web .
RUN pnpm build

EXPOSE 3000

CMD ["pnpm", "start", "--hostname", "0.0.0.0", "--port", "3000"]
