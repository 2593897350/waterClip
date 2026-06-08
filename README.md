# WaterClip

Services:
- `web/`: Next.js frontend
- `api/`: Go orchestration API
- `processor/`: Python image-processing service

Local dev:
1. Create a Python virtualenv: `python3 -m venv .venv`
2. Install Python deps: `./.venv/bin/pip install fastapi httpx pytest uvicorn`
3. Install frontend deps: `pnpm install --dir web`
4. Start the processor: `cd processor && ../.venv/bin/uvicorn app.main:app --reload --port 8000`
5. Start the API: `cd api && GOCACHE=/Users/liuyafeng/workspace/waterClip/.cache/go-build go run ./cmd/server`
6. Start the web app: `cd web && pnpm dev`

Tests:
- Web: `cd web && pnpm test`
- API: `cd api && GOCACHE=/Users/liuyafeng/workspace/waterClip/.cache/go-build go test ./...`
- Processor: `cd processor && ../.venv/bin/pytest`

Cloud deploy:
1. Copy `.env.production.example` to `.env.production`
2. Adjust `WEB_PORT` if you do not want to expose port `3000`
3. Run `bash scripts/deploy.sh deploy`
4. Check status with `bash scripts/deploy.sh status`
5. Stream logs with `bash scripts/deploy.sh logs web`

Production notes:
- The browser only talks to the `web` container.
- `web` proxies `/api/*` to the internal Go API via `INTERNAL_API_PROXY_TARGET`.
- `api` talks to `processor` through `PROCESSOR_BASE_URL`.
- For a public server, you only need to expose the web port.
