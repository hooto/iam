<script>
  import { fetchSession } from "../../lib/session.svelte.js";
  import UserLayout from "./Layout.svelte";
  import {
    fetchApps,
    createApp,
    updateApp,
    deleteApp,
  } from "../../lib/api.js";

  fetchSession();

  let alertMsg = $state("");
  let alertType = $state("info");
  /** @type {any[]} */
  let apps = $state([]);
  let loading = $state(true);

  // Create modal state
  let showCreateModal = $state(false);
  let createName = $state("");
  let createUrl = $state("");
  let createSaving = $state(false);
  let createAlert = $state("");

  // Edit modal state
  let showEditModal = $state(false);
  let editId = $state("");
  let editName = $state("");
  let editUrl = $state("");
  let editSaving = $state(false);
  let editAlert = $state("");

  // Delete modal state
  let showDeleteModal = $state(false);
  let deleteId = $state("");
  let deleteName = $state("");
  let deleteSaving = $state(false);
  let deleteAlert = $state("");

  // Secret reveal
  let revealSecret = $state("");

  /** @param {unknown} err @returns {string} */
  function toErrMsg(err) {
    return err instanceof Error ? err.message : String(err);
  }

  async function loadApps() {
    loading = true;
    try {
      apps = await fetchApps();
    } catch (err) {
      showAlert("danger", toErrMsg(err));
    } finally {
      loading = false;
    }
  }

  loadApps();

  /** @param {string} type @param {string} msg */
  function showAlert(type, msg) {
    alertType = type;
    alertMsg = msg;
    window.scrollTo({ top: 0, behavior: "smooth" });
    setTimeout(() => (alertMsg = ""), 8000);
  }

  /** @param {string} text */
  function copyText(text) {
    navigator.clipboard.writeText(text).then(
      () => showAlert("success", "Copied to clipboard"),
      () => {},
    );
  }

  // -- Create --
  function openCreateModal() {
    createName = "";
    createUrl = "";
    createAlert = "";
    showCreateModal = true;
  }

  async function onCreateSave() {
    createSaving = true;
    createAlert = "";
    try {
      const result = await createApp(createName, createUrl);
      showCreateModal = false;
      if (result.app?.secret_key) {
        revealSecret = result.app.secret_key;
        showAlert(
          "success",
          `App created! App ID: ${result.app.id}`,
        );
      } else {
        showAlert("success", "App created");
      }
      await loadApps();
    } catch (err) {
      createAlert = toErrMsg(err);
    } finally {
      createSaving = false;
    }
  }

  // -- Edit --
  /** @param {any} app */
  function openEditModal(app) {
    editId = app.id;
    editName = app.name;
    editUrl = app.url || "";
    editAlert = "";
    showEditModal = true;
  }

  async function onEditSave() {
    editSaving = true;
    editAlert = "";
    try {
      await updateApp(editId, editName, editUrl);
      showEditModal = false;
      showAlert("success", "App updated");
      await loadApps();
    } catch (err) {
      editAlert = toErrMsg(err);
    } finally {
      editSaving = false;
    }
  }

  // -- Delete --
  /** @param {any} app */
  function openDeleteModal(app) {
    deleteId = app.id;
    deleteName = app.name;
    deleteAlert = "";
    showDeleteModal = true;
  }

  async function onDeleteConfirm() {
    deleteSaving = true;
    deleteAlert = "";
    try {
      await deleteApp(deleteId);
      showDeleteModal = false;
      showAlert("success", "App deleted");
      await loadApps();
    } catch (err) {
      deleteAlert = toErrMsg(err);
    } finally {
      deleteSaving = false;
    }
  }

  /** @param {KeyboardEvent} e */
  function onGlobalKeydown(e) {
    if (e.key !== "Escape") return;
    if (showCreateModal) showCreateModal = false;
    else if (showEditModal) showEditModal = false;
    else if (showDeleteModal) showDeleteModal = false;
    else revealSecret = "";
  }
</script>

<svelte:window onkeydown={onGlobalKeydown} />

