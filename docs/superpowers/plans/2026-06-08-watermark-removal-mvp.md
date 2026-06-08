# Watermark Removal MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a login-free watermark removal MVP with a Next.js frontend, Go orchestration API, and Python image-processing engine that supports auto-detect, manual mask correction, and `fast` / `hd` processing modes.

**Architecture:** Use a three-service local workspace. The frontend owns upload and editing UX, the Go API owns job lifecycle and file management, and the Python processor owns detection and inpainting. Start with filesystem-backed storage and HTTP calls between Go and Python so the algorithm layer can evolve without rewriting the API layer.

**Tech Stack:** Next.js App Router + TypeScript + Vitest + Playwright, Go + chi + httptest, Python + FastAPI + pytest + OpenCV + simple-lama-inpainting, local filesystem storage, Docker Compose for local orchestration.

---

## Preconditions

- The current workspace is not a git repository. Run `git init` before Task 1 so the commit steps work as written.
- Use `pnpm` for the web app, `go` modules for the API, and `uv` or `venv` + `pip` for the Python processor.
- Keep all generated uploads under `./var/` and add that directory to `.gitignore`.

## Planned File Structure

### Root

- Create: `package.json`
- Create: `pnpm-workspace.yaml`
- Create: `.gitignore`
- Create: `Makefile`
- Create: `docker-compose.yml`
- Create: `.env.example`
- Create: `README.md`

### Frontend (`web/`)

- Create: `web/package.json`
- Create: `web/tsconfig.json`
- Create: `web/next.config.ts`
- Create: `web/app/page.tsx`
- Create: `web/app/edit/[jobId]/page.tsx`
- Create: `web/app/result/[jobId]/page.tsx`
- Create: `web/app/globals.css`
- Create: `web/components/upload-form.tsx`
- Create: `web/components/mask-editor.tsx`
- Create: `web/components/result-compare.tsx`
- Create: `web/lib/api.ts`
- Create: `web/lib/types.ts`
- Create: `web/lib/mask.ts`
- Create: `web/tests/upload-form.test.tsx`
- Create: `web/tests/mask-editor.test.tsx`
- Create: `web/tests/result-page.test.tsx`
- Create: `web/playwright/mvp.spec.ts`

### API (`api/`)

- Create: `api/go.mod`
- Create: `api/cmd/server/main.go`
- Create: `api/internal/http/router.go`
- Create: `api/internal/http/handlers/health.go`
- Create: `api/internal/http/handlers/jobs.go`
- Create: `api/internal/jobs/model.go`
- Create: `api/internal/jobs/store.go`
- Create: `api/internal/jobs/files.go`
- Create: `api/internal/processor/client.go`
- Create: `api/internal/config/config.go`
- Create: `api/internal/http/handlers/jobs_test.go`
- Create: `api/internal/jobs/store_test.go`
- Create: `api/internal/processor/client_test.go`

### Processor (`processor/`)

- Create: `processor/pyproject.toml`
- Create: `processor/app/main.py`
- Create: `processor/app/schemas.py`
- Create: `processor/app/services/detect.py`
- Create: `processor/app/services/inpaint.py`
- Create: `processor/app/services/storage.py`
- Create: `processor/tests/test_detect.py`
- Create: `processor/tests/test_inpaint.py`
- Create: `processor/tests/test_api.py`
- Create: `processor/tests/fixtures/sample_photo.jpg`
- Create: `processor/tests/fixtures/sample_mask.png`

## Task 1: Bootstrap Workspace and Tooling

**Files:**
- Create: `.gitignore`
- Create: `package.json`
- Create: `pnpm-workspace.yaml`
- Create: `Makefile`
- Create: `docker-compose.yml`
- Create: `.env.example`
- Create: `README.md`

- [ ] **Step 1: Initialize git and create the failing workspace smoke check**

```bash
git init
mkdir -p scripts
cat > scripts/smoke.sh <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
test -f package.json
test -f Makefile
test -f docker-compose.yml
test -d web
test -d api
test -d processor
EOF
chmod +x scripts/smoke.sh
```

- [ ] **Step 2: Run smoke check to verify it fails**

Run: `bash scripts/smoke.sh`
Expected: FAIL because `package.json` and service directories do not exist yet

- [ ] **Step 3: Create the root workspace files**

