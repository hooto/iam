// @ts-nocheck
/**
 * Shared reactive session state for authenticated pages.
 *
 * Uses Svelte 5 runes in a .svelte.js module so that any component
 * importing from here gets reactive updates when session data changes.
 */
import { routePath } from "./config.js";

const session = $state({
  username: "",
  roles: [],
  loaded: false,
});

/**
 * Fetch the current user session from the server.
 * Safe to call multiple times; subsequent calls are no-ops once loaded.
 *
 * The /v2/auth/session response returns DECODED claims directly:
 *   { status, auth_claims: { sub, ... }, identity_token: { roles, ... } }
 * It does NOT return a raw access_token, so we read from those fields.
 */
export async function fetchSession() {
  if (session.loaded) return;
  try {
    const resp = await fetch(routePath + "/v2/auth/session", {
      credentials: "same-origin",
    });
    const data = await resp.json();
    if (data.status?.code === "200" && data.auth_claims) {
      session.username = data.auth_claims.sub || "";
      session.roles = data.identity_token?.roles || [];
    }
  } catch {
    // ignore network errors
  }
  session.loaded = true;
}

/**
 * Return the reactive session object.
 * Access `getSession().username` or `getSession().roles` in components
 * for reactive binding.
 */
export function getSession() {
  return session;
}

/**
 * Whether the current session holds the sysadmin role.
 * Reads the reactive session state, so it stays current after fetch.
 */
export function isAdminSession() {
  return Array.isArray(session.roles) && session.roles.includes("sa");
}

/**
 * Sign out the current user and redirect to the sign-in page.
 */
export async function signOut() {
  try {
    await fetch(routePath + "/v2/auth/sign-out", {
      method: "POST",
      credentials: "same-origin",
    });
  } catch {
    // ignore network errors, proceed to redirect
  }
  window.location.href = routePath + "/auth/sign-in";
}