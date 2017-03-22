<style type="text/css">
.pagination {
  margin: 10px 0;
}
</style>

<table width="100%">
  <tr>
    <td>
      <!-- <form id="kqqfqk" action="#" class="form-inlines">
        <input id="query_text" type="text"
          class="form-control iam-input-query" 
          placeholder="Enter to search" 
          value="">
      </form> -->
    </td>
    <td align="right">
      <button type="button" 
        class="btn btn-primary btn-sm" 
        onclick="iamUserMgr.RoleSetForm()">
        New Role
      </button>
    </td>
  </tr>
</table>

<table class="table table-hover">
  <thead>
    <tr>
      <th>ID</th>
      <th>Name</th>
      <th>Description</th>
      <th>Status</th>
      <th>Created</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="iam-usermgr-rolelist"></tbody>
</table>

<script id="iam-usermgr-rolelist-tpl" type="text/html">
{[~it.items :v]}
    <tr>
      <td>{[=v.meta.id]}</td>
      <td>{[=v.meta.name]}</td>
      <td>{[=v.desc]}</td>
      <td>
        {[~it._statusls :sv]}
        {[ if (v.status == sv.status) { ]}{[=sv.title]}{[ } ]}
        {[~]}
      </td>
      <td>{[=l4i.TimeParseFormat(v.meta.created, "Y-m-d")]}</td>
      <td>{[=l4i.TimeParseFormat(v.meta.updated, "Y-m-d")]}</td>
      <td align="right">
        <a href="#user-mgr/role-set" onclick="iamUserMgr.RoleSetForm('{[=v.meta.id]}')" class="btn btn-default btn-xs">Setting</a>
      </td>
    </tr>
{[~]}
</script>