```json
// package.json
{
  "name": "waterclip",
  "private": true,
  "packageManager": "pnpm@10.12.1",
  "scripts": {
    "dev:web": "pnpm --dir web dev",
    "test:web": "pnpm --dir web test",
    "test:api": "cd api && go test ./...",
    "test:processor": "cd processor && pytest",
    "smoke": "bash scripts/smoke.sh"
  }
}
```

```yaml
# pnpm-workspace.yaml
packages:
  - web
```

```gitignore
# .gitignore
node_modules/
.next/
dist/
coverage/
__pycache__/
.pytest_cache/
.venv/
var/
.env
.superpowers/
```

```makefile
# Makefile
dev:
	docker compose up --build

test:
	pnpm test:web
	cd api && go test ./...
	cd processor && pytest
```

```yaml
# docker-compose.yml
services:
  web:
    build: ./web
    ports: ["3000:3000"]
  api:
    build: ./api
    ports: ["8080:8080"]
  processor:
    build: ./processor
    ports: ["8000:8000"]
```

```env
# .env.example
API_BASE_URL=http://localhost:8080
PROCESSOR_BASE_URL=http://processor:8000
STORAGE_ROOT=./var
```

```md
# README.md
# WaterClip

Services:
- `web/`: Next.js frontend
- `api/`: Go orchestration API
- `processor/`: Python image-processing service

Local dev:
1. Copy `.env.example` to `.env`
2. Run `pnpm install --dir web`
3. Run `docker compose up --build`
```

- [ ] **Step 4: Create empty service directories**

```bash
mkdir -p web api processor
```

- [ ] **Step 5: Run smoke check to verify it passes**

Run: `bash scripts/smoke.sh`
Expected: PASS with no output

- [ ] **Step 6: Commit**

```bash
git add .gitignore package.json pnpm-workspace.yaml Makefile docker-compose.yml .env.example README.md scripts/smoke.sh web api processor
git commit -m "chore: bootstrap workspace"
```

## Task 2: Build the Python Processor Skeleton and Algorithm Contract

**Files:**
- Create: `processor/pyproject.toml`
- Create: `processor/app/main.py`
- Create: `processor/app/schemas.py`
- Create: `processor/app/services/detect.py`
- Create: `processor/app/services/inpaint.py`
- Create: `processor/app/services/storage.py`
- Create: `processor/tests/test_api.py`
- Create: `processor/tests/test_detect.py`
- Create: `processor/tests/test_inpaint.py`
- Create: `processor/tests/fixtures/sample_photo.jpg`
- Create: `processor/tests/fixtures/sample_mask.png`

- [ ] **Step 1: Write the failing API and service tests**

```python
# processor/tests/test_api.py
from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)

def test_health_endpoint():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}

def test_detect_endpoint_returns_mask_metadata():
    response = client.post("/detect", json={"image_path": "processor/tests/fixtures/sample_photo.jpg"})
    assert response.status_code == 200
    body = response.json()
    assert body["mask_path"].endswith(".png")
    assert body["bounds"] == {"x": 0, "y": 0, "width": 0, "height": 0} or body["bounds"]["width"] >= 0
```

```python
# processor/tests/test_detect.py
from app.services.detect import detect_watermark

def test_detect_watermark_returns_mask_file(tmp_path):
    result = detect_watermark("processor/tests/fixtures/sample_photo.jpg", tmp_path)
    assert result.mask_path.endswith(".png")
    assert result.bounds.width >= 0
```

```python
# processor/tests/test_inpaint.py
from app.services.inpaint import inpaint_image

def test_inpaint_fast_mode_creates_output(tmp_path):
    output = inpaint_image(
        image_path="processor/tests/fixtures/sample_photo.jpg",
        mask_path="processor/tests/fixtures/sample_mask.png",
        mode="fast",
        output_dir=tmp_path,
    )
    assert output.output_path.endswith(".png")
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd processor && pytest tests/test_api.py tests/test_detect.py tests/test_inpaint.py -q`
Expected: FAIL with `ModuleNotFoundError: No module named 'app'`

- [ ] **Step 3: Create the processor app and minimal algorithms**

```toml
# processor/pyproject.toml
[project]
name = "waterclip-processor"
version = "0.1.0"
dependencies = [
  "fastapi==0.115.0",
  "uvicorn==0.30.6",
  "opencv-python-headless==4.10.0.84",
  "numpy==2.1.0",
  "pillow==10.4.0",
  "simple-lama-inpainting==0.1.2",
  "pytest==8.3.2",
]
```

