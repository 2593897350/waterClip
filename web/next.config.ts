import type { NextConfig } from "next";

const nextConfig: NextConfig = {
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
