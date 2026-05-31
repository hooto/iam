// @ts-nocheck
/**
 * API utility module for v2 backend endpoints.
 * All requests use credentials: "same-origin" for cookie-based auth.
 */

import { routePath } from "./config.js";

/**
 * Core request helper. Appends routePath prefix and checks response status.
 * @param {string} path - API path relative to routePath
 * @param {RequestInit} [options] - fetch options
 * @param {string} [errMsg] - fallback error message
 */
async function request(path, options = {}, errMsg = "Request failed") {
  const opts = { credentials: "same-origin", ...options };
  if (opts.body && typeof opts.body === "string") {
    opts.headers = { "Content-Type": "application/json", ...opts.headers };
  }
  const resp = await fetch(routePath + path, opts);
  const data = await resp.json();
  if (data.status?.code !== "200") {
    throw new Error(data.status?.message || errMsg);
  }
  return data;
}

// User Profile API

export async function fetchProfile() {
  const data = await request("/v2/user/profile", {}, "Failed to fetch profile");
  return data.item;
}

export async function updateProfile(params) {
  return request("/v2/user/profile-set", {
    method: "POST",
    body: JSON.stringify(params),
  }, "Failed to update profile");
}

export async function changePassword(params) {
  return request("/v2/user/pass-set", {
    method: "POST",
    body: JSON.stringify(params),
  }, "Failed to change password");
}

export async function changeEmail(params) {
  return request("/v2/user/email-set", {
    method: "POST",
    body: JSON.stringify(params),
  }, "Failed to change email");
}

export async function uploadPhoto(dataUrl) {
  return request("/v2/user/photo-set", {
    method: "POST",
    body: JSON.stringify({ data: dataUrl }),
  }, "Failed to upload photo");
}

// Access Key Management API

export async function fetchAccessKeys() {
  const data = await request("/v2/user/keys/list", {}, "Failed to fetch access keys");
  return data.items || [];
}

export async function fetchAccessKey(id) {
  const data = await request(
    "/v2/user/keys/entry?access_key_id=" + encodeURIComponent(id),
    {},
    "Failed to fetch access key",
  );
  return data.item;
}

export async function setAccessKey(req) {
  return request("/v2/user/keys/set", {
    method: "PUT",
    body: JSON.stringify(req),
  }, "Failed to save access key");
}

export async function deleteAccessKey(id) {
  return request(
    "/v2/user/keys/delete?access_key_id=" + encodeURIComponent(id),
    {},
    "Failed to delete access key",
  );
}

function scopeParams(id, scopeContent) {
  return "access_key_id=" + encodeURIComponent(id) +
    "&scope_content=" + encodeURIComponent(scopeContent);
}

export async function bindAccessKeyScope(id, scopeContent) {
  return request("/v2/user/keys/bind?" + scopeParams(id, scopeContent), {}, "Failed to bind scope");
}

export async function unbindAccessKeyScope(id, scopeContent) {
  return request("/v2/user/keys/unbind?" + scopeParams(id, scopeContent), {}, "Failed to unbind scope");
}

// App Management API

export async function fetchApps() {
  const data = await request("/v2/user/apps/list", {}, "Failed to fetch apps");
  return data.items || [];
}

export async function createApp(name, url) {
  return request("/v2/user/apps/register", {
    method: "POST",
    body: JSON.stringify({ name, url }),
  }, "Failed to create app");
}

export async function updateApp(id, name, url) {
  return request("/v2/user/apps/update", {
    method: "POST",
    body: JSON.stringify({ id, name, url }),
  }, "Failed to update app");
}

export async function deleteApp(id) {
  return request("/v2/user/apps/delete", {
    method: "POST",
    body: JSON.stringify({ id }),
  }, "Failed to delete app");
}