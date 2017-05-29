<style type="text/css">
.pagination {
  margin: 10px 0;
}
</style>

<div class="iam-div-light">
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
</div>

<script id="iam-usermgr-rolelist-tpl" type="text/html">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">{[=v.id]}</td>
  <td>{[=v.name]}</td>
  <td>{[=v.desc]}</td>
  <td>
    {[~it._statusls :sv]}
    {[ if (v.status == sv.status) { ]}{[=sv.title]}{[ } ]}
    {[~]}
  </td>
  <td>{[=l4i.MetaTimeParseFormat(v.created, "Y-m-d")]}</td>
  <td>{[=l4i.MetaTimeParseFormat(v.updated, "Y-m-d")]}</td>
  <td align="right">
    <button class="pure-button button-xsmall"
      onclick="iamUserMgr.RoleSet('{[=v.id]}')">
      Setting
    </button>
  </td>
</tr>
{[~]}
</script>

<script type="text/html" id="iam-usermgr-rolelist-optools">
<li class="iam-btn iam-btn-primary">
  <a href="#" onclick="iamUserMgr.RoleSet()">
     New Role
  </a>
</li>
</script>

