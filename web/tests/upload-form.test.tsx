import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

import { UploadForm } from "../components/upload-form";

test("submits a file and shows processing state", async () => {
  const user = userEvent.setup();
  render(<UploadForm onUploaded={async () => "job-123"} />);

  const input = screen.getByLabelText(/upload image/i);
  await user.upload(input, new File(["demo"], "sample.png", { type: "image/png" }));
  await user.click(screen.getByRole("button", { name: /start removal/i }));

  expect(screen.getByText(/detecting watermark/i)).toBeInTheDocument();
});
