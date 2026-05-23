// @ts-nocheck
/**
 * 前端路径配置
 *
 * 路径前缀统一由 index.html 中的 <base href="/iam/"> 控制
 * Vite 构建时通过 base 配置（默认 '/iam/'）设置 import.meta.env.BASE_URL
 *
 * 修改路径前缀：同步修改 index.html 的 <base> 和 vite.config.js 的 base
 */

// Vite 在构建时注入，值为 base 配置，如 '/iam/'
export const basePath = import.meta.env.BASE_URL || "/iam/";

// 去掉尾部斜杠，用于路由匹配，如 '/iam'
export const routePath = basePath.replace(/\/+$/, "");
