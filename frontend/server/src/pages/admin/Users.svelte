<script>
  import { fetchSession } from "../../lib/session.svelte.js";
  import AdminLayout from "./Layout.svelte";
  import {
    fetchUsers,
    fetchUserEntry,
    adminSaveUser,
    fetchAdminRoles,
  } from "../../lib/api.js";

  fetchSession();

  let alertMsg = $state("");
  let alertType = $state("info");
  /** @type {any[]} */
  let users = $state([]);
  /** @type {any[]} */
  let roles = $state([]);
  let loading = $state(true);
  let qryText = $state("");

  // Set modal state (create/edit)
  let showSetModal = $state(false);
  let setIsEdit = $state(false);
  let setFormName = $state("");
  let setFormEmail = $state("");
  let setFormPassword = $state("");
  let setFormDisplayName = $state("");
  let setFormBirthday = $state("");
  let setFormAbout = $state("");
  /** @type {string[]} */
  let setFormRoles = $state([]);
  let setFormStatus = $state("1");
  let setSaving = $state(false);
  let setAlert = $state("");

  // The sysadmin super-user is protected: its roles and status cannot be
  // modified (the backend lockout guard enforces this; the UI reflects it).
  let isProtectedUser = $derived(setIsEdit && setFormName === "sysadmin");

  const roleLabel = {
    sa: "Sysadmin",
    user: "User",
    dev: "Developer",
    guest: "Guest",
  };

  /** @param {unknown} err @returns {string} */
  function toErrMsg(err) {
    return err instanceof Error ? err.message : String(err);
  }

  /** @param {number} unixSec @returns {string} */
  function fmtDate(unixSec) {
    if (!unixSec) return "-";
    const d = new Date(unixSec * 1000);
    const p = (n) => String(n).padStart(2, "0");
    return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}`;
  }

  async function loadRoles() {
    try {
      roles = await fetchAdminRoles();
    } catch {
      roles = [];
    }
  }

  /** @param {string} [q] */
  async function loadUsers(q) {
    loading = true;
    try {
      users = await fetchUsers(q || "");
    } catch (err) {
      showAlert("danger", toErrMsg(err));
    } finally {
      loading = false;
    }
  }

  // init
  loadRoles();
  loadUsers();

  /** @param {string} type @param {string} msg */
  function showAlert(type, msg) {
    alertType = type;
    alertMsg = msg;
    window.scrollTo({ top: 0, behavior: "smooth" });
    setTimeout(() => (alertMsg = ""), 5000);
  }

  /** @param {number} status @returns {string} */
  function statusBadgeClass(status) {
    return status === 2 ? "bg-danger" : "bg-success";
  }

  /** @param {number} status @returns {string} */
  function statusLabel(status) {
    return status === 2 ? "Banned" : "Active";
  }

  // -- Set (Create/Edit) --
  function openCreateModal() {
    setIsEdit = false;
    setFormName = "";
    setFormEmail = "";
    setFormPassword = "";
    setFormDisplayName = "";
    setFormBirthday = "";
    setFormAbout = "";
    setFormRoles = [];
    setFormStatus = "1";
    setAlert = "";
    showSetModal = true;
  }

  /** @param {string} name */
  async function openEditModal(name) {
    setIsEdit = true;
    setFormName = name;
    setFormEmail = "";
    setFormPassword = "";
    setFormDisplayName = "";
    setFormBirthday = "";
    setFormAbout = "";
    setFormRoles = [];
    setFormStatus = "1";
    setAlert = "";
    showSetModal = true;
    try {
      const item = await fetchUserEntry(name);
      setFormEmail = item.email || "";
      setFormDisplayName = item.display_name || "";
      setFormBirthday = item.birthday || "";
      setFormAbout = item.about || "";
      setFormRoles = Array.isArray(item.roles) ? item.roles.slice() : [];
      setFormStatus = String(item.status || 1);
    } catch (err) {
      setAlert = toErrMsg(err);
    }
  }

  /** @param {string} role @returns {boolean} */
  function hasRole(role) {
    return setFormRoles.includes(role);
  }

  /** @param {string} role @param {boolean} checked */
  function toggleRole(role, checked) {
    if (checked) {
      if (!setFormRoles.includes(role)) setFormRoles = [...setFormRoles, role];
    } else {
      setFormRoles = setFormRoles.filter((r) => r !== role);
    }
  }

  async function onSetSave() {
    setAlert = "";

    if (!setFormName) {
      setAlert = "Username is required";
      return;
    }
    if (!setIsEdit) {
      if (setFormPassword.length < 8 || setFormPassword.length > 30) {
        setAlert = "Password must be between 8 and 30 characters long";
        return;
      }
    } else if (
      setFormPassword &&
      (setFormPassword.length < 8 || setFormPassword.length > 30)
    ) {
      setAlert = "Password must be between 8 and 30 characters long";
      return;
    }

    setSaving = true;
    try {
      await adminSaveUser({
        name: setFormName,
        email: setFormEmail,
        password: setFormPassword,
        display_name: setFormDisplayName,
        roles: setFormRoles,
        status: Number(setFormStatus),
        birthday: setFormBirthday,
        about: setFormAbout,
      });
      showSetModal = false;
      showAlert("success", setIsEdit ? "User updated" : "User created");
      await loadUsers(qryText);
    } catch (err) {
      setAlert = toErrMsg(err);
    } finally {
      setSaving = false;
    }
  }

  /** @param {KeyboardEvent} e */
  function onSearchKeydown(e) {
    if (e.key === "Enter") {
      e.preventDefault();
      loadUsers(qryText);
    }
  }

  /** @param {KeyboardEvent} e */
  function onGlobalKeydown(e) {
    if (e.key === "Escape" && showSetModal) showSetModal = false;
  }
</script>

<svelte:window onkeydown={onGlobalKeydown} />

<AdminLayout bind:alertMsg {alertType}>
  <div class="card shadow-sm border-0">
    <div
      class="card-header bg-transparent border-bottom d-flex justify-content-between align-items-center py-3"
    >
      <h5 class="mb-0">Users</h5>
      <div class="d-flex align-items-center gap-2">
        <input
          type="text"
          class="form-control form-control-sm"
          style="width:220px"
          placeholder="Press Enter to Search"
          bind:value={qryText}
          onkeydown={onSearchKeydown}
        />
        <button class="btn btn-primary btn-sm" onclick={openCreateModal}>
          New User
        </button>
      </div>
    </div>
    <div class="card-body p-3">
      {#if loading}
        <div class="text-center py-5">
          <div
            class="spinner-border"
            style="width:2rem;height:2rem"
            role="status"
          ></div>
          <p class="mt-2 text-muted">Loading...</p>
        </div>
      {:else if users.length === 0}
        <div class="text-center py-5 text-muted">
          <p class="mb-2">No users found</p>
        </div>
      {:else}
        <div class="table-responsive">
          <table class="table table-last-row-borderless">
            <thead>
              <tr>
                <th>Username</th>
                <th>Display Name</th>
                <th>Email</th>
                <th>Roles</th>
                <th>Status</th>
                <th>Updated</th>
                <th class="text-end">Actions</th>
              </tr>
            </thead>
            <tbody>
              {#each users as u}
                <tr>
                  <td><code class="user-select-all" style="font-size:0.85em"
                      >{u.name}</code
                    ></td>
                  <td>{u.display_name || "-"}</td>
                  <td>{u.email || "-"}</td>
                  <td>
                    {#each (u.roles || []) as r}
                      <span class="badge bg-secondary me-1 mb-1">
                        {roleLabel[r] || r}
                      </span>
                    {/each}
                    {#if !(u.roles || []).length}
                      <span class="text-muted small">-</span>
                    {/if}
                  </td>
                  <td>
                    <span class="badge {statusBadgeClass(u.status)}">
                      {statusLabel(u.status)}
                    </span>
                  </td>
                  <td>{fmtDate(u.updated)}</td>
                  <td class="text-end text-nowrap">
                    <button
                      class="btn btn-outline-primary btn-sm"
                      onclick={() => openEditModal(u.name)}
                      title="Edit"
                    >
                      Edit
                    </button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  </div>
</AdminLayout>

<!-- Set Modal (Create/Edit) -->
{#if showSetModal}
  <div class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,0.5)" role="dialog">
    <div class="modal-dialog modal-dialog-centered modal-dialog-900">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{setIsEdit ? "Edit User" : "New User"}</h5>
          <button
            type="button"
            class="btn-close"
            aria-label="Close"
            onclick={() => (showSetModal = false)}
          ></button>
        </div>
        <div class="modal-body">
          {#if setAlert}
            <div class="alert alert-danger py-2 small">{setAlert}</div>
          {/if}

          <div class="row mb-2 align-items-center">
            <label class="col-sm-4 col-form-label" for="setName">Username</label>
            <div class="col-sm-8">
              <input
                id="setName"
                type="text"
                class="form-control"
                bind:value={setFormName}
                placeholder="3-30 chars, a-z, 0-9, -"
                disabled={setIsEdit}
              />
            </div>
          </div>

          <div class="row mb-2 align-items-center">
            <label class="col-sm-4 col-form-label" for="setEmail">Email</label>
            <div class="col-sm-8">
              <input
                id="setEmail"
                type="email"
                class="form-control"
                bind:value={setFormEmail}
                placeholder="user@example.com"
              />
            </div>
          </div>

          <div class="row mb-2 align-items-center">
            <label class="col-sm-4 col-form-label" for="setPassword">Password</label>
            <div class="col-sm-8">
              <input
                id="setPassword"
                type="password"
                class="form-control"
                bind:value={setFormPassword}
                placeholder={setIsEdit
                  ? "Leave blank to keep current"
                  : "8-30 characters"}
                autocomplete="new-password"
              />
            </div>
          </div>

          <div class="row mb-2 align-items-center">
            <label class="col-sm-4 col-form-label" for="setDisplayName"
              >Display Name</label
            >
            <div class="col-sm-8">
              <input
                id="setDisplayName"
                type="text"
                class="form-control"
                bind:value={setFormDisplayName}
              />
            </div>
          </div>

          <div class="row mb-2">
            <span class="col-sm-4 col-form-label">
              Roles
              {#if isProtectedUser}
                <span class="text-muted small">(locked)</span>
              {/if}
            </span>
            <div class="col-sm-8">
              {#if roles.length === 0}
                <span class="text-muted small">No roles available</span>
              {/if}
              {#each roles as role}
                <div class="form-check form-check-inline">
                  <input
                    class="form-check-input"
                    type="checkbox"
                    id="role_{role.name}"
                    value={role.name}
                    checked={hasRole(role.name)}
                    disabled={isProtectedUser}
                    onchange={(e) => toggleRole(role.name, e.currentTarget.checked)}
                  />
                  <label class="form-check-label" for="role_{role.name}">
                    {roleLabel[role.name] || role.name}
                  </label>
                </div>
              {/each}
            </div>
          </div>

          <div class="row mb-2 align-items-center">
            <span class="col-sm-4 col-form-label">
              Status
              {#if isProtectedUser}
                <span class="text-muted small">(locked)</span>
              {/if}
            </span>
            <div class="col-sm-8">
              <div class="form-check form-check-inline">
                <input
                  class="form-check-input"
                  type="radio"
                  name="userStatus"
                  id="userStatusActive"
                  value="1"
                  bind:group={setFormStatus}
                  disabled={isProtectedUser}
                />
                <label class="form-check-label" for="userStatusActive">Active</label>
              </div>
              <div class="form-check form-check-inline">
                <input
                  class="form-check-input"
                  type="radio"
                  name="userStatus"
                  id="userStatusBanned"
                  value="2"
                  bind:group={setFormStatus}
                  disabled={isProtectedUser}
                />
                <label class="form-check-label" for="userStatusBanned">Banned</label>
              </div>
            </div>
          </div>

          <div class="row mb-2 align-items-center">
            <label class="col-sm-4 col-form-label" for="setBirthday">Birthday</label>
            <div class="col-sm-8">
              <input
                id="setBirthday"
                type="text"
                class="form-control"
                bind:value={setFormBirthday}
                placeholder="YYYY-MM-DD"
              />
            </div>
          </div>

          <div class="row mb-2">
            <label class="col-sm-4 col-form-label" for="setAbout">About</label>
            <div class="col-sm-8">
              <textarea
                id="setAbout"
                class="form-control"
                rows="2"
                bind:value={setFormAbout}
              ></textarea>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button
            class="btn btn-secondary btn-sm"
            onclick={() => (showSetModal = false)}>Cancel</button
          >
          <button
            class="btn btn-primary btn-sm"
            onclick={onSetSave}
            disabled={setSaving}
          >
            {#if setSaving}
              <span class="spinner-border spinner-border-sm me-1" role="status"
              ></span>
              Saving...
            {:else}
              Save
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
