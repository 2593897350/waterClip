import { render, screen } from "@testing-library/react";

import { ResultCompare } from "../components/result-compare";

test("shows original and processed image download action", () => {
  render(
    <ResultCompare
      onRetry={() => undefined}
      originalUrl="/original.ppm"
      resultUrl="/result.ppm"
    />,
  );

  expect(screen.getByRole("img", { name: /original image/i })).toBeInTheDocument();
  expect(screen.getByRole("img", { name: /processed image/i })).toBeInTheDocument();
  expect(screen.getByRole("link", { name: /download result/i })).toBeInTheDocument();
});
