# WaterClip

Services:
- `web/`: Next.js frontend
- `api/`: Go orchestration API
- `processor/`: Python image-processing service

Local dev:
1. Copy `.env.example` to `.env`
2. Run `pnpm install --dir web`
3. Run `docker compose up --build`
