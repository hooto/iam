// @ts-nocheck
/**
 * Clipboard utility with fallback for non-secure contexts (HTTP).
 *
 * navigator.clipboard is undefined when the page is served over HTTP on a
 * non-localhost origin. This module transparently degrades to the legacy
 * execCommand("copy") path so copy actions never throw synchronously.
 */

/**
 * Copy text to clipboard, falling back to a hidden textarea + execCommand
 * when the async Clipboard API is unavailable or rejected.
 * @param {string} text - the text to copy
 * @returns {Promise<boolean>} resolves true on success, false on failure
 */
export async function copyToClipboard(text) {
  // Prefer the modern async Clipboard API (only available in secure contexts).
  if (navigator.clipboard && typeof navigator.clipboard.writeText === "function") {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch {
      // Fall through to the legacy approach below.
    }
  }

  // Legacy fallback: a transient textarea selected and copied via execCommand.
  try {
    const ta = document.createElement("textarea");
    ta.value = text;
    ta.setAttribute("readonly", "");
    ta.style.position = "fixed";
    ta.style.top = "0";
    ta.style.left = "0";
    ta.style.opacity = "0";
    document.body.appendChild(ta);
    ta.focus();
    ta.select();
    const ok = document.execCommand("copy");
    document.body.removeChild(ta);
    return ok;
  } catch {
    return false;
  }
}