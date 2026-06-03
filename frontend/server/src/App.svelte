<script>
  import Auth_SignIn from "./pages/auth/SignIn.svelte";
  import Auth_SignUp from "./pages/auth/SignUp.svelte";
  import Auth_Password_Forgot from "./pages/auth/password/Forgot.svelte";
  import Auth_Password_Reset from "./pages/auth/password/Reset.svelte";

  import User_Profile from "./pages/user/Profile.svelte";
  import User_Keys from "./pages/user/AccessKey.svelte";
  import User_Apps from "./pages/user/Apps.svelte";

  import { routePath } from "./lib/config.js";

  let currentRoute = $state(window.location.pathname);
  let authChecked = $state(false);
  let isLoggedIn = $state(false);

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
    "": User_Profile,
    "/": User_Profile,
    "/user/profile": User_Profile,
    "/user/keys": User_Keys,
    "/user/apps": User_Apps,
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

  function getComponent() {
    const relPath = getRelPath();

    if (authRoutes[relPath]) {
      if (isLoggedIn) {
        // logged-in user on a public route, redirect to profile
        window.__navigate(routePath + "/");
      }
      return authRoutes[relPath];
    }

    if (userRoutes[relPath]) {
      if (!isLoggedIn) {
        // not logged in, redirect to sign-in
        window.location.href = routePath + signInPath;
        return Auth_SignIn;
      }
      return userRoutes[relPath];
    }

    // default: show sign-in or profile based on auth state
    return isLoggedIn ? User_Profile : Auth_SignIn;
  }

  window.__navigate = (path) => {
    window.history.pushState({}, "", path);
    currentRoute = path;
  };

  // check sign-in status via /iam/v2/auth/session
  (async () => {
    try {
      const resp = await fetch(routePath + "/v2/auth/session", {
        credentials: "same-origin",
      });
      const data = await resp.json();
      isLoggedIn = data.status?.code === "200" && !!data.auth_claims;
    } catch {
      isLoggedIn = false;
    }

    const relPath = getRelPath();
    if (isLoggedIn && authRoutes[relPath]) {
      window.location.href = routePath + "/";
    } else if (!isLoggedIn && userRoutes[relPath]) {
      window.location.href = routePath + signInPath;
    }

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
    {@const Component = getComponent()}
    <Component />
  {/key}
{/if}
