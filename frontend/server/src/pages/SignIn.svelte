<script>
  import { routePath } from "../lib/config.js";

  let username = $state("");
  let password = $state("");
  let alertMsg = $state("");
  let alertType = $state("info");
  let loading = $state(false);

  /** @param {SubmitEvent} e */
  async function onSubmit(e) {
    e.preventDefault();
    alertMsg = "";

    if (!username.trim()) {
      alertMsg = "Please enter your username";
      alertType = "warning";
      return;
    }
    if (!password.trim()) {
      alertMsg = "Please enter your password";
      alertType = "warning";
      return;
    }

    loading = true;

    try {
      // pass redirect_token from URL query if present
      const urlParams = new URLSearchParams(window.location.search);
      const redirectToken = urlParams.get("redirect_token");

      const body = {
        username: username,
        password: password,
        ...(redirectToken && { redirect_token: redirectToken }),
      };

      const resp = await fetch(routePath + "/v2/service/user-sign-in", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });

      const data = await resp.json();

      if (!data.status || data.status.code !== "200") {
        alertMsg = (data.status && data.status.message) || "Sign in failed";
        alertType = "danger";
        return;
      }

      alertMsg = "Sign in successful! Redirecting...";
      alertType = "success";

      setTimeout(() => {
        window.location.href = data.redirect_uri || routePath + "/";
      }, 500);
    } catch (err) {
      alertMsg = "Network error, please try again";
      alertType = "danger";
    } finally {
      loading = false;
    }
  }
</script>

<div
  class="bg-auth min-vh-100 d-flex flex-column justify-content-center align-items-center p-3"
>
  <div class="w-100 bg-white rounded-4 shadow p-4" style="max-width: 420px;">
    <div class="text-center mb-4">
      <div class="mb-3">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="96"
          height="96"
          fill="#0d6efd"
          class="bi bi-person-circle"
          viewBox="0 0 16 16"
          ><path d="M11 6a3 3 0 1 1-6 0 3 3 0 0 1 6 0" /><path
            d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8m8-7a7 7 0 0 0-5.468 11.37C3.242 11.226 4.805 10 8 10s4.757 1.225 5.468 2.37A7 7 0 0 0 8 1"
          /></svg
        >
      </div>
      <p class="text-secondary mb-0">Sign in to your account</p>
    </div>

    <form class="mb-4" onsubmit={onSubmit}>
      {#if alertMsg}
        <div
          class="alert alert-{alertType} alert-dismissible fade show"
          role="alert"
        >
          {alertMsg}
          <button
            type="button"
            class="btn-close"
            aria-label="Close"
            onclick={() => (alertMsg = "")}
          ></button>
        </div>
      {/if}

      <div class="mb-3">
        <input
          type="text"
          class="form-control"
          id="username"
          bind:value={username}
          placeholder="Username"
          autocomplete="username"
        />
      </div>

      <div class="mb-3">
        <input
          type="password"
          class="form-control"
          id="password"
          bind:value={password}
          placeholder="Password"
          autocomplete="current-password"
        />
      </div>

      <div class="mb-3 d-flex justify-content-end align-items-center">
        <a
          href="{routePath}/service/forgot-password"
          class="link-primary"
          onclick={(e) => {
            e.preventDefault();
            window.__navigate(routePath + "/service/forgot-password");
          }}>Forgot password?</a
        >
      </div>

      <button
        type="submit"
        class="btn btn-primary w-100 py-2 fw-medium"
        disabled={loading}
      >
        {#if loading}
          <span class="spinner-border spinner-border-sm me-2" role="status"
          ></span>
          Signing in...
        {:else}
          Sign In
        {/if}
      </button>
    </form>

    <div class="text-center border-top pt-3">
      <p class="mb-0">
        <span class="text-muted">Don't have an account? </span>
        <a
          href="{routePath}/service/sign-up"
          onclick={(e) => {
            e.preventDefault();
            window.__navigate(routePath + "/service/sign-up");
          }}>Create account</a
        >
      </p>
    </div>
  </div>

  <p class="text-white-50 small mt-3 mb-0">
    <a
      href="https://github.com/hooto/iam"
      target="_blank"
      class="text-white text-decoration-none">hooto IAM</a
    >
  </p>
</div>
