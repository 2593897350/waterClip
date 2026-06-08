"use client";

import { useEffect, useState } from "react";

import { getJob } from "../lib/api";
import type { ProcessJob } from "../lib/types";
import { ResultCompare } from "./result-compare";

export function ResultClient({ jobId }: { jobId: string }) {
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

  if (!job?.resultPath) {
    return <p>Processing...</p>;
  }

  return (
    <ResultCompare
      onRetry={() => {
        window.location.href = `/edit/${jobId}`;
      }}
      originalUrl={job.sourcePath ?? "/original.ppm"}
      resultUrl={job.resultPath}
    />
  );
}
