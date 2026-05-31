<script>
  import { routePath } from "../../../lib/config.js";

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

  /** @param {SubmitEvent} e */
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
        routePath + "/v2/auth/password/reset-confirm",
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
        window.location.href = routePath + "/auth/sign-in";
      }, 2000);
    } catch (err) {
      alertMsg = "Network error, please try again";
      alertType = "danger";
    } finally {
      loading = false;
    }
  }

  /** @param {MouseEvent} e */
  function goToSignIn(e) {
    e.preventDefault();
    window.__navigate(routePath + "/auth/sign-in");
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
          class="bi bi-key-fill"
          viewBox="0 0 16 16"
        >
          <path
            d="M3.5 11.5a3.5 3.5 0 1 1 3.163-5H14L15.5 8 14 9.5l-1-1-1 1-1-1-1 1-1-1-1 1H6.663a3.5 3.5 0 0 1-3.163 2M2.5 9a1 1 0 1 0 0-2 1 1 0 0 0 0 2"
          />
        </svg>
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
        <a href="{routePath}/auth/sign-in" onclick={goToSignIn}
          >Back to Sign In</a
        >
      </p>
    </div>
  </div>

  <p class="text-white-50 small mt-3 mb-0">
    <a
      href="https://github.com/hooto/iam"
      target="_blank"
      class="text-white text-decoration-none">Powered by hooto IAM</a
    >
  </p>
</div>