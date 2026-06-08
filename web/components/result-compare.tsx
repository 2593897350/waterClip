type ResultCompareProps = {
  originalUrl: string;
  resultUrl: string;
  onRetry: () => void;
};

export function ResultCompare({ originalUrl, resultUrl, onRetry }: ResultCompareProps) {
  return (
    <section>
      <img alt="Original image" src={originalUrl} />
      <img alt="Processed image" src={resultUrl} />
      <a download href={resultUrl}>
        Download result
      </a>
      <button onClick={onRetry} type="button">
        Refine mask
      </button>
    </section>
  );
}
