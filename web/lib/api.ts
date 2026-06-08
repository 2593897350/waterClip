import type { DetectJob, Mode, ProcessJob } from "./types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export async function createDetectJob(file: File): Promise<DetectJob> {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetch(`${API_BASE_URL}/api/jobs/detect`, {
    method: "POST",
    body: formData
  });
  if (!response.ok) {
    throw new Error("Failed to create detect job");
  }

  const body = await response.json();
  return {
    jobId: body.job_id,
    status: body.status,
    maskPath: body.mask_path,
    sourcePath: body.source_path
  };
}

export async function processJob(sourcePath: string, maskPath: string, mode: Mode): Promise<ProcessJob> {
  const response = await fetch(`${API_BASE_URL}/api/jobs/process`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      source_path: sourcePath,
      mask_path: maskPath,
      mode
    })
  });
  if (!response.ok) {
    throw new Error("Failed to create process job");
  }

  const body = await response.json();
  return {
    jobId: body.job_id,
    status: body.status,
    resultPath: body.result_path
  };
}

export async function getJob(jobId: string): Promise<ProcessJob> {
  const response = await fetch(`${API_BASE_URL}/api/jobs/${jobId}`);
  if (!response.ok) {
    throw new Error("Failed to load job");
  }

  const body = await response.json();
  return {
    jobId: body.job_id,
    sourcePath: body.source_path,
    maskPath: body.mask_path,
    status: body.status,
    resultPath: body.result_path
  };
}
