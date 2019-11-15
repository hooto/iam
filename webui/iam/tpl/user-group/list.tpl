<style type="text/css">
.pagination {
  margin: 10px 0;
}
</style>

<div id="iam-user-group-list-alert" class="hide" style="margin:20px 0;"></div>
<div class="iam-div-light" id="iam-user-group-list"></div>
<div id="iam-user-group-list-pager"></div>

<script id="iam-user-group-list-tpl" type="text/html">
<table class="table table-hover valign-middle">
  <thead>
    <tr>
      <th>Name</th>
      <th>Display Name</th>
      <th>Owner</th>
      <th>Member</th>
      <th>Status</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody >
{[~it.items :v]}
<tr>
  <td class="iam-monofont">{[=v.name]}</td>
  <td>{[=v.display_name]}</td>
  <td>{[=v._owners]}</td>
  <td>{[=v._members]}</td>
  <td>
    {[~it._statusls :sv]}
    {[ if (v.status == sv.status) { ]}{[=sv.title]}{[ } ]}
    {[~]}
  </td>
  <td>{[=l4i.MetaTimeParseFormat(v.updated, "Y-m-d")]}</td>
  <td align="right">
    <button class="pure-button button-small"
      onclick="iamUserGroup.SetForm('{[=v.name]}')">
      {[=l4i.T("Settings")]}
    </button>
  </td>
</tr>
{[~]}
</tbody>
</table>
</script>

<script type="text/html" id="iam-user-group-list-optools">
<li class="iam-btn iam-btn-primary">
  <a href="#" onclick="iamUserGroup.SetForm()">
     {[=l4i.T("New %s", l4i.T("Group"))]} 
  </a>
</li>
</script>
