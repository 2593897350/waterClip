import { UploadForm } from "../components/upload-form";
import { createDetectJob } from "../lib/api";

export default function HomePage() {
  return (
    <main>
      <h1>Remove watermarks from photos</h1>
      <p>Auto-detect first, then refine the mask yourself.</p>
      <UploadForm onUploaded={async (file) => createDetectJob(file).then((job) => job.jobId)} />
    </main>
  );
}
