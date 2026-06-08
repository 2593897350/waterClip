import { afterEach, expect, test, vi } from "vitest";

import { createDetectJob, getJob } from "../lib/api";

afterEach(() => {
  vi.unstubAllGlobals();
});

test("createDetectJob posts multipart form data and returns the job id", async () => {
  const fetchMock = vi.fn(async () => ({
    ok: true,
    json: async () => ({ job_id: "job-123", status: "done", mask_path: "var/masks/a.pgm" })
  }));
  vi.stubGlobal("fetch", fetchMock);

  const result = await createDetectJob(new File(["demo"], "sample.ppm", { type: "image/x-portable-pixmap" }));

  expect(result).toEqual({
    jobId: "job-123",
    status: "done",
    maskPath: "var/masks/a.pgm"
  });
  expect(fetchMock).toHaveBeenCalledTimes(1);
});

test("getJob returns the current job payload", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => ({
      ok: true,
      json: async () => ({
        job_id: "job-123",
        status: "done",
        source_path: "var/uploads/a.ppm",
        mask_path: "var/masks/a.pgm",
        result_path: "var/results/a-fast.ppm"
      })
    })),
  );

  const result = await getJob("job-123");

  expect(result).toEqual({
    jobId: "job-123",
    sourcePath: "var/uploads/a.ppm",
    maskPath: "var/masks/a.pgm",
    status: "done",
    resultPath: "var/results/a-fast.ppm"
  });
});

test("processJob posts the selected mode and returns a result path", async () => {
  const fetchMock = vi.fn(async () => ({
    ok: true,
    json: async () => ({ job_id: "job-456", status: "done", result_path: "var/results/b-hd.ppm" })
  }));
  vi.stubGlobal("fetch", fetchMock);

  const result = await import("../lib/api").then(({ processJob }) =>
    processJob("var/uploads/b.ppm", "var/masks/b.pgm", "hd"),
  );

  expect(result).toEqual({
    jobId: "job-456",
    status: "done",
    resultPath: "var/results/b-hd.ppm"
  });
  expect(fetchMock).toHaveBeenCalledTimes(1);
});
