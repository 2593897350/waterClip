import { EditClient } from "../../../components/edit-client";

export default async function EditPage({ params }: { params: Promise<{ jobId: string }> }) {
  const { jobId } = await params;
  return <EditClient jobId={jobId} />;
}
