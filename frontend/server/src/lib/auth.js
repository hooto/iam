// @ts-nocheck
/**
 * Auth utilities for session management.
 */

/**
 * Decode the payload section of a JWT token string.
 * Returns the parsed claims object or null on failure.
 */
export function parseJwtPayload(token) {
  if (!token || typeof token !== "string") return null;
  const parts = token.split(".");
  if (parts.length < 2) return null;
  try {
    const payload = parts[1].replace(/-/g, "+").replace(/_/g, "/");
    return JSON.parse(decodeURIComponent(escape(atob(payload))));
  } catch {
    return null;
  }
}