"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { getJob, processJob } from "../lib/api";
import type { ProcessJob } from "../lib/types";
import { MaskEditor } from "./mask-editor";

export function EditClient({ jobId }: { jobId: string }) {
  const router = useRouter();
  const [job, setJob] = useState<ProcessJob | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getJob(jobId)
      .then(setJob)
      .catch((err: Error) => setError(err.message));
  }, [jobId]);

  if (error) {
    return <p>{error}</p>;
  }

  if (!job) {
    return <p>Loading detect result...</p>;
  }

  return (
    <MaskEditor
      onProcess={async ({ mode }) => {
        if (!job.sourcePath || !job.maskPath) {
          throw new Error("Missing source or mask path");
        }
        const processResult = await processJob(job.sourcePath, job.maskPath, mode);
        router.push(`/result/${processResult.jobId}`);
      }}
    />
  );
}
