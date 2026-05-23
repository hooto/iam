<script>
  import { fetchSession, getSession } from "../lib/session.svelte.js";
  import UserLayout from "../layouts/UserLayout.svelte";
  import {
    fetchProfile,
    updateProfile,
    changePassword,
    changeEmail,
    uploadPhoto,
  } from "../lib/api.js";

  const session = getSession();
  fetchSession();

  let alertMsg = $state("");
  let alertType = $state("info");
  /** @type {{ display_name?: string; birthday?: string; about?: string; email?: string; photo?: string } | null} */
  let profile = $state(null);
  let loading = $state(true);

  // profile form
  let displayName = $state("");
  let birthday = $state("");
  let about = $state("");
  let profileSaving = $state(false);

  // password form
  let currentPassword = $state("");
  let newPassword = $state("");
  let confirmPassword = $state("");
  let passwordSaving = $state(false);

  // email form
  let email = $state("");
  let emailPassword = $state("");
  let emailSaving = $state(false);

  // photo
  let photoSaving = $state(false);

  async function loadProfile() {
    loading = true;
    try {
      profile = await fetchProfile();
      if (profile) {
        displayName = profile.display_name || "";
        birthday = profile.birthday || "";
        about = profile.about || "";
        email = profile.email || "";
      }
    } catch (/** @type {any} */ err) {
      showAlert("danger", err.message);
    } finally {
      loading = false;
    }
  }

  // fetch profile on mount
  loadProfile();

  /** @param {string} type @param {string} msg */
  function showAlert(type, msg) {
    alertType = type;
    alertMsg = msg;
    window.scrollTo({ top: 0, behavior: "smooth" });
    setTimeout(() => {
      alertMsg = "";
    }, 5000);
  }

  async function onProfileSave() {
    profileSaving = true;
    try {
      await updateProfile({
        display_name: displayName,
        birthday: birthday,
        about: about,
      });
      showAlert("success", "Profile updated successfully");
      await loadProfile();
    } catch (err) {
      showAlert("danger", err instanceof Error ? err.message : String(err));
    } finally {
      profileSaving = false;
    }
  }

  async function onPasswordSave() {
    if (newPassword !== confirmPassword) {
      showAlert("danger", "Passwords do not match");
      return;
    }
    if (newPassword.length < 8) {
      showAlert("danger", "Password must be at least 8 characters long");
      return;
    }
    passwordSaving = true;
    try {
      await changePassword({
        current_password: currentPassword,
        new_password: newPassword,
      });
      currentPassword = "";
      newPassword = "";
      confirmPassword = "";
      showAlert("success", "Password changed successfully");
    } catch (err) {
      showAlert("danger", err instanceof Error ? err.message : String(err));
    } finally {
      passwordSaving = false;
    }
  }

  async function onEmailSave() {
    emailSaving = true;
    try {
      await changeEmail({
        email: email,
        auth: emailPassword,
      });
      emailPassword = "";
      showAlert("success", "Email changed successfully");
      await loadProfile();
    } catch (err) {
      showAlert("danger", err instanceof Error ? err.message : String(err));
    } finally {
      emailSaving = false;
    }
  }

  /** @param {Event & { currentTarget: EventTarget & HTMLInputElement }} e */
  async function onPhotoChange(e) {
    const file = e.currentTarget.files?.[0];
    if (!file) return;

    if (file.size > 2 * 1024 * 1024) {
      showAlert("danger", "The file is too large (max 2MB)");
      return;
    }

    if (!file.type.match(/^image\/(jpeg|png|gif)$/)) {
      showAlert("danger", "Please upload a JPG, GIF, or PNG file");
      return;
    }

    photoSaving = true;
    try {
      const dataUrl = await readFileAsDataUrl(file);
      await uploadPhoto(dataUrl);
      showAlert("success", "Photo updated successfully");
      await loadProfile();
    } catch (err) {
      showAlert("danger", err instanceof Error ? err.message : String(err));
    } finally {
      photoSaving = false;
      e.currentTarget.value = "";
    }
  }

  /** @param {File} file */
  function readFileAsDataUrl(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(reader.result);
      reader.onerror = () => reject(new Error("Failed to read file"));
      reader.readAsDataURL(file);
    });
  }

  const defaultPhoto =
    "data:image/svg+xml," +
    encodeURIComponent(
      '<svg xmlns="http://www.w3.org/2000/svg" width="96" height="96" viewBox="0 0 16 16">' +
        '<path fill="#6c757d" d="M11 6a3 3 0 1 1-6 0 3 3 0 0 1 6 0"/>' +
        '<path fill="#6c757d" d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8m8-7a7 7 0 0 0-5.468 11.37C3.242 11.226 4.805 10 8 10s4.757 1.225 5.468 2.37A7 7 0 0 0 8 1"/>' +
        "</svg>",
    );
