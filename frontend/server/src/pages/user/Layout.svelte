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
        href="{routePath}/"
        onclick={(e) => onNavClick(e, "/")}
        style="display:flex; align-items:center; gap:1rem; margin-right:2rem;"
        ><svg
          xmlns="http://www.w3.org/2000/svg"
          width="32"
          height="32"
          fill="currentColor"
          class="bi bi-person-circle"
          viewBox="0 0 16 16"
        >
          <path d="M11 6a3 3 0 1 1-6 0 3 3 0 0 1 6 0" />
          <path
            fill-rule="evenodd"
            d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8m8-7a7 7 0 0 0-5.468 11.37C3.242 11.226 4.805 10 8 10s4.757 1.225 5.468 2.37A7 7 0 0 0 8 1"
          />
        </svg> <span>Account</span></a
      >
      <ul class="navbar-nav me-auto">
        <li class="nav-item">
          <a
            class="nav-link nav-dot"
            class:active={currentPath === "/user/profile" ||
              currentPath === "/"}
            href="{routePath}/user/profile"
            onclick={(e) => onNavClick(e, "/user/profile")}>Profile</a
          >
        </li>
        <li class="nav-item">
          <a
            class="nav-link nav-dot"
            class:active={currentPath === "/user/keys"}
            href="{routePath}/user/keys"
            onclick={(e) => onNavClick(e, "/user/keys")}
            >Keys</a
          >
        </li>
        <li class="nav-item">
          <a
            class="nav-link nav-dot"
            class:active={currentPath === "/user/apps"}
            href="{routePath}/user/apps"
            onclick={(e) => onNavClick(e, "/user/apps")}
            >Apps</a
          >
        </li>
      </ul>
      <div class="d-flex">
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