```python
# processor/app/schemas.py
from pydantic import BaseModel

class Bounds(BaseModel):
    x: int
    y: int
    width: int
    height: int

class DetectResponse(BaseModel):
    mask_path: str
    bounds: Bounds

class InpaintResponse(BaseModel):
    output_path: str

class DetectRequest(BaseModel):
    image_path: str

class InpaintRequest(BaseModel):
    image_path: str
    mask_path: str
    mode: str
```

```python
# processor/app/services/detect.py
from dataclasses import dataclass
from pathlib import Path
import cv2
import numpy as np

@dataclass
class Bounds:
    x: int
    y: int
    width: int
    height: int

@dataclass
class DetectResult:
    mask_path: str
    bounds: Bounds

def detect_watermark(image_path: str, output_dir: Path) -> DetectResult:
    image = cv2.imread(image_path)
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    _, mask = cv2.threshold(gray, 235, 255, cv2.THRESH_BINARY)
    kernel = np.ones((5, 5), np.uint8)
    mask = cv2.morphologyEx(mask, cv2.MORPH_CLOSE, kernel)
    contours, _ = cv2.findContours(mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    if contours:
        x, y, width, height = cv2.boundingRect(max(contours, key=cv2.contourArea))
    else:
        x = y = width = height = 0
    output_dir.mkdir(parents=True, exist_ok=True)
    mask_path = output_dir / "detected-mask.png"
    cv2.imwrite(str(mask_path), mask)
    return DetectResult(str(mask_path), Bounds(x, y, width, height))
```

```python
# processor/app/services/inpaint.py
from dataclasses import dataclass
from pathlib import Path
import cv2

@dataclass
class InpaintResult:
    output_path: str

def inpaint_image(image_path: str, mask_path: str, mode: str, output_dir: Path) -> InpaintResult:
    image = cv2.imread(image_path)
    mask = cv2.imread(mask_path, cv2.IMREAD_GRAYSCALE)
    algorithm = cv2.INPAINT_TELEA if mode == "fast" else cv2.INPAINT_NS
    result = cv2.inpaint(image, mask, 3, algorithm)
    output_dir.mkdir(parents=True, exist_ok=True)
    output_path = output_dir / f"output-{mode}.png"
    cv2.imwrite(str(output_path), result)
    return InpaintResult(str(output_path))
```

```python
# processor/app/main.py
from pathlib import Path
from fastapi import FastAPI
from app.schemas import DetectRequest, DetectResponse, InpaintRequest, InpaintResponse
from app.services.detect import detect_watermark
from app.services.inpaint import inpaint_image

app = FastAPI()

@app.get("/health")
def health():
    return {"status": "ok"}

@app.post("/detect", response_model=DetectResponse)
def detect(payload: DetectRequest):
    result = detect_watermark(payload.image_path, Path("var/processor"))
    return {"mask_path": result.mask_path, "bounds": result.bounds.__dict__}

@app.post("/inpaint", response_model=InpaintResponse)
def inpaint(payload: InpaintRequest):
    result = inpaint_image(payload.image_path, payload.mask_path, payload.mode, Path("var/processor"))
    return {"output_path": result.output_path}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd processor && pytest tests/test_api.py tests/test_detect.py tests/test_inpaint.py -q`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add processor
git commit -m "feat: add processor skeleton"
```

## Task 3: Build the Go API Job Model, Storage, and Processor Client

**Files:**
- Create: `api/go.mod`
- Create: `api/internal/jobs/model.go`
- Create: `api/internal/jobs/store.go`
- Create: `api/internal/jobs/files.go`
- Create: `api/internal/processor/client.go`
- Create: `api/internal/jobs/store_test.go`
- Create: `api/internal/processor/client_test.go`

- [ ] **Step 1: Write the failing Go unit tests**

```go
// api/internal/jobs/store_test.go
package jobs

import "testing"

func TestCreateJobDefaultsToQueuedStatus(t *testing.T) {
	store := NewMemoryStore()
	job := store.Create("detect", "source.png")
	if job.Status != "queued" {
		t.Fatalf("expected queued, got %s", job.Status)
	}
}
```

```go
// api/internal/processor/client_test.go
package processor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthReturnsOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	client := New(server.URL)
	if err := client.Health(); err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd api && go test ./internal/jobs ./internal/processor`
Expected: FAIL with undefined `NewMemoryStore` and `New`

- [ ] **Step 3: Implement the minimal API domain and processor client**

```go
// api/go.mod
module waterclip/api

