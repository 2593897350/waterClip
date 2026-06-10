FROM node:20-bookworm-slim

WORKDIR /app

ARG NPM_REGISTRY=https://registry.npmjs.org/
ARG NPM_BOOTSTRAP_REGISTRY=https://registry.npmjs.org/
ARG NPM_STRICT_SSL=true

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY web/package.json package.json
COPY pnpm-lock.yaml pnpm-lock.yaml
RUN npm config set registry "$NPM_REGISTRY" \
    && npm config set strict-ssl "$NPM_STRICT_SSL" \
    && npm install -g pnpm --registry "$NPM_BOOTSTRAP_REGISTRY" \
    && pnpm config set registry "$NPM_REGISTRY" \
    && pnpm config set strict-ssl "$NPM_STRICT_SSL" \
    && pnpm install --frozen-lockfile

COPY web .
RUN pnpm build

EXPOSE 3000

CMD ["pnpm", "start", "--hostname", "0.0.0.0", "--port", "3000"]
