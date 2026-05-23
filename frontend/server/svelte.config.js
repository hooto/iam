import { vitePreprocess } from '@sveltejs/vite-plugin-svelte'

/** @type {import("@sveltejs/vite-plugin-svelte").SvelteConfig} */
export default {
  // Consult https://svelte.dev/docs#compile-time-svelte-preprocess
  // for more information about preprocessors
  preprocess: vitePreprocess(),
  build: {
    // 1. 指定输出路径到后端的文件夹
    outDir: '../../cmd/server/dist',
    // 2. 每次构建先清空旧文件
    emptyOutDir: true,
    // 3. (可选) 如果你希望文件名不带哈希值，可以额外配置，但默认带哈希对缓存更友好
  }
}
