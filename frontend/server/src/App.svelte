<script>
  // @ts-nocheck
  import Auth_SignIn from "./pages/auth/SignIn.svelte";
  import Auth_SignUp from "./pages/auth/SignUp.svelte";
  import Auth_Password_Forgot from "./pages/auth/password/Forgot.svelte";
  import Auth_Password_Reset from "./pages/auth/password/Reset.svelte";

  import User_Profile from "./pages/user/Profile.svelte";
  import User_Keys from "./pages/user/AccessKey.svelte";
  import User_Apps from "./pages/user/Apps.svelte";

  import Admin_Users from "./pages/admin/Users.svelte";

  import { routePath } from "./lib/config.js";

  let currentRoute = $state(window.location.pathname);
  let authChecked = $state(false);
  let isLoggedIn = $state(false);
  let username = $state("");
  let userRoles = $state([]);
  let isAdmin = $derived(
    username === "sysadmin" || userRoles.includes("sa"),
  );

  // dynamic body background based on auth state
  $effect(() => {
    if (isLoggedIn) {
      document.body.style.background = "var(--bs-gray-200, #e9ecef)";
    } else {
      document.body.style.background = "linear-gradient(135deg, #0d6efd 0%, #0a58ca 100%)";
    }
  });

  const signInPath = "/auth/sign-in";
  const signUpPath = "/auth/sign-up";
  const passwordForgotPath = "/auth/password/forgot";
  const passwordResetPath = "/auth/password/reset";

  /** @type {Record<string, typeof Auth_SignIn>} */
  const authRoutes = {
    [signInPath]: Auth_SignIn,
    [signUpPath]: Auth_SignUp,
    [passwordForgotPath]: Auth_Password_Forgot,
    [passwordResetPath]: Auth_Password_Reset,
  };

  /** @type {Record<string, typeof User_Profile>} */
  const userRoutes = {
    "/user/profile": User_Profile,
    "/user/keys": User_Keys,
    "/user/apps": User_Apps,
  };

  /** @type {Record<string, typeof Admin_Users>} */
  const adminRoutes = {
    "/admin/users": Admin_Users,
  };

  window.addEventListener("popstate", () => {
    currentRoute = window.location.pathname;
  });

  function getRelPath() {
    let rel = currentRoute.startsWith(routePath)
      ? currentRoute.slice(routePath.length)
      : currentRoute;
    return rel || "/";
  }

  // Pure function: resolve which component to render, no side effects.
  // Redirects are handled separately in the $effect below to avoid
  // mutating reactive state during render (state_unsafe_mutation).
  function resolveComponent() {
    const relPath = getRelPath();

    // root path
    if (relPath === "" || relPath === "/") {
      return isLoggedIn ? User_Profile : Auth_SignIn;
    }

    // third-party app sign-in flow: always render the sign-in form,
    // regardless of the local login state.
    if (relPath === signInPath && hasAppIdParam()) {
      return Auth_SignIn;
    }

    if (authRoutes[relPath]) {
      return isLoggedIn ? User_Profile : authRoutes[relPath];
    }

    if (userRoutes[relPath]) {
      return isLoggedIn ? userRoutes[relPath] : Auth_SignIn;
    }

    if (adminRoutes[relPath]) {
      if (!isLoggedIn) return Auth_SignIn;
      // non-admins are not allowed into the admin area
      if (!isAdmin) return User_Profile;
      return adminRoutes[relPath];
    }

    // /admin (and /admin/) is the admin entry; it resolves to the Users page.
    // applyRedirect() normalizes the URL to /admin/users.
    if (relPath === "/admin" || relPath === "/admin/") {
      if (!isLoggedIn) return Auth_SignIn;
      if (!isAdmin) return User_Profile;
      return Admin_Users;
    }

    // default: show sign-in or profile based on auth state
    return isLoggedIn ? User_Profile : Auth_SignIn;
  }

  window.__navigate = (path) => {
    window.history.pushState({}, "", path);
    currentRoute = path;
  };

  // Detect if the current request is a third-party app sign-in flow.
  // When app_id is present, the sign-in form MUST be shown regardless of the
  // local login state, so that the user re-authenticates and is redirected
  // back to the third-party callback-url.
  function hasAppIdParam() {
    const params = new URLSearchParams(window.location.search);
    return !!params.get("app_id");
  }

  // Redirect to the appropriate route based on auth state.
  // Uses replaceState so the initial redirect does not pollute history.
  // replace: true to replace history entry (initial load),
  //          false to push a new entry (runtime navigation).
  function applyRedirect(replace) {
    const relPath = getRelPath();
    const profilePath = routePath + "/user/profile";
    const signInUrl = routePath + signInPath;

    // third-party app sign-in flow: always force the sign-in form,
    // even when the user is already logged in locally.
    if (relPath === signInPath && hasAppIdParam()) {
      return false;
    }

    // root path: redirect based on auth state
    if (relPath === "" || relPath === "/") {
      if (isLoggedIn) {
        if (replace) {
          window.history.replaceState({}, "", profilePath);
          currentRoute = profilePath;
        } else {
          window.__navigate(profilePath);
        }
      } else {
        window.location.href = signInUrl;
      }
      return true;
    }

    // logged-in user on a public/auth route -> redirect to profile
    if (authRoutes[relPath] && isLoggedIn) {
      if (replace) {
        window.history.replaceState({}, "", profilePath);
        currentRoute = profilePath;
      } else {
        window.__navigate(profilePath);
      }
      return true;
    }

    // not logged in on a user route -> redirect to sign-in
    if (userRoutes[relPath] && !isLoggedIn) {
      window.location.href = signInUrl;
      return true;
    }

    // admin routes: require login + sysadmin role
    if (adminRoutes[relPath]) {
      if (!isLoggedIn) {
        window.location.href = signInUrl;
        return true;
      }
      if (!isAdmin) {
        if (replace) {
          window.history.replaceState({}, "", profilePath);
          currentRoute = profilePath;
        } else {
          window.__navigate(profilePath);
        }
        return true;
      }
    }

    // /admin entry -> default to the Users sub-module
    if (relPath === "/admin" || relPath === "/admin/") {
      if (!isLoggedIn) {
        window.location.href = signInUrl;
        return true;
      }
      if (!isAdmin) {
        if (replace) {
          window.history.replaceState({}, "", profilePath);
          currentRoute = profilePath;
        } else {
          window.__navigate(profilePath);
        }
        return true;
      }
      const adminUsersUrl = routePath + "/admin/users";
      if (replace) {
        window.history.replaceState({}, "", adminUsersUrl);
        currentRoute = adminUsersUrl;
      } else {
        window.__navigate(adminUsersUrl);
      }
      return true;
    }

    return false;
  }

  // Handle runtime redirects (after popstate / nav clicks) in an effect so
  // that reactive state mutations never occur during the render phase.
  $effect(() => {
    if (!authChecked) return;
    applyRedirect(false);
  });

  // check sign-in status via /iam/v2/auth/session
  (async () => {
    try {
      const resp = await fetch(routePath + "/v2/auth/session", {
        credentials: "same-origin",
      });
      const data = await resp.json();
      isLoggedIn = data.status?.code === "200" && !!data.auth_claims;
      username = data.auth_claims?.sub || "";
      userRoles = data.identity_token?.roles || [];
    } catch {
      isLoggedIn = false;
      username = "";
      userRoles = [];
    }

    // Pre-correct the URL before the first render so the target page
    // component mounts only once (avoids duplicate data fetches).
    applyRedirect(true);

    authChecked = true;
  })();
</script>

{#if !authChecked}
  <div
    class="bg-auth min-vh-100 d-flex flex-column justify-content-center align-items-center"
  >
    <div class="text-center text-white">
      <div
        class="spinner-border"
        style="width:2rem;height:2rem"
        role="status"
      ></div>
      <p class="mt-2">Loading...</p>
    </div>
  </div>
{:else}
  {#key currentRoute}
    {@const Component = resolveComponent()}
    <Component />
  {/key}
{/if}
