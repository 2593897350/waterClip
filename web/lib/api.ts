import type { DetectJob } from "./types";

export async function createDetectJob(_file: File): Promise<DetectJob> {
  return { jobId: "job-123" };
}
