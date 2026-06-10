FROM node:20-bookworm-slim

WORKDIR /app

ARG NPM_REGISTRY=https://registry.npmjs.org/
ARG NPM_STRICT_SSL=true

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY web/package.json package.json
COPY web/package.docker.json package.json
RUN npm config set registry "$NPM_REGISTRY" \
    && npm config set strict-ssl "$NPM_STRICT_SSL" \
    && npm install

COPY web .
RUN npm run build

EXPOSE 3000

CMD ["npm", "run", "start", "--", "--hostname", "0.0.0.0", "--port", "3000"]
