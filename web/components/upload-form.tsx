"use client";

import React from "react";
import { useState } from "react";

type UploadFormProps = {
  onUploaded: (file: File) => Promise<string>;
};

export function UploadForm({ onUploaded }: UploadFormProps) {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState<"idle" | "detecting">("idle");

  return (
    <form
      onSubmit={async (event) => {
        event.preventDefault();
        if (!file) {
          return;
        }
        setStatus("detecting");
        await onUploaded(file);
      }}
    >
      <label htmlFor="upload">Upload image</label>
      <input
        id="upload"
        name="upload"
        type="file"
        accept="image/png,image/jpeg,image/webp"
        onChange={(event) => setFile(event.target.files?.[0] ?? null)}
      />
      <button type="submit">Start removal</button>
      {status === "detecting" ? <p>Detecting watermark...</p> : null}
    </form>
  );
}
