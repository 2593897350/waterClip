import { EditClient } from "../../../components/edit-client";

export default function EditPage({ params }: { params: { jobId: string } }) {
  return <EditClient jobId={params.jobId} />;
}
