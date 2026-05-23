<script>
  import { routePath } from "../lib/config.js";

  let username = $state("");
  let email = $state("");
  let alertMsg = $state("");
  let alertType = $state("info");
  let loading = $state(false);
  let submitted = $state(false);

  async function onSubmit(e) {
    e.preventDefault();
    alertMsg = "";

    if (!username.trim()) {
      alertMsg = "Please enter your username";
      alertType = "warning";
      return;
    }
    if (!email.trim()) {
      alertMsg = "Please enter your email";
      alertType = "warning";
      return;
    }

    loading = true;

    try {
      const resp = await fetch(
        routePath + "/v2/service/reset-password-ticket",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            username: username,
            email: email,
          }),
        },
      );

      const data = await resp.json();

      if (data.status?.code !== "200") {
        alertMsg = data.status?.message || "Request failed";
        alertType = "danger";
        return;
      }

      submitted = true;
      alertMsg =
        "If an account with that username and email exists, a verification code has been sent to your email.";
      alertType = "success";
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
          class="bi bi-key-fill"
          viewBox="0 0 16 16"
        >
          <path
            d="M3.5 11.5a3.5 3.5 0 1 1 3.163-5H14L15.5 8 14 9.5l-1-1-1 1-1-1-1 1-1-1-1 1H6.663a3.5 3.5 0 0 1-3.163 2M2.5 9a1 1 0 1 0 0-2 1 1 0 0 0 0 2"
          />
        </svg>
      </div>
      <p class="text-secondary mb-0">Reset your password</p>
    </div>

    {#if submitted}
      <div class="mb-4">
        <div class="alert alert-success" role="alert">
          {alertMsg}
        </div>
        <p class="text-muted small">
          Please check your email for the verification code, then use it to set
          a new password.
        </p>
        <button
          class="btn btn-outline-primary w-100 py-2"
          onclick={() => {
            window.__navigate(routePath + "/service/reset-password");
          }}>Enter verification code</button
        >
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

        <p class="text-muted small mb-3">
          Enter your username and the email address associated with your
          account. We will send you a verification code to reset your password.
        </p>

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
            type="email"
            class="form-control"
            id="email"
            bind:value={email}
            placeholder="Email address"
            autocomplete="email"
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
            Sending...
          {:else}
            Send verification code
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
