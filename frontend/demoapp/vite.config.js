import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import purgecss from "vite-plugin-purgecss";

export default defineConfig({
  plugins: [
    svelte(),
    purgecss({
      content: [
        "./src/**/*.html",
        "./src/**/*.svelte",
        "./src/**/*.jsx",
        "./src/**/*.tsx",
      ],
    }),
  ],
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