<UserLayout bind:alertMsg {alertType}>
  {#if revealSecret}
    <div class="alert alert-warning alert-dismissible fade show">
      <strong>Save your App Secret now!</strong> It will not be shown again.
      <div class="mt-2 p-2 bg-white rounded">
        <code class="user-select-all" style="word-break:break-all">{revealSecret}</code>
        <button
          class="btn btn-sm btn-outline-secondary ms-2"
          onclick={() => copyText(revealSecret)}>Copy</button
        >
      </div>
      <button
        type="button"
        class="btn-close"
        aria-label="Close"
        onclick={() => (revealSecret = "")}
      ></button>
    </div>
  {/if}

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
        <h5 class="mb-0">Applications</h5>
        <button class="btn btn-primary btn-sm" onclick={openCreateModal}>
          New Application
        </button>
      </div>
      <div class="card-body p-3">
        {#if apps.length === 0}
          <div class="text-center py-5 text-muted">
            <p class="mb-2">No applications found</p>
            <p class="small">Create an application to integrate with your services</p>
          </div>
        {:else}
          <div class="table-responsive">
            <table class="table">
              <thead>
                <tr>
                  <th>App ID</th>
                  <th>Name</th>
                  <th>Callback URL</th>
                  <th>Secret</th>
                  <th class="text-end">Actions</th>
                </tr>
              </thead>
              <tbody>
                {#each apps as app}
                  <tr>
                    <td>
                      <code class="user-select-all" style="font-size:0.8em"
                        >{app.id.slice(0, 8)}...</code
                      >
                      <button
                        class="btn btn-sm btn-link p-0 ms-1"
                        style="font-size:0.75em"
                        onclick={() => copyText(app.id)}
                        title="Copy App ID">Copy</button
                      >
                    </td>
                    <td>{app.name}</td>
                    <td>
                      {#if app.url}
                        <span class="small" style="word-break:break-all">{app.url}</span>
                      {:else}
                        <span class="text-muted">-</span>
                      {/if}
                    </td>
                    <td>
                      <code class="small">{app.secret_key || "****"}</code>
                    </td>
                    <td class="text-end text-nowrap">
                      <button
                        class="btn btn-outline-primary btn-sm"
                        onclick={() => openEditModal(app)}
                        >Edit</button
                      >
                      <button
                        class="btn btn-outline-danger btn-sm"
                        onclick={() => openDeleteModal(app)}
                        >Delete</button
                      >
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

<!-- Create Modal -->
{#if showCreateModal}
  <div class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,0.5)" role="dialog">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">New Application</h5>
          <button
            type="button"
            class="btn-close"
            aria-label="Close"
            onclick={() => (showCreateModal = false)}
          ></button>
        </div>
        <div class="modal-body">
          {#if createAlert}
            <div class="alert alert-danger py-2 small">{createAlert}</div>
          {/if}
          <div class="mb-3">
            <label class="form-label" for="appName">App Name</label>
            <input
              id="appName"
              type="text"
              class="form-control"
              bind:value={createName}
              placeholder="e.g. My Website"
              required
            />
          </div>
          <div class="mb-3">
            <label class="form-label" for="appUrl">Callback URL</label>
            <input
              id="appUrl"
              type="url"
              class="form-control"
              bind:value={createUrl}
              placeholder="https://yourdomain.com/api/auth/callback"
            />
            <div class="form-text">IAM redirects here after user login</div>
          </div>
        </div>
        <div class="modal-footer">
          <button
            class="btn btn-secondary btn-sm"
            onclick={() => (showCreateModal = false)}>Cancel</button
          >
          <button
            class="btn btn-primary btn-sm"
            onclick={onCreateSave}
            disabled={createSaving}
          >
            {#if createSaving}
              <span class="spinner-border spinner-border-sm me-1" role="status"></span>
              Creating...
            {:else}
              Create
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<!-- Edit Modal -->
{#if showEditModal}
  <div class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,0.5)" role="dialog">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Edit Application</h5>
          <button
            type="button"
            class="btn-close"
            aria-label="Close"
            onclick={() => (showEditModal = false)}
          ></button>
        </div>
        <div class="modal-body">
          {#if editAlert}
            <div class="alert alert-danger py-2 small">{editAlert}</div>
          {/if}
          <div class="mb-3">
            <label class="form-label" for="editName">App Name</label>
            <input
              id="editName"
              type="text"
              class="form-control"
              bind:value={editName}
              required
            />
          </div>
          <div class="mb-3">
            <label class="form-label" for="editUrl">Callback URL</label>
            <input
              id="editUrl"
              type="url"
              class="form-control"
              bind:value={editUrl}
              placeholder="https://yourdomain.com/api/auth/callback"
            />
          </div>
        </div>
        <div class="modal-footer">
          <button
            class="btn btn-secondary btn-sm"
            onclick={() => (showEditModal = false)}>Cancel</button
          >
          <button
            class="btn btn-primary btn-sm"
            onclick={onEditSave}
            disabled={editSaving}
          >
            {#if editSaving}
              <span class="spinner-border spinner-border-sm me-1" role="status"></span>
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

<!-- Delete Modal -->
{#if showDeleteModal}
  <div class="modal d-block" tabindex="-1" style="background:rgba(0,0,0,0.5)" role="dialog">
    <div class="modal-dialog modal-dialog-centered modal-sm">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Delete Application</h5>
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
            Delete application <strong>{deleteName}</strong>? This cannot be undone.
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
              <span class="spinner-border spinner-border-sm me-1" role="status"></span>
            {:else}
              Confirm Delete
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
