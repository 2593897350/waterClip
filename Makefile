dev:
	docker compose up --build

test:
	pnpm test:web
	cd api && GOCACHE=/Users/liuyafeng/workspace/waterClip/.cache/go-build go test ./...
	cd processor && ../.venv/bin/pytest
