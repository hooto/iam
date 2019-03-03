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
        <th>Access Key</th>
        <th>Description</th>
        <th width="35%">Bounds</th>
        <th>Action</th>
        <th>Created</th>
        <th></th>
      </tr>
    </thead>
    <tbody id="iam-aklist"></tbody>
  </table>
</div>

<script id="iam-aklist-tpl" type="text/html">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">
    <a href="#info" onclick="iamAccessKey.Info('{[=v.access_key]}')">{[=v.access_key]}</a>
  </td>
  <td>{[=v.desc]}</td>
  <td>
    {[~v.bounds :bv]}
    <span class="label label-success boundop">{[=bv.name]} <a href="#" onclick="iamAccessKey.UnBind('{[=v.access_key]}', '{[=bv.name]}')">&times;</a></span>
    {[~]}
  </td>
  <td>
    {[~it._actionls :sv]}
    {[ if (v.action == sv.action) { ]}{[=sv.title]}{[ } ]}
    {[~]}
  </td>
  <td>{[=l4i.MetaTimeParseFormat(v.created, "Y-m-d")]}</td>
  <td align="right">
    <button class="pure-button button-small"
      onclick="iamAccessKey.Del('{[=v.access_key]}')">
      Delete
    </button>
    <button class="pure-button button-small"
      onclick="iamAccessKey.Bind('{[=v.access_key]}')">
      Bind New
    </button>
    <button class="pure-button button-small"
      onclick="iamAccessKey.Set('{[=v.access_key]}')">
      Setting
    </button>
  </td>
</tr>
{[~]}
</script>

<script type="text/html" id="iam-aklist-optools">
<li class="iam-btn iam-btn-primary">
  <a href="#" onclick="iamAccessKey.Set()">
     New Access Key
  </a>
</li>
</script>

