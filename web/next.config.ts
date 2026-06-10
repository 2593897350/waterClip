import type { NextConfig } from "next";

const ignoreTypecheck = process.env.NEXT_IGNORE_TYPECHECK === "true";

const nextConfig: NextConfig = {
  typescript: {
    ignoreBuildErrors: ignoreTypecheck
  },
  async rewrites() {
    const internalApiProxyTarget = process.env.INTERNAL_API_PROXY_TARGET ?? "http://localhost:8080";
    return [
      {
        source: "/api/:path*",
        destination: `${internalApiProxyTarget}/api/:path*`
      }
    ];
  }
};

export default nextConfig;
