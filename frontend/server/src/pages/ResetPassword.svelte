<script>
  import { routePath } from "../lib/config.js";

  let token = $state("");
  let password = $state("");
  let confirmPassword = $state("");
  let alertMsg = $state("");
  let alertType = $state("info");
  let loading = $state(false);
  let success = $state(false);

  // check for token in URL query params on mount
  $effect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const t = urlParams.get("token");
    if (t) {
      token = t;
    }
  });

  async function onSubmit(e) {
    e.preventDefault();
    alertMsg = "";

    if (!token.trim()) {
      alertMsg = "Please enter the verification code";
      alertType = "warning";
      return;
    }

    if (!password.trim()) {
      alertMsg = "Please enter a new password";
      alertType = "warning";
      return;
    }
    if (password.length < 8) {
      alertMsg = "Password must be at least 8 characters";
      alertType = "warning";
      return;
    }
    if (password !== confirmPassword) {
      alertMsg = "Passwords do not match";
      alertType = "warning";
      return;
    }

    loading = true;

    try {
      const resp = await fetch(
        routePath + "/v2/service/reset-password-confirm",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            token: token,
            password: password,
          }),
        },
      );

      const data = await resp.json();

      if (data.status?.code !== "200") {
        alertMsg = data.status?.message || "Reset failed";
        alertType = "danger";
        return;
      }

      success = true;
      alertMsg =
        "Your password has been reset successfully! Redirecting to sign in...";
      alertType = "success";

      setTimeout(() => {
        window.location.href = routePath + "/service/sign-in";
      }, 2000);
    } catch (err) {
      alertMsg = "Network error, please try again";
      alertType = "danger";
    } finally {
      loading = false;
    }
  }

  function goToSignIn(e) {
    e.preventDefault();
    window.__navigate(routePath + "/service/sign-in");
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
          class="bi bi-shield-lock"
          viewBox="0 0 16 16"
          ><path
            d="M5.338 1.59a61 61 0 0 0-2.837.856.081.081 0 0 0-.038.119.48.48 0 0 1 .039.06c.286.478.628.904 1.017 1.27a5.7 5.7 0 0 1 1.819-2.305m7.324 0a5.7 5.7 0 0 1 1.819 2.306c.389-.367.731-.792 1.017-1.27a.48.48 0 0 1 .039-.06.08.08 0 0 0-.038-.119 61 61 0 0 0-2.837-.856zM8 0a6 6 0 0 0-2.714.657C3.67 1.71 2.5 3.22 1.866 4.908.687 7.988.755 11.942 2.542 14.87A2 2 0 0 0 4.283 16h7.434a2 2 0 0 0 1.741-1.129c1.787-2.928 1.855-6.882.676-9.962C10.5 3.22 9.33 1.71 8.714.657A6 6 0 0 0 8 0m0 5a1.5 1.5 0 0 1 .5 2.915l.385 1.99a.5.5 0 0 1-.491.595h-.788a.5.5 0 0 1-.49-.595l.384-1.99A1.5 1.5 0 0 1 8 5"
          /></svg
        >
      </div>
      <p class="text-secondary mb-0">Set new password</p>
    </div>

    {#if success}
      <div class="mb-4">
        <div class="alert alert-success" role="alert">
          {alertMsg}
        </div>
      </div>
    {:else}
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
          <label for="token" class="form-label text-muted small"
            >Verification code</label
          >
          <input
            type="text"
            class="form-control"
            id="token"
            bind:value={token}
            placeholder="Enter the code from your email"
          />
        </div>

        <div class="mb-3">
          <label for="password" class="form-label text-muted small"
            >New password</label
          >
          <input
            type="password"
            class="form-control"
            id="password"
            bind:value={password}
            placeholder="Password (8-30 characters)"
            autocomplete="new-password"
          />
        </div>

        <div class="mb-3">
          <label for="confirmPassword" class="form-label text-muted small"
            >Confirm password</label
          >
          <input
            type="password"
            class="form-control"
            id="confirmPassword"
            bind:value={confirmPassword}
            placeholder="Confirm new password"
            autocomplete="new-password"
          />
        </div>

        <button
          type="submit"
          class="btn btn-primary w-100 py-2 fw-medium"
          disabled={loading}
        >
          {#if loading}
            <span class="spinner-border spinner-border-sm me-2" role="status"
            ></span>
            Resetting...
          {:else}
            Reset password
          {/if}
        </button>
      </form>
    {/if}

    <div class="text-center border-top pt-3">
      <p class="mb-0">
        <a href="{routePath}/service/sign-in" onclick={goToSignIn}
          >Back to Sign In</a
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