</script>

<UserLayout contentClass="" bind:alertMsg {alertType}>
  {#if loading}
    <div class="text-center py-5">
      <div
        class="spinner-border"
        style="width:2rem;height:2rem"
        role="status"
      ></div>
      <p class="mt-2 text-muted">Loading...</p>
    </div>
  {:else}
    <div class="row g-4">
      <!-- Left Column: Personal Info Card -->
      <div class="col-md-4">
        <div class="card shadow-sm border-0">
          <div class="card-body text-center p-4">
            <div class="mb-3">
              <img
                src={profile?.photo || defaultPhoto}
                alt="Avatar"
                class="rounded-circle"
                width="96"
                height="96"
                style="object-fit:cover;"
              />
            </div>
            <h5 class="mb-1">{profile?.display_name || session.username}</h5>
            <p class="text-muted small mb-3">@{session.username}</p>

            <!-- Photo Upload -->
            <div class="mb-3">
              <label
                class="btn btn-outline-primary btn-sm"
                class:disabled={photoSaving}
              >
                {#if photoSaving}
                  <span
                    class="spinner-border spinner-border-sm me-1"
                    role="status"
                  ></span>
                  Uploading...
                {:else}
                  Change Photo
                {/if}
                <input
                  type="file"
                  class="d-none"
                  accept="image/jpeg,image/png,image/gif"
                  onchange={onPhotoChange}
                  disabled={photoSaving}
                />
              </label>
            </div>

            <div>
              <div class="d-flex align-items-center justify-content-center mb-2">
                <span class="text-muted me-2"
                  ><svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="16"
                    height="16"
                    fill="currentColor"
                    class="bi bi-envelope"
                    viewBox="0 0 16 16"
                  >
                    <path
                      d="M0 4a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2zm2-1a1 1 0 0 0-1 1v.217l7 4.2 7-4.2V4a1 1 0 0 0-1-1zm13 2.383-4.708 2.825L15 11.105zm-.034 6.876-5.64-3.471L8 9.583l-1.326-.795-5.64 3.47A1 1 0 0 0 2 13h12a1 1 0 0 0 .966-.741M1 11.105l4.708-2.897L1 5.383z"
                    />
                  </svg></span
                >
                <small>{profile?.email || "Not set"}</small>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Settings Forms -->
      <div class="col-md-8">
        <!-- Profile Edit Card -->
        <div class="card shadow-sm border-0 mb-4">
          <div class="card-header bg-transparent border-bottom-0 pt-3">
            <h5 class="mb-0">Personal Info</h5>
          </div>
          <div class="card-body">
            <form
              onsubmit={(e) => {
                e.preventDefault();
                onProfileSave();
              }}
            >
              <div class="row mb-3">
                <label for="displayName" class="col-sm-4 col-form-label">Display Name</label>
                <div class="col-sm-8">
                  <input
                    type="text"
                    class="form-control"
                    id="displayName"
                    bind:value={displayName}
                    placeholder="Enter display name"
                    maxlength="30"
                    required
                  />
                  <div class="form-text">1-30 characters</div>
                </div>
              </div>

              <div class="row mb-3">
                <label for="birthday" class="col-sm-4 col-form-label">Birthday</label>
                <div class="col-sm-8">
                  <input
                    type="date"
                    class="form-control"
                    id="birthday"
                    bind:value={birthday}
                    max="2099-12-31"
                  />
                  <div class="form-text">Format: YYYY-MM-DD</div>
                </div>
              </div>

              <div class="row mb-3">
                <label for="about" class="col-sm-4 col-form-label">About Me</label>
                <div class="col-sm-8">
                  <textarea
                    class="form-control"
                    id="about"
                    rows="3"
                    bind:value={about}
                    placeholder="Tell us about yourself"
                    required
                  ></textarea>
                </div>
              </div>

              <div class="text-end">
                <button
                  type="submit"
                  class="btn btn-primary"
                  disabled={profileSaving}
                >
                  {#if profileSaving}
                    <span
                      class="spinner-border spinner-border-sm me-1"
                      role="status"
                    ></span>
                    Saving...
                  {:else}
                    Save Profile
                  {/if}
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- Password Change Card -->
        <div class="card shadow-sm border-0 mb-4">
          <div class="card-header bg-transparent border-bottom-0 pt-3">
            <h5 class="mb-0">Change Password</h5>
          </div>
          <div class="card-body">
            <form
              onsubmit={(e) => {
                e.preventDefault();
                onPasswordSave();
              }}
            >
              <div class="row mb-3">
                <label for="currentPassword" class="col-sm-4 col-form-label"
                  >Current Password</label
                >
                <div class="col-sm-8">
                  <input
                    type="password"
                    class="form-control"
                    id="currentPassword"
                    bind:value={currentPassword}
                    placeholder=""
                    autocomplete="current-password"
                    required
                  />
                </div>
              </div>

              <div class="row mb-3">
                <label for="newPassword" class="col-sm-4 col-form-label">New Password</label>
                <div class="col-sm-8">
                  <input
                    type="password"
                    class="form-control"
                    id="newPassword"
                    bind:value={newPassword}
                    placeholder=""
                    autocomplete="new-password"
                    required
                    minlength="8"
                    maxlength="30"
                  />
                  <div class="form-text">8-30 characters</div>
                </div>
              </div>

              <div class="row mb-3">
                <label for="confirmPassword" class="col-sm-4 col-form-label"
                  >Confirm New Password</label
                >
                <div class="col-sm-8">
                  <input
                    type="password"
                    class="form-control"
                    id="confirmPassword"
                    bind:value={confirmPassword}
                    placeholder=""
                    autocomplete="new-password"
                    required
                  />
                </div>
              </div>

              <div class="text-end">
                <button
                  type="submit"
                  class="btn btn-primary"
                  disabled={passwordSaving}
                >
                  {#if passwordSaving}
                    <span
                      class="spinner-border spinner-border-sm me-1"
                      role="status"
                    ></span>
                    Changing...
                  {:else}
                    Change Password
                  {/if}
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- Email Change Card -->
        <div class="card shadow-sm border-0 mb-4">
          <div class="card-header bg-transparent border-bottom-0 pt-3">
            <h5 class="mb-0">Change Email</h5>
          </div>
          <div class="card-body">
            <form
              onsubmit={(e) => {
                e.preventDefault();
                onEmailSave();
              }}
            >
              <div class="row mb-3">
                <label for="email" class="col-sm-4 col-form-label">Email</label>
                <div class="col-sm-8">
                  <input
                    type="email"
                    class="form-control"
                    id="email"
                    bind:value={email}
                    placeholder="Enter email address"
                    required
                  />
                </div>
              </div>

              <div class="row mb-3">
                <label for="emailPassword" class="col-sm-4 col-form-label"
                  >Current Password</label
                >
                <div class="col-sm-8">
                  <input
                    type="password"
                    class="form-control"
                    id="emailPassword"
                    bind:value={emailPassword}
                    placeholder=""
                    autocomplete="current-password"
                    required
                  />
                  <div class="form-text">Required to verify your identity</div>
                </div>
              </div>

              <div class="text-end">
                <button
                  type="submit"
                  class="btn btn-primary"
                  disabled={emailSaving}
                >
                  {#if emailSaving}
                    <span
                      class="spinner-border spinner-border-sm me-1"
                      role="status"
                    ></span>
                    Changing...
                  {:else}
                    Change Email
                  {/if}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  {/if}
</UserLayout>
