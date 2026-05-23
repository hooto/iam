<script>
  import { routePath } from "../lib/config.js";

  let username = $state("");
  let email = $state("");
  let password = $state("");
  let confirmPassword = $state("");
  let alertMsg = $state("");
  let alertType = $state("info");
  let loading = $state(false);
  let signUpDisabled = $state(false);

  $effect(() => {
    checkSignUpAllowed();
  });

  async function checkSignUpAllowed() {
    try {
      const resp = await fetch(routePath + "/v2/sys/info");
      const data = await resp.json();
      if (data.status?.code === "200" && data.allow_user_sign_up !== true) {
        signUpDisabled = true;
        alertMsg = "Self-registration is not allowed. Please contact the administrator to create an account.";
        alertType = "warning";
      }
    } catch (err) {
      // network error, allow sign up attempt
    }
  }

  async function onSubmit(e) {
    e.preventDefault();
    if (signUpDisabled) return;
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
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      alertMsg = "Invalid email format";
      alertType = "warning";
      return;
    }
    if (!password.trim()) {
      alertMsg = "Please enter your password";
      alertType = "warning";
      return;
    }
    if (password.length < 6) {
      alertMsg = "Password must be at least 6 characters";
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
      const resp = await fetch(routePath + "/v2/service/sign-up", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: username,
          email: email,
          password: password,
        }),
      });

      const data = await resp.json();

      if (data.status?.code !== "200") {
        alertMsg = data.status?.message || "Sign up failed";
        alertType = "danger";
        return;
      }

      alertMsg = "Account created successfully! Redirecting to sign in...";
      alertType = "success";

      setTimeout(() => {
        window.location.href = routePath + "/service/sign-in";
      }, 1500);
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
          class="bi bi-person-plus-fill"
          viewBox="0 0 16 16"
          ><path
            d="M1 14s-1 0-1-1 1-4 6-4 6 3 6 4-1 1-1 1zm5-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6"
          /><path
            fill-rule="evenodd"
            d="M13.5 5a.5.5 0 0 1 .5.5V7h1.5a.5.5 0 0 1 0 1H14v1.5a.5.5 0 0 1-1 0V8h-1.5a.5.5 0 0 1 0-1H13V5.5a.5.5 0 0 1 .5-.5"
          /></svg
        >
      </div>
      <p class="text-secondary mb-0">Create a new account</p>
    </div>

    <form class="mb-4" onsubmit={onSubmit}>
      {#if alertMsg}
        <div
          class="alert alert-{alertType} {signUpDisabled ? '' : 'alert-dismissible'} fade show"
          role="alert"
        >
          {alertMsg}
          {#if !signUpDisabled}
            <button
              type="button"
              class="btn-close"
              aria-label="Close"
              onclick={() => (alertMsg = "")}
            ></button>
          {/if}
        </div>
      {/if}

      {#if !signUpDisabled}
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
          placeholder="Email"
          autocomplete="email"
        />
      </div>

      <div class="mb-3">
        <input
          type="password"
          class="form-control"
          id="password"
          bind:value={password}
          placeholder="Password (at least 6 characters)"
          autocomplete="new-password"
        />
      </div>

      <div class="mb-3">
        <input
          type="password"
          class="form-control"
          id="confirmPassword"
          bind:value={confirmPassword}
          placeholder="Confirm password"
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
          Signing up...
        {:else}
          Sign Up
        {/if}
      </button>
      {/if}
    </form>

    <div class="text-center border-top pt-3">
      <p class="mb-0">
        <span class="text-muted">Already have an account? </span>
        <a href="{routePath}/service/sign-in" onclick={goToSignIn}>Sign in</a>
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
