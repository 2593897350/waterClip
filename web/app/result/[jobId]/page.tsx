import { ResultClient } from "../../../components/result-client";

export default function ResultPage({ params }: { params: { jobId: string } }) {
  return <ResultClient jobId={params.jobId} />;
}
