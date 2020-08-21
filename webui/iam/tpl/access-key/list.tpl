<style type="text/css">
.pagination {
  margin: 10px 0;
}

#iam-aklist .boundop {
  margin-right: 0 !important;
  padding-right: 0 !important;
}
#iam-aklist .boundop > a {
  margin: 0 0 0 5px;
  padding: 2px 5px;
  color: #fff;
  text-align: center;
  width: 20px;
  text-decoration: none;
}
#iam-aklist .boundop > a:hover {
  color: #fff;
  background-color: #d9534f;
}


</style>

<div class="iam-div-light">

<table class="table table-hover valign-middle">
<thead>
  <tr>
    <th>Access Key ID</th>
    <th>{[=l4i.T("Description")]}</th>
    <th>Status</th>
    <th width="35%">{[=l4i.T("Scopes")]}</th>
    <th></th>
  </tr>
</thead>
<tbody id="iam-aklist">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">
    <a href="#info" onclick="iamAccessKey.Info('{[=v.id]}')">{[=v.id]}</a>
  </td>
  <td>{[=v.description]}</td>
  <td>
    {[~it._statuses :sv]}
    {[ if (v.status == sv.status) { ]}{[=sv.title]}{[ } ]}
    {[~]}
  </td>
  <td>
    {[~v.scopes :bv]}
    <span class="label label-success boundop">{[=bv.name]} = {[=bv.value]} <a href="#" onclick="iamAccessKey.UnBind('{[=v.id]}', '{[=bv.name]}')">&times;</a></span>
    {[~]}
  </td>
  <td align="right">
   <button class="pure-button button-small"
      onclick="iamAccessKey.Bind('{[=v.id]}')">
      <span class="fa fa-plus"></span>
      {[=l4i.T("New Scope")]}
    </button>
    <button class="pure-button button-small"
      onclick="iamAccessKey.Del('{[=v.id]}')">
      <span class="fa fa-times-circle"></span>
      {[=l4i.T("Delete")]}
    </button>
    <button class="pure-button button-small"
      onclick="iamAccessKey.Set('{[=v.id]}')">
      <span class="fa fa-cog"></span>
      {[=l4i.T("Settings")]}
    </button>
  </td>
</tr>
{[~]}
</tbody>
</table>

</div>

<script type="text/html" id="iam-aklist-optools">
<li class="iam-btn iam-btn-primary">
  <a href="#" onclick="iamAccessKey.Set()">
     {[=l4i.T("New %s", "Access Key")]}
  </a>
</li>
</script>

