import nextConfig from "../next.config";

test("rewrites api requests to the configured internal target", async () => {
  process.env.INTERNAL_API_PROXY_TARGET = "http://api:8080";

  const rewrites = await nextConfig.rewrites?.();
  expect(rewrites).toEqual([
    {
      source: "/api/:path*",
      destination: "http://api:8080/api/:path*"
    }
  ]);
});
