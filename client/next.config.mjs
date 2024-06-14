/** @type {import('next').NextConfig} */
const nextConfig = {
  async redirects() {
    return [
      {
        source: "/",
        destination: "/files",
        permanent: true,
      },
    ];
  },
};

export default nextConfig;
