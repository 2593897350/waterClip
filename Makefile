dev:
	docker compose up --build

test:
	pnpm test:web
	cd api && go test ./...
	cd processor && pytest
