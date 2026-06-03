<script>
  import { fetchSession } from "../../lib/session.svelte.js";
  import UserLayout from "./Layout.svelte";
  import {
    fetchAccessKeys,
    fetchAccessKey,
    setAccessKey,
    deleteAccessKey,
    bindAccessKeyScope,
    unbindAccessKeyScope,
  } from "../../lib/api.js";

  fetchSession();

  let alertMsg = $state("");
  let alertType = $state("info");
  /** @type {any[]} */
  let keys = $state([]);
  let loading = $state(true);

  // Set modal state
  let showSetModal = $state(false);
  let setModalTitle = $state("New Access Key");
  let setFormId = $state("");
  let setFormDesc = $state("");
  let setFormState = $state("active");
  let setSaving = $state(false);
  let setAlert = $state("");

  // Delete confirm modal
  let showDeleteModal = $state(false);
  let deleteId = $state("");
  let deleteSaving = $state(false);
  let deleteAlert = $state("");

  // Bind modal
  let showBindModal = $state(false);
  let bindId = $state("");
  let bindScopeContent = $state("");
  let bindSaving = $state(false);
  let bindAlert = $state("");

  // Unbind confirm modal
  let showUnbindModal = $state(false);
  let unbindId = $state("");
  let unbindScopeName = $state("");
  let unbindSaving = $state(false);
  let unbindAlert = $state("");

  // Secret display modal
  let showSecretModal = $state(false);
  let newSecret = $state("");
  let secretCopied = $state(false);

  /** @param {unknown} err @returns {string} */
  function toErrMsg(err) {
    return err instanceof Error ? err.message : String(err);
  }

  async function loadKeys() {
    loading = true;
    try {
      keys = await fetchAccessKeys();
    } catch (err) {
      showAlert("danger", toErrMsg(err));
    } finally {
      loading = false;
    }
  }

  // init
  loadKeys();

  /** @param {string} type @param {string} msg */
  function showAlert(type, msg) {
    alertType = type;
    alertMsg = msg;
    window.scrollTo({ top: 0, behavior: "smooth" });
    setTimeout(() => (alertMsg = ""), 5000);
  }

  /** @param {string} state */
  function stateBadgeClass(state) {
    return state === "active" ? "bg-success" : "bg-secondary";
  }

  // -- Set (Create/Edit) --
  function openCreateModal() {
    setModalTitle = "New Access Key";
    setFormId = "";
    setFormDesc = "";
    setFormState = "active";
    setAlert = "";
    showSetModal = true;
  }

  /** @param {string} id */
  async function openEditModal(id) {
    setModalTitle = "Access Key Settings";
    setFormId = id;
    setFormDesc = "";
    setFormState = "active";
    setAlert = "";
    showSetModal = true;
    try {
      const item = await fetchAccessKey(id);
      setFormDesc = item.description || "";
      setFormState = item.state || "active";
    } catch (err) {
      setAlert = toErrMsg(err);
    }
  }

  async function onSetSave() {
    setSaving = true;
    setAlert = "";
    try {
      const result = await setAccessKey({
        id: setFormId,
        description: setFormDesc,
        state: setFormState,
      });
      showSetModal = false;

      // Show the secret only on creation (first time) in a dedicated modal
      if (result.item?.secret) {
        newSecret = result.item.secret;
        secretCopied = false;
        showSecretModal = true;
      } else {
        showAlert("success", setFormId ? "Access key updated" : "Access key created");
      }
      await loadKeys();
    } catch (err) {
      setAlert = toErrMsg(err);
    } finally {
      setSaving = false;
    }
  }

  async function copySecret() {
    try {
      await navigator.clipboard.writeText(newSecret);
      secretCopied = true;
      setTimeout(() => (secretCopied = false), 2000);
    } catch {
      // fallback for non-secure contexts
      const ta = document.createElement("textarea");
      ta.value = newSecret;
      ta.style.position = "fixed";
      ta.style.opacity = "0";
      document.body.appendChild(ta);
      ta.select();
      document.execCommand("copy");
      document.body.removeChild(ta);
      secretCopied = true;
      setTimeout(() => (secretCopied = false), 2000);
    }
  }

  // -- Delete --
  /** @param {string} id */
  function openDeleteModal(id) {
    deleteId = id;
    deleteAlert = "";
    showDeleteModal = true;
  }

  async function onDeleteConfirm() {
    deleteSaving = true;
    deleteAlert = "";
    try {
      await deleteAccessKey(deleteId);
      showDeleteModal = false;
      showAlert("success", "Access key deleted");
      await loadKeys();
    } catch (err) {
      deleteAlert = toErrMsg(err);
    } finally {
      deleteSaving = false;
    }
  }

  // -- Bind --
  /** @param {string} id */
  function openBindModal(id) {
    bindId = id;
    bindScopeContent = "";
    bindAlert = "";
    showBindModal = true;
  }

  async function onBindSave() {
    if (!bindScopeContent.includes("=")) {
      bindAlert = 'Scope format: "name=value"';
      return;
    }
    bindSaving = true;
    bindAlert = "";
    try {
      await bindAccessKeyScope(bindId, bindScopeContent);
      showBindModal = false;
      showAlert("success", "Scope added");
      await loadKeys();
    } catch (err) {
      bindAlert = toErrMsg(err);
    } finally {
      bindSaving = false;
    }
  }

  // -- Unbind --
  /** @param {string} id @param {string} scopeName */
  function openUnbindModal(id, scopeName) {
    unbindId = id;
    unbindScopeName = scopeName;
    unbindAlert = "";
    showUnbindModal = true;
  }

  async function onUnbindConfirm() {
    unbindSaving = true;
    unbindAlert = "";
    try {
      await unbindAccessKeyScope(unbindId, unbindScopeName);
      showUnbindModal = false;
      showAlert("success", "Scope removed");
      await loadKeys();
    } catch (err) {
      unbindAlert = toErrMsg(err);
    } finally {
      unbindSaving = false;
    }
  }

  // Global ESC key handler for all modals
  /** @param {KeyboardEvent} e */
  function onGlobalKeydown(e) {
    if (e.key !== "Escape") return;
    if (showSecretModal) showSecretModal = false;
    else if (showSetModal) showSetModal = false;
    else if (showDeleteModal) showDeleteModal = false;
    else if (showBindModal) showBindModal = false;
    else if (showUnbindModal) showUnbindModal = false;
  }

  /** @param {Record<string, any>} ak @param {string} field @returns {string} */
  function akField(ak, field) {
    return ak[field] ?? ak[field.charAt(0).toUpperCase() + field.slice(1)] ?? "";
  }
