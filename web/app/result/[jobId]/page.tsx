import { ResultClient } from "../../../components/result-client";

export default async function ResultPage({ params }: { params: Promise<{ jobId: string }> }) {
  const { jobId } = await params;
  return <ResultClient jobId={jobId} />;
}
