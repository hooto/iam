import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import purgecss from "vite-plugin-purgecss";

// https://vite.dev/config/
// base 需与 index.html 中 <base href="..."> 保持一致
// 修改路径前缀时，同步修改两处即可：
//   1. index.html 的 <base href="/iam/">
//   2. 此处的 base 配置（也可通过环境变量 VITE_BASE_PATH 覆盖）
const basePath = process.env.VITE_BASE_PATH || "/iam/";

export default defineConfig({
  base: basePath,
  plugins: [
    {
      name: "redirect-base-path",
      configureServer(server) {
        const baseNoSlash = basePath.replace(/\/+$/, "");
        server.middlewares.use((req, res, next) => {
          if (req.url === baseNoSlash) {
            res.writeHead(302, { Location: basePath });
            res.end();
            return;
          }
          next();
        });
      },
    },
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
  server: {
    historyApiFallback: true,
    proxy: {
      "/iam/v2": {
        target: "http://localhost:3000",
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "../../cmd/server/dist",
    emptyOutDir: true,
  },
});
