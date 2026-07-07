<script>
  import { routePath } from "../../lib/config.js";
  import { getSession, signOut } from "../../lib/session.svelte.js";

  let {
    contentClass = "",
    alertMsg = $bindable(""),
    alertType = "info",
    children,
  } = $props();

  const session = getSession();

  // derive current relative path for active nav highlighting
  let currentPath = $state(
    (() => {
      let p = window.location.pathname;
      if (p.startsWith(routePath)) {
        p = p.slice(routePath.length) || "/";
      }
      return p;
    })()
  );

  /** @param {string} path */
  function navigateTo(path) {
    window.__navigate(routePath + path);
  }

  /** @param {MouseEvent} e @param {string} path */
  function onNavClick(e, path) {
    e.preventDefault();
    navigateTo(path);
  }
</script>

<div class="min-vh-100 d-flex flex-column">
  <!-- Navbar -->
  <nav
    class="navbar navbar-expand bg-dark border-bottom border-body"
    data-bs-theme="dark"
  >
    <div class="container" style="max-width:1100px">
      <a
        class="navbar-brand"
        href="{routePath}/admin/users"
        onclick={(e) => onNavClick(e, "/admin/users")}
        style="display:flex; align-items:center; gap:1rem; margin-right:2rem;"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="32"
          height="32"
          fill="currentColor"
          class="bi bi-gear-fill"
          viewBox="0 0 16 16"
        >
          <path
            d="M9.405 1.05c-.413-1.4-2.397-1.4-2.81 0l-.1.34a1.464 1.464 0 0 1-2.105.872l-.51-.212c-1.353-.56-2.66.907-1.95 2.03l.23.334c.264.356.246.85-.072 1.18l-.27.293a1.466 1.466 0 0 0 0 2.03l.27.293c.318.33.334.824.072 1.18l-.23.334c-.71 1.123.597 2.59 1.95 2.03l.51-.212a1.464 1.464 0 0 1 2.105.872l.1.34c.413 1.4 2.397 1.4 2.81 0l.1-.34a1.464 1.464 0 0 1 2.105-.872l.51.212c1.353.56 2.66-.907 1.95-2.03l-.23-.334a1.464 1.464 0 0 0-.072-1.18l.27-.293a1.466 1.466 0 0 0 0-2.03l-.27-.293a1.464 1.464 0 0 1 .072-1.18l.23-.334c.71-1.123-.597-2.59-1.95-2.03l-.51.212a1.464 1.464 0 0 1-2.105-.872zM8 5a3 3 0 1 1 0 6 3 3 0 0 1 0-6"
          />
        </svg>
        <span>Admin</span>
      </a>
      <ul class="navbar-nav me-auto">
        <li class="nav-item">
          <a
            class="nav-link nav-dot"
            class:active={currentPath === "/admin/users"}
            href="{routePath}/admin/users"
            onclick={(e) => onNavClick(e, "/admin/users")}>Users</a
          >
        </li>
      </ul>
      <div class="d-flex align-items-center gap-2">
        <a
          class="btn btn-outline-light btn-sm"
          href="{routePath}/user/profile"
          onclick={(e) => onNavClick(e, "/user/profile")}>Account</a
        >
        <button class="btn btn-outline-light btn-sm" onclick={() => signOut()}
          >Sign Out</button
        >
      </div>
    </div>
  </nav>

  <!-- Alert -->
  {#if alertMsg}
    <div class="container mt-3" style="max-width:1100px">
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
    </div>
  {/if}

  <!-- Main Content -->
  <div class="container py-4 flex-grow-1" style="max-width:1100px">
    <div class="row justify-content-center">
      <div class={contentClass}>
        {@render children()}
      </div>
    </div>
  </div>

  <!-- Footer -->
  <footer class="text-center py-3 mt-auto">
    <p class="mb-0 text-muted small">
      <a
        href="https://github.com/hooto/iam"
        target="_blank"
        class="text-decoration-none">Powered by hooto IAM</a
      >
    </p>
  </footer>
</div>

<style>
  .nav-dot {
    position: relative;
  }
  .nav-dot.active::after {
    content: "";
    position: absolute;
    bottom: 2px;
    left: 50%;
    transform: translateX(-50%);
    width: 4px;
    height: 4px;
    border-radius: 50%;
    background: currentColor;
  }
</style>
