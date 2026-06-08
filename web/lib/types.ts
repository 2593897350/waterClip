export type Mode = "fast" | "hd";

export type DetectJob = {
  jobId: string;
  status: string;
  maskPath?: string;
  sourcePath?: string;
};

export type ProcessJob = {
  jobId: string;
  status: string;
  sourcePath?: string;
  maskPath?: string;
  resultPath?: string;
};
