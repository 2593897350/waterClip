import { HomeUpload } from "../components/home-upload";

export default function HomePage() {
  return (
    <main>
      <h1>Remove watermarks from photos</h1>
      <p>Auto-detect first, then refine the mask yourself.</p>
      <HomeUpload />
    </main>
  );
}
