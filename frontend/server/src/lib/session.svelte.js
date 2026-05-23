// @ts-nocheck
/**
 * Shared reactive session state for authenticated pages.
 *
 * Uses Svelte 5 runes in a .svelte.js module so that any component
 * importing from here gets reactive updates when session data changes.
 */
import { routePath } from "./config.js";
import { parseJwtPayload } from "./auth.js";

const session = $state({
  username: "",
  loaded: false,
});

/**
 * Fetch the current user session from the server.
 * Safe to call multiple times; subsequent calls are no-ops once loaded.
 */
export async function fetchSession() {
  if (session.loaded) return;
  try {
    const resp = await fetch(routePath + "/v2/service/user-session", {
      credentials: "same-origin",
    });
    const data = await resp.json();
    if (data.status?.code === "200" && data.access_token) {
      const claims = parseJwtPayload(data.access_token);
      if (claims?.sub) {
        session.username = claims.sub;
      }
    }
  } catch {
    // ignore network errors
  }
  session.loaded = true;
}

/**
 * Return the reactive session object.
 * Access `getSession().username` in components for reactive binding.
 */
export function getSession() {
  return session;
}

/**
 * Sign out the current user and redirect to the sign-in page.
 */
export async function signOut() {
  try {
    await fetch(routePath + "/v2/service/user-sign-out", {
      method: "POST",
      credentials: "same-origin",
    });
  } catch {
    // ignore network errors, proceed to redirect
  }
  window.location.href = routePath + "/service/sign-in";
}