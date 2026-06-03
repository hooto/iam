import { routePath } from "./config.js";

const DEFAULT_TIMEOUT = 15_000;

/**
 * Core request helper with auto JSON headers, timeout, and status checking.
 *
 * Custom options (extracted before passing to fetch):
 *   - errMsg  {string}  Fallback error message (default: "Request failed")
 *   - raw     {boolean} If true, skip status check and return raw JSON (default: false)
 *   - timeout {number}  Request timeout in milliseconds (default: 15000)
 *
 * All other options are forwarded to fetch() as RequestInit.
 *
 * @param {string} path   - API path relative to routePath
 * @param {object} options - fetch options mixed with custom options above
 */
async function request(path, options = {}) {
  const {
    errMsg = "Request failed",
    raw = false,
    timeout = DEFAULT_TIMEOUT,
    ...fetchOpts
  } = options;

  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeout);

  fetchOpts.credentials = "same-origin";
  fetchOpts.signal = controller.signal;

  if (fetchOpts.body && typeof fetchOpts.body === "string") {
    fetchOpts.headers = { "Content-Type": "application/json", ...fetchOpts.headers };
  }

  try {
    const resp = await fetch(routePath + path, fetchOpts);

    let data;
    try {
      data = await resp.json();
    } catch {
      throw new Error(errMsg);
    }

    if (!raw && data.status?.code !== "200") {
      throw new Error(data.status?.message || errMsg);
    }

    return data;
  } catch (err) {
    if (err.name === "AbortError") {
      throw new Error(`Request timeout after ${timeout}ms: ${path}`);
    }
    throw err;
  } finally {
    clearTimeout(timer);
  }
}

// Fetch user auth session without throwing on auth errors (401).
// Uses raw mode so the caller can inspect auth_claims and auth_endpoint.
export async function userAuthSession() {
  return request("/api/user-auth/session", {
    method: "POST",
    body: JSON.stringify({ current_url: window.location.pathname }),
    raw: true,
  });
}

export async function signOut() {
  return request("/api/user-auth/sign-out", {
    method: "POST",
    errMsg: "Sign out failed",
  });
}