</script>

<svelte:window onkeydown={onGlobalKeydown} />

<UserLayout bind:alertMsg {alertType}>
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
    <div class="card shadow-sm border-0">
      <div
        class="card-header bg-transparent border-bottom d-flex justify-content-between align-items-center py-3"
      >
        <h5 class="mb-0">Access Keys</h5>
        <button class="btn btn-primary btn-sm" onclick={openCreateModal}>
          New Access Key
        </button>
      </div>
      <div class="card-body p-3">
        {#if keys.length === 0}
          <div class="text-center py-5 text-muted">
            <p class="mb-2">No access keys found</p>
          </div>
        {:else}
          <div class="table-responsive">
            <table class="table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Description</th>
                  <th>Status</th>
                  <th>Scopes</th>
                  <th class="text-end">Actions</th>
                </tr>
              </thead>
              <tbody>
                {#each keys as ak}
                  {@const akId = akField(ak, "id")}
                  {@const akScopes = ak.scopes || ak.Scopes || []}
                  <tr>
                    <td>
                      <code class="user-select-all" style="font-size:0.85em"
                        >{akId}</code
                      >
                    </td>
                    <td>{akField(ak, "description")}</td>
                    <td>
                      <span
                        class="badge {stateBadgeClass(akField(ak, 'state'))}"
                      >
                        {akField(ak, "state") === "active" ? "Active" : "Disabled"}
                      </span>
                    </td>
                    <td>
                      {#each akScopes as scope}
                        {@const scopeName = scope.includes("=") ? scope.split("=")[0] : scope}
                        <span
                          class="badge bg-success me-1 mb-1 d-inline-flex align-items-center"
                        >
                          {scope}
                          <button
                            type="button"
                            class="btn-close btn-close-white ms-1"
                            style="font-size:0.65em"
                            aria-label="Remove scope"
                            onclick={() => openUnbindModal(akId, scopeName)}
                          ></button>
                        </span>
                      {/each}
                      {#if !akScopes.length}
                        <span class="text-muted small">-</span>
                      {/if}
                    </td>
                    <td class="text-end text-nowrap">
                      <button
                        class="btn btn-outline-primary btn-sm"
                        onclick={() => openBindModal(akId)}
                        title="Add Scope"
                      >
                        Add Scope
                      </button>
                      <button
                        class="btn btn-outline-primary btn-sm"
                        onclick={() => openEditModal(akId)}
                        title="Settings"
                      >
                        Setting
                      </button>
                      <button
                        class="btn btn-outline-danger btn-sm"
                        onclick={() => openDeleteModal(akId)}
                        title="Delete"
                      >
                        Delete
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
  {/if}
</UserLayout>

<!-- Set Modal (Create/Edit) -->
{#if showSetModal}
  <div class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,0.5)" role="dialog">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{setModalTitle}</h5>
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
          <div class="mb-3">
            <label class="form-label" for="setDescription">Description</label>
            <input
              id="setDescription"
              type="text"
              class="form-control"
              bind:value={setFormDesc}
              placeholder="Access key description"
            />
          </div>
          <div class="mb-3">
            <span class="form-label">Status</span>
            <div>
              <div class="form-check form-check-inline">
                <input
                  class="form-check-input"
                  type="radio"
                  name="akState"
                  id="akStateActive"
                  value="active"
                  bind:group={setFormState}
                />
                <label class="form-check-label" for="akStateActive">Active</label>
              </div>
              <div class="form-check form-check-inline">
                <input
                  class="form-check-input"
                  type="radio"
                  name="akState"
                  id="akStateDisable"
                  value="disabled"
                  bind:group={setFormState}
                />
                <label class="form-check-label" for="akStateDisable">Disabled</label>
              </div>
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

<!-- Delete Confirm Modal -->
{#if showDeleteModal}
  <div class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,0.5)" role="dialog">
    <div class="modal-dialog modal-dialog-centered modal-sm">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Delete</h5>
          <button
            type="button"
            class="btn-close"
            aria-label="Close"
            onclick={() => (showDeleteModal = false)}
          ></button>
        </div>
        <div class="modal-body">
          {#if deleteAlert}
            <div class="alert alert-danger py-2 small">{deleteAlert}</div>
          {/if}
          <div class="alert alert-danger mb-0">
            Are you sure you want to delete this access key?
          </div>
        </div>
        <div class="modal-footer">
          <button
            class="btn btn-secondary btn-sm"
            onclick={() => (showDeleteModal = false)}>Cancel</button
          >
          <button
            class="btn btn-danger btn-sm"
            onclick={onDeleteConfirm}
            disabled={deleteSaving}
          >
            {#if deleteSaving}
              <span class="spinner-border spinner-border-sm me-1" role="status"
              ></span>
            {:else}
              Confirm Delete
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<!-- Bind Scope Modal -->
{#if showBindModal}
  <div class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,0.5)" role="dialog">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Add Scope</h5>
          <button
            type="button"
            class="btn-close"
            aria-label="Close"
            onclick={() => (showBindModal = false)}
          ></button>
        </div>
        <div class="modal-body">
          {#if bindAlert}
            <div class="alert alert-danger py-2 small">{bindAlert}</div>
          {/if}
          <div class="mb-3">
            <label class="form-label" for="bindScope">Scope</label>
            <input
              id="bindScope"
              type="text"
              class="form-control"
              bind:value={bindScopeContent}
              placeholder="e.g. app=instance-id"
            />
            <div class="form-text">Format: name = value</div>
          </div>
        </div>
        <div class="modal-footer">
          <button
            class="btn btn-secondary btn-sm"
            onclick={() => (showBindModal = false)}>Cancel</button
          >
          <button
            class="btn btn-primary btn-sm"
            onclick={onBindSave}
            disabled={bindSaving}
          >
            {#if bindSaving}
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

<!-- Secret Display Modal -->
{#if showSecretModal}
  <div class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,0.5)" role="dialog">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Access Key Created</h5>
          <button
            type="button"
            class="btn-close"
            aria-label="Close"
            onclick={() => (showSecretModal = false)}
          ></button>
        </div>
        <div class="modal-body">
          <div class="alert alert-warning py-2 small">
            Please save your Secret Key now. This is the only time it will be shown. You will not be able to view it again after closing this dialog.
          </div>
          <div class="mb-2">
            <label class="form-label fw-semibold" for="newSecretInput">Secret Key</label>
            <div class="input-group">
              <input
                id="newSecretInput"
                type="text"
                class="form-control font-monospace user-select-all"
                readonly
                value={newSecret}
              />
              <button
                class="btn btn-outline-secondary"
                type="button"
                onclick={copySecret}
              >
                {#if secretCopied}
                  Copied!
                {:else}
                  Copy
                {/if}
              </button>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button
            class="btn btn-primary btn-sm"
            onclick={() => (showSecretModal = false)}
          >Close</button>
        </div>
      </div>
    </div>
  </div>
{/if}

<!-- Unbind Confirm Modal -->
{#if showUnbindModal}
  <div class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,0.5)" role="dialog">
    <div class="modal-dialog modal-dialog-centered modal-sm">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Remove Scope</h5>
          <button
            type="button"
            class="btn-close"
            aria-label="Close"
            onclick={() => (showUnbindModal = false)}
          ></button>
        </div>
        <div class="modal-body">
          {#if unbindAlert}
            <div class="alert alert-danger py-2 small">{unbindAlert}</div>
          {/if}
          <div class="alert alert-warning mb-0">
            Remove scope <strong>{unbindScopeName}</strong>?
          </div>
        </div>
        <div class="modal-footer">
          <button
            class="btn btn-secondary btn-sm"
            onclick={() => (showUnbindModal = false)}>Cancel</button
          >
          <button
            class="btn btn-danger btn-sm"
            onclick={onUnbindConfirm}
            disabled={unbindSaving}
          >
            {#if unbindSaving}
              <span class="spinner-border spinner-border-sm me-1" role="status"
              ></span>
            {:else}
              Confirm Remove
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  tbody tr:last-child td {
    border-bottom: none;
  }
</style>