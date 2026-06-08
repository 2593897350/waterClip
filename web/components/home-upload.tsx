"use client";

import { useRouter } from "next/navigation";

import { createDetectJob } from "../lib/api";
import { UploadForm } from "./upload-form";

export function HomeUpload() {
  const router = useRouter();

  return (
    <UploadForm
      onUploaded={async (file) => {
        const job = await createDetectJob(file);
        router.push(`/edit/${job.jobId}`);
        return job.jobId;
      }}
    />
  );
}