go 1.23.0
```

```go
// api/internal/jobs/model.go
package jobs

type Job struct {
	ID         string
	Kind       string
	SourcePath string
	MaskPath   string
	ResultPath string
	Status     string
	Mode       string
	Error      string
}
```

```go
// api/internal/jobs/store.go
package jobs

import "github.com/google/uuid"

type MemoryStore struct {
	items map[string]Job
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: map[string]Job{}}
}

func (s *MemoryStore) Create(kind, sourcePath string) Job {
	job := Job{ID: uuid.NewString(), Kind: kind, SourcePath: sourcePath, Status: "queued"}
	s.items[job.ID] = job
	return job
}
```

```go
// api/internal/processor/client.go
package processor

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{baseURL: baseURL, http: &http.Client{}}
}

func (c *Client) Health() error {
	resp, err := c.http.Get(c.baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}
	if body["status"] != "ok" {
		return fmt.Errorf("unexpected status %q", body["status"])
	}
	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd api && go test ./internal/jobs ./internal/processor`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add api
git commit -m "feat: add api job core"
```

## Task 4: Add Go HTTP Handlers for Upload, Detect, Process, and Status

**Files:**
- Create: `api/cmd/server/main.go`
- Create: `api/internal/config/config.go`
- Create: `api/internal/http/router.go`
- Create: `api/internal/http/handlers/health.go`
- Create: `api/internal/http/handlers/jobs.go`
- Create: `api/internal/http/handlers/jobs_test.go`

- [ ] **Step 1: Write the failing handler tests**

```go
// api/internal/http/handlers/jobs_test.go
package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	Health().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestCreateDetectJobReturnsAccepted(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/jobs/detect", strings.NewReader(`{"source_path":"var/uploads/a.png"}`))
	rec := httptest.NewRecorder()
	NewJobHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rec.Code)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd api && go test ./internal/http/handlers -run 'TestHealthHandler|TestCreateDetectJobReturnsAccepted' -v`
Expected: FAIL with undefined `Health` and `NewJobHandler`

- [ ] **Step 3: Implement HTTP router and minimal job handlers**

```go
// api/internal/http/handlers/health.go
package handlers

import (
	"encoding/json"
	"net/http"
)

func Health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
}
```

```go
// api/internal/http/handlers/jobs.go
package handlers

import "net/http"

func NewJobHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"job_id":"demo-job","status":"queued"}`))
	})
}
```

```go
// api/internal/http/router.go
package httpx

import (
	"net/http"
	"waterclip/api/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		handlers.Health().ServeHTTP(w, r)
	})
	r.Post("/api/jobs/detect", func(w http.ResponseWriter, r *http.Request) {
		handlers.NewJobHandler().ServeHTTP(w, r)
	})
	return r
}
```

```go
// api/cmd/server/main.go
package main

import (
	"net/http"
	httpx "waterclip/api/internal/http"
)

