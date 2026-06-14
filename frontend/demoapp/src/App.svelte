<script>
  import { onMount } from "svelte";
  import { userAuthSession, signOut } from "./lib/api.js";
  import Dashboard from "./pages/Dashboard.svelte";

  const REDIRECT_URI_KEY = "user_auth_sign_in_redirect_uri";

  let page = $state("loading"); // loading | error | dashboard
  let authInfo = $state(null);
  let alertMsg = $state("");

  onMount(async () => {
    await loadApp();
  });

  async function loadApp() {
    try {
      const data = await userAuthSession();

      // logged in: AuthClaims is present
      if (data.auth_claims) {
        authInfo = data;

        // check if there is a saved redirect URI from a previous session
        const redirectUri = sessionStorage.getItem(REDIRECT_URI_KEY);
        if (redirectUri && redirectUri !== window.location.pathname) {
          sessionStorage.removeItem(REDIRECT_URI_KEY);
          window.location.href = redirectUri;
          return;
        }
        sessionStorage.removeItem(REDIRECT_URI_KEY);

        page = "dashboard";
        return;
      }

      // not logged in but IAM sign-in URL is available: redirect to sign-in
      if (data.auth_sign_in_url) {
        // save current path for post-login redirect
        sessionStorage.setItem(REDIRECT_URI_KEY, window.location.pathname);
        window.location.href = data.auth_sign_in_url;
        return;
      }

      // no IAM sign-in URL configured
      alertMsg = "IAM service is not configured. Please check the app settings.";
      page = "error";
    } catch (err) {
      alertMsg = "Failed to load authentication info: " + (err.message || "Unknown error");
      page = "error";
    }
  }

  async function handleSignOut() {
    try {
      await signOut();
    } catch {
      // ignore
    }
    authInfo = null;
    page = "loading";
    loadApp();
  }
</script>

{#if page === "loading"}
  <div class="min-vh-100 d-flex align-items-center justify-content-center">
    <div class="spinner-border text-primary" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
{:else if page === "error"}
  <div class="min-vh-100 bg-light">
    <div class="container py-4">
      <div class="alert alert-warning alert-dismissible fade show">
        {alertMsg}
        <button
          type="button"
          class="btn-close"
          aria-label="Close"
          onclick={() => (alertMsg = "")}
        ></button>
      </div>
    </div>
  </div>
{:else if page === "dashboard"}
  <Dashboard
    username={authInfo?.auth_claims?.sub || ""}
    onSignOut={handleSignOut}
  />
{/if}
