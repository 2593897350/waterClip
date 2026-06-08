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
