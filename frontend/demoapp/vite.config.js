import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

export default defineConfig({
  plugins: [svelte()],
  base: "/demoapp/",
  server: {
    port: 5175,
    proxy: {
      "/demoapp/api": "http://localhost:3001",
    },
  },
  build: {
    outDir: "../../cmd/demoapp/dist",
    emptyOutDir: true,
  },
});