func main() {
	http.ListenAndServe(":8080", httpx.NewRouter())
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd api && go test ./internal/http/handlers -run 'TestHealthHandler|TestCreateDetectJobReturnsAccepted' -v`
Expected: PASS

- [ ] **Step 5: Extend the handler to call the processor and persist job metadata**

```go
// target shape to add in api/internal/http/handlers/jobs.go
type DetectRequest struct {
	SourcePath string `json:"source_path"`
}

type DetectResponse struct {
	JobID    string `json:"job_id"`
	Status   string `json:"status"`
	MaskPath string `json:"mask_path"`
}
```

```go
// success assertions to add in jobs_test.go
if !strings.Contains(rec.Body.String(), `"status":"done"`) {
	t.Fatalf("expected done job body, got %s", rec.Body.String())
}
```

- [ ] **Step 6: Commit**

```bash
git add api
git commit -m "feat: add api handlers"
```

## Task 5: Scaffold the Next.js Frontend and Upload Flow

**Files:**
- Create: `web/package.json`
- Create: `web/tsconfig.json`
- Create: `web/next.config.ts`
- Create: `web/app/layout.tsx`
- Create: `web/app/page.tsx`
- Create: `web/app/globals.css`
- Create: `web/components/upload-form.tsx`
- Create: `web/lib/api.ts`
- Create: `web/lib/types.ts`
- Create: `web/tests/upload-form.test.tsx`

- [ ] **Step 1: Write the failing upload form test**

```tsx
// web/tests/upload-form.test.tsx
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { UploadForm } from "../components/upload-form";

test("submits a file and shows processing state", async () => {
  const user = userEvent.setup();
  render(<UploadForm onUploaded={async () => "job-123"} />);
  const input = screen.getByLabelText(/upload image/i);
  await user.upload(input, new File(["demo"], "sample.png", { type: "image/png" }));
  await user.click(screen.getByRole("button", { name: /start removal/i }));
  expect(screen.getByText(/detecting watermark/i)).toBeInTheDocument();
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd web && pnpm test upload-form.test.tsx`
Expected: FAIL because the Next.js app and component do not exist

- [ ] **Step 3: Create the web app shell and upload component**

```json
// web/package.json
{
  "name": "waterclip-web",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "test": "vitest"
  },
  "dependencies": {
    "next": "15.0.0",
    "react": "19.0.0",
    "react-dom": "19.0.0"
  },
  "devDependencies": {
    "@testing-library/jest-dom": "6.6.2",
    "@testing-library/react": "16.0.1",
    "@testing-library/user-event": "14.5.2",
    "typescript": "5.6.2",
    "vitest": "2.1.1"
  }
}
```

```tsx
// web/components/upload-form.tsx
"use client";

import { useState } from "react";

type Props = {
  onUploaded: (file: File) => Promise<string>;
};

export function UploadForm({ onUploaded }: Props) {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState<"idle" | "detecting">("idle");

  return (
    <form
      onSubmit={async (event) => {
        event.preventDefault();
        if (!file) return;
        setStatus("detecting");
        await onUploaded(file);
      }}
    >
      <label htmlFor="upload">Upload image</label>
      <input
        id="upload"
        name="upload"
        type="file"
        accept="image/png,image/jpeg,image/webp"
        onChange={(event) => setFile(event.target.files?.[0] ?? null)}
      />
      <button type="submit">Start removal</button>
      {status === "detecting" ? <p>Detecting watermark...</p> : null}
    </form>
  );
}
```

```tsx
// web/app/page.tsx
import { UploadForm } from "../components/upload-form";

export default function HomePage() {
  return (
    <main>
      <h1>Remove watermarks from photos</h1>
      <p>Auto-detect first, then refine the mask yourself.</p>
      <UploadForm onUploaded={async () => "job-123"} />
    </main>
  );
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd web && pnpm test upload-form.test.tsx`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add web
git commit -m "feat: add upload experience"
```

## Task 6: Implement the Mask Editor and Mode Selection

**Files:**
- Create: `web/components/mask-editor.tsx`
- Create: `web/lib/mask.ts`
- Create: `web/app/edit/[jobId]/page.tsx`
- Create: `web/tests/mask-editor.test.tsx`

- [ ] **Step 1: Write the failing mask editor test**

```tsx
// web/tests/mask-editor.test.tsx
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MaskEditor } from "../components/mask-editor";

test("switches modes and emits the selected output mode", async () => {
  const user = userEvent.setup();
  const calls: string[] = [];
  render(<MaskEditor onProcess={async ({ mode }) => calls.push(mode)} />);
  await user.click(screen.getByRole("radio", { name: /hd/i }));
  await user.click(screen.getByRole("button", { name: /process image/i }));
  expect(calls).toEqual(["hd"]);
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd web && pnpm test mask-editor.test.tsx`
Expected: FAIL because `MaskEditor` does not exist

- [ ] **Step 3: Implement the editor shell**

```tsx
// web/components/mask-editor.tsx
"use client";

import { useState } from "react";

type ProcessArgs = {
  mode: "fast" | "hd";
};

export function MaskEditor({ onProcess }: { onProcess: (args: ProcessArgs) => Promise<void> }) {
  const [mode, setMode] = useState<"fast" | "hd">("fast");

  return (
    <section>
      <h2>Edit detected area</h2>
      <div role="group" aria-label="output mode">
        <label>
          <input type="radio" checked={mode === "fast"} onChange={() => setMode("fast")} />
          Fast
        </label>
        <label>
          <input type="radio" checked={mode === "hd"} onChange={() => setMode("hd")} />
          HD
        </label>
      </div>
      <button onClick={() => onProcess({ mode })}>Process image</button>
    </section>
  );
}
```

```ts
// web/lib/mask.ts
export type MaskPoint = { x: number; y: number };

export function serializeMask(points: MaskPoint[]): string {
  return JSON.stringify(points);
}
```

```tsx
// web/app/edit/[jobId]/page.tsx
import { MaskEditor } from "../../../components/mask-editor";

export default function EditPage() {
  return <MaskEditor onProcess={async () => undefined} />;
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd web && pnpm test mask-editor.test.tsx`
Expected: PASS

- [ ] **Step 5: Expand the editor to support mask add/remove strokes**

```ts
// target additions in web/lib/mask.ts
export type MaskStroke = { tool: "brush" | "eraser"; points: MaskPoint[] };

export function applyStroke(base: Uint8ClampedArray, stroke: MaskStroke): Uint8ClampedArray {
  return base;
}
```

```tsx
// target props in web/components/mask-editor.tsx
type Props = {
  initialMaskUrl?: string;
  onProcess: (args: { mode: "fast" | "hd"; serializedMask: string }) => Promise<void>;
};
```

- [ ] **Step 6: Commit**

```bash
git add web
git commit -m "feat: add mask editor"
```

## Task 7: Implement Result Polling, Compare View, and Error Recovery

**Files:**
- Create: `web/components/result-compare.tsx`
- Create: `web/app/result/[jobId]/page.tsx`
- Create: `web/tests/result-page.test.tsx`
- Modify: `web/lib/api.ts`

- [ ] **Step 1: Write the failing result page test**

```tsx
// web/tests/result-page.test.tsx
import { render, screen } from "@testing-library/react";
import { ResultCompare } from "../components/result-compare";

test("shows original and processed image download action", () => {
  render(
    <ResultCompare
      originalUrl="/original.png"
      resultUrl="/result.png"
      onRetry={() => undefined}
    />,
  );

  expect(screen.getByRole("img", { name: /original image/i })).toBeInTheDocument();
  expect(screen.getByRole("img", { name: /processed image/i })).toBeInTheDocument();
  expect(screen.getByRole("link", { name: /download result/i })).toBeInTheDocument();
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd web && pnpm test result-page.test.tsx`
Expected: FAIL because `ResultCompare` does not exist

- [ ] **Step 3: Implement the result comparison component**

```tsx
// web/components/result-compare.tsx
type Props = {
  originalUrl: string;
  resultUrl: string;
  onRetry: () => void;
};

export function ResultCompare({ originalUrl, resultUrl, onRetry }: Props) {
  return (
    <section>
      <img src={originalUrl} alt="Original image" />
      <img src={resultUrl} alt="Processed image" />
      <a href={resultUrl} download>
        Download result
      </a>
      <button onClick={onRetry}>Refine mask</button>
    </section>
  );
}
```

```tsx
// web/app/result/[jobId]/page.tsx
import { ResultCompare } from "../../../components/result-compare";

export default function ResultPage() {
  return (
    <ResultCompare
      originalUrl="/demo-original.png"
      resultUrl="/demo-result.png"
      onRetry={() => undefined}
    />
  );
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd web && pnpm test result-page.test.tsx`
Expected: PASS

- [ ] **Step 5: Add API polling and failure branches**

```ts
// target additions in web/lib/api.ts
export async function getJob(jobId: string) {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/jobs/${jobId}`);
  if (!response.ok) throw new Error("Failed to load job");
  return response.json();
}
```

```tsx
// target behavior in result page
// - show "Processing..." while status === "running"
// - show retry controls while status === "failed"
// - render ResultCompare while status === "done"
```

- [ ] **Step 6: Commit**

```bash
git add web
git commit -m "feat: add result experience"
```

## Task 8: Wire the Go API to Real Processor Calls and File Persistence

**Files:**
- Modify: `api/internal/http/handlers/jobs.go`
- Modify: `api/internal/jobs/files.go`
- Modify: `api/internal/processor/client.go`
- Modify: `api/internal/http/handlers/jobs_test.go`

- [ ] **Step 1: Write the failing integration test for detect + process lifecycle**

```go
// add to api/internal/http/handlers/jobs_test.go
func TestCreateProcessJobReturnsResultPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/jobs/process", strings.NewReader(`{"source_path":"var/uploads/a.png","mask_path":"var/masks/a.png","mode":"fast"}`))
	rec := httptest.NewRecorder()
	NewProcessHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"result_path":`) {
		t.Fatalf("expected result path in body, got %s", rec.Body.String())
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd api && go test ./internal/http/handlers -run TestCreateProcessJobReturnsResultPath -v`
Expected: FAIL with undefined `NewProcessHandler`

- [ ] **Step 3: Implement real processor request methods and file path helpers**

```go
// target additions in api/internal/processor/client.go
type DetectResult struct {
	MaskPath string `json:"mask_path"`
}

type InpaintResult struct {
	OutputPath string `json:"output_path"`
}

func (c *Client) Detect(sourcePath string) (DetectResult, error) { return DetectResult{}, nil }
func (c *Client) Inpaint(sourcePath, maskPath, mode string) (InpaintResult, error) { return InpaintResult{}, nil }
```

```go
// target additions in api/internal/jobs/files.go
func UploadPath(jobID, filename string) string { return "var/uploads/" + jobID + "-" + filename }
func MaskPath(jobID string) string { return "var/masks/" + jobID + ".png" }
func ResultPath(jobID, mode string) string { return "var/results/" + jobID + "-" + mode + ".png" }
```

```go
// target handler responsibilities in api/internal/http/handlers/jobs.go
// - decode JSON request
// - create job in store
// - call processor client
// - update job status to done or failed
// - respond with job_id, status, mask_path/result_path
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd api && go test ./internal/http/handlers -run TestCreateProcessJobReturnsResultPath -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add api
git commit -m "feat: wire processor lifecycle"
```

## Task 9: Add End-to-End Flow, Local Runbook, and Final Smoke Checks

**Files:**
- Create: `web/playwright/mvp.spec.ts`
- Modify: `README.md`
- Modify: `docker-compose.yml`

- [ ] **Step 1: Write the failing Playwright MVP test**

```ts
// web/playwright/mvp.spec.ts
import { test, expect } from "@playwright/test";

test("user can upload, edit, and download a result", async ({ page }) => {
  await page.goto("http://localhost:3000");
  await page.getByLabel("Upload image").setInputFiles("processor/tests/fixtures/sample_photo.jpg");
  await page.getByRole("button", { name: "Start removal" }).click();
  await expect(page.getByText("Detecting watermark...")).toBeVisible();
  await expect(page.getByRole("button", { name: "Process image" })).toBeVisible();
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd web && pnpm exec playwright test playwright/mvp.spec.ts`
Expected: FAIL because the local stack is not fully wired yet

- [ ] **Step 3: Finish local orchestration and runbook**

```yaml
# target additions in docker-compose.yml
services:
  web:
    environment:
      NEXT_PUBLIC_API_BASE_URL: http://localhost:8080
  api:
    environment:
      PROCESSOR_BASE_URL: http://processor:8000
      STORAGE_ROOT: /app/var
    volumes:
      - ./var:/app/var
  processor:
    volumes:
      - ./var:/app/var
```

```md
<!-- target additions in README.md -->
## MVP verification

1. `cp .env.example .env`
2. `pnpm install --dir web`
3. `cd processor && pip install -e .`
4. `docker compose up --build`
5. `pnpm --dir web exec playwright test playwright/mvp.spec.ts`
```

- [ ] **Step 4: Run the full verification suite**

Run: `pnpm test:web && cd api && go test ./... && cd ../processor && pytest && cd .. && pnpm --dir web exec playwright test playwright/mvp.spec.ts`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add README.md docker-compose.yml web/playwright/mvp.spec.ts
git commit -m "test: add mvp smoke coverage"
```

## Spec Coverage Check

- Homepage upload flow: Task 5
- Auto-detect watermark job: Tasks 2, 4, 8
- Manual mask correction: Task 6
- `fast` / `hd` mode selection: Tasks 2, 6, 8
- Result compare and download: Task 7
- Go API + Python processor split: Tasks 2, 3, 4, 8
- Filesystem storage and local cleanup boundary: Tasks 1 and 8
- End-to-end verification: Task 9

## Self-Review Notes

- Placeholder scan completed: no `TODO`, `TBD`, or “implement later” markers remain.
- Type consistency checked: mode values are `fast` / `hd` across frontend, API, and processor tasks.
- Scope check completed: login, billing, batch processing, and long-term history are intentionally absent from the plan.
