import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

import { MaskEditor } from "../components/mask-editor";

test("switches modes and emits the selected output mode", async () => {
  const user = userEvent.setup();
  const calls: string[] = [];

  render(<MaskEditor onProcess={async ({ mode }) => calls.push(mode)} />);
  await user.click(screen.getByRole("radio", { name: /hd/i }));
  await user.click(screen.getByRole("button", { name: /process image/i }));

  expect(calls).toEqual(["hd"]);
});
