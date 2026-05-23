<script>
  import SignIn from "./pages/SignIn.svelte";
  import SignUp from "./pages/SignUp.svelte";
  import ForgotPassword from "./pages/ForgotPassword.svelte";
  import ResetPassword from "./pages/ResetPassword.svelte";
  import UserProfile from "./pages/UserProfile.svelte";
  import AccessKey from "./pages/AccessKey.svelte";
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

  const signInPath = "/service/sign-in";
  const signUpPath = "/service/sign-up";
  const forgotPasswordPath = "/service/forgot-password";
  const resetPasswordPath = "/service/reset-password";

  /** @type {Record<string, typeof SignIn>} */
  const publicRoutes = {
    [signInPath]: SignIn,
    [signUpPath]: SignUp,
    [forgotPasswordPath]: ForgotPassword,
    [resetPasswordPath]: ResetPassword,
  };

  /** @type {Record<string, typeof UserProfile>} */
  const protectedRoutes = {
    "": UserProfile,
    "/": UserProfile,
    "/service/profile": UserProfile,
    "/service/access-keys": AccessKey,
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

    if (publicRoutes[relPath]) {
      if (isLoggedIn) {
        // logged-in user on a public route, redirect to profile
        window.__navigate(routePath + "/");
      }
      return publicRoutes[relPath];
    }

    if (protectedRoutes[relPath]) {
      if (!isLoggedIn) {
        // not logged in, redirect to sign-in
        window.location.href = routePath + signInPath;
        return SignIn;
      }
      return protectedRoutes[relPath];
    }

    // default: show sign-in or profile based on auth state
    return isLoggedIn ? UserProfile : SignIn;
  }

  window.__navigate = (path) => {
    window.history.pushState({}, "", path);
    currentRoute = path;
  };

  // check sign-in status via /iam/v2/service/user-session
  (async () => {
    try {
      const resp = await fetch(routePath + "/v2/service/user-session", {
        credentials: "same-origin",
      });
      const data = await resp.json();
      isLoggedIn = data.status?.code === "200" && !!data.access_token;
    } catch {
      isLoggedIn = false;
    }

    const relPath = getRelPath();
    if (isLoggedIn && publicRoutes[relPath]) {
      window.location.href = routePath + "/";
    } else if (!isLoggedIn && protectedRoutes[relPath]) {
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
