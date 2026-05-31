<script>
  import { signOut } from "../lib/api.js";

  let { username, onSignOut } = $props();

  let loading = $state(false);

  async function onSignOutClick() {
    loading = true;
    try {
      await signOut();
    } catch {
      // ignore
    }
    onSignOut();
  }
</script>

<div class="min-vh-100 bg-light">
  <nav class="navbar navbar-expand bg-white shadow-sm">
    <div class="container">
      <span class="navbar-brand fw-bold">Demo App</span>
      <div class="d-flex align-items-center">
        <span class="text-secondary me-3">{username}</span>
        <button
          class="btn btn-outline-secondary btn-sm"
          onclick={onSignOutClick}
          disabled={loading}
        >
          {#if loading}
            <span class="spinner-border spinner-border-sm"></span>
          {:else}
            Sign Out
          {/if}
        </button>
      </div>
    </div>
  </nav>

  <div class="container py-4">
    <div class="row g-4">
      <div class="col-md-6">
        <div class="card">
          <div class="card-header">User Info</div>
          <div class="card-body">
            <table class="table table-borderless mb-0">
              <tbody>
                <tr>
                  <th class="text-secondary ps-0" style="width: 120px">Username</th>
                  <td>{username}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div class="col-md-6">
        <div class="card">
          <div class="card-header">How It Works</div>
          <div class="card-body">
            <ol class="mb-0">
              <li>App detects user is not logged in</li>
              <li>Current page path is saved to sessionStorage</li>
              <li>User is redirected to IAM sign-in with <code>app_id</code></li>
              <li>User authenticates on the IAM server</li>
              <li>IAM generates a one-time auth <code>code</code></li>
              <li>IAM redirects back to the app callback with the <code>code</code></li>
              <li>App backend exchanges <code>code</code> for <code>access_token</code></li>
              <li>App sets <code>http-only</code> cookie with the token</li>
              <li>User is redirected to the originally saved page</li>
            </ol>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>