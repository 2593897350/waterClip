"use client";

import { useState } from "react";

import type { Mode } from "../lib/types";

type ProcessArgs = {
  mode: Mode;
};

type MaskEditorProps = {
  onProcess: (args: ProcessArgs) => Promise<void>;
};

export function MaskEditor({ onProcess }: MaskEditorProps) {
  const [mode, setMode] = useState<Mode>("fast");

  return (
    <section>
      <h2>Edit detected area</h2>
      <div aria-label="output mode" role="group">
        <label>
          <input checked={mode === "fast"} name="mode" onChange={() => setMode("fast")} type="radio" />
          Fast
        </label>
        <label>
          <input checked={mode === "hd"} name="mode" onChange={() => setMode("hd")} type="radio" />
          HD
        </label>
      </div>
      <button onClick={() => onProcess({ mode })} type="button">
        Process image
      </button>
    </section>
  );
}
