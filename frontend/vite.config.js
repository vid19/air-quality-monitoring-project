import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // This forwards any request starting with /api
      // to your Gin server
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
