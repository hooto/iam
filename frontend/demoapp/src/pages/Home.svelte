<script>
  let { config } = $props();

  async function onLogin() {
    if (!config?.auth_sign_in_url) {
      return;
    }
    // save current path for post-login redirect
    sessionStorage.setItem("user_auth_sign_in_redirect_uri", window.location.pathname);

    window.location.href = config.auth_sign_in_url;
  }
</script>

<div class="min-vh-100 d-flex align-items-center justify-content-center bg-light">
  <div class="w-100 bg-white rounded-4 shadow p-4" style="max-width: 480px;">
    <div class="text-center mb-4">
      <h3 class="fw-bold">Demo App</h3>
      <p class="text-secondary">
        A sample application demonstrating IAM integration
      </p>
    </div>

    <div class="text-center mb-3">
      <p class="text-secondary">
        Click below to sign in through the IAM service
      </p>
      <div class="small text-muted mb-3">
        <p>IAM: {config?.auth_base_url || "N/A"}</p>
        <p>App: {config?.app_id || "N/A"}</p>
      </div>
    </div>

    <button
      class="btn btn-primary w-100 py-2"
      onclick={onLogin}
      disabled={!config?.auth_sign_in_url}
    >
      Sign in with IAM
    </button>
  </div>
</div>
