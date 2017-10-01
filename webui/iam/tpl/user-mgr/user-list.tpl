<style type="text/css">
.pagination {
  margin: 10px 0;
}
</style>

<div id="iam-usermgr-list-alert" class="hide" style="margin:20px 0;"></div>

<div class="iam-div-light">
<table class="table table-hover">
  <thead>
    <tr>
      <th>Username</th>
      <th>Display Name</th>
      <th>Email</th>
      <th>Status</th>
      <th>Roles</th>
      <th>ECoin</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="iam-usermgr-list"></tbody>
</table>
<div id="iam-usermgr-list-pager"></div>
</div>

<script id="iam-usermgr-list-tpl" type="text/html">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">{[=v.name]}</td>
  <td>{[=v.display_name]}</td>
  <td>{[=v.email]}</td>
  <td>
    {[~it._statusls :sv]}
    {[ if (v.status == sv.status) { ]}{[=sv.title]}{[ } ]}
    {[~]}
  </td>
  <td>
    {[~v.roles :rv]}
    {[~it._roles.items :drv]}
    {[ if (drv.id == rv) { ]}
    <div>{[=drv.name]}</div>
    {[ } ]}
    {[~]}
    {[~]}
  </td>
  <td>{[=v.ecoin_balance]}</td>
  <td>{[=l4i.MetaTimeParseFormat(v.updated, "Y-m-d")]}</td>
  <td align="right">
    <button class="pure-button button-xsmall"
      onclick="iamAccMgr.Recharge('{[=v.name]}')">
      Recharge
    </button>
    <button class="pure-button button-xsmall"
      onclick="iamUserMgr.UserSetForm('{[=v.name]}')">
      Setting
    </button>
  </td>
</tr>
{[~]}
</script>

<script id="iam-usermgr-list-pager-tpl" type="text/html">
<ul class="pagination pagination-sm">
  {[ if (it.FirstPageNumber > 0) { ]}
  <li><a href="#{[=it.FirstPageNumber]}">First</a></li>
  {[ } ]}
  {[~it.RangePages :v]}
  <li {[ if (v == it.CurrentPageNumber) { ]}class="active"{[ } ]}><a href="#{[=v]}">{[=v]}</a></li>
  {[~]}
  {[ if (it.LastPageNumber > 0) { ]}
  <li><a href="#{[=it.LastPageNumber]}">Last</a></li>
  {[ } ]}
</ul>
</script>


<script type="text/html" id="iam-usermgr-list-optools">
<li>
  <form id="iam-usermgr-list-query" action="#" onsubmit="iamUserMgr.UserList()" class="form-inlines">
    <input id="iam_usermgr_list_qry_text" type="text"
      class="form-control iam-query-input" 
      placeholder="Press Enter to Search" 
      value="">
  </form>
</li>
<li class="iam-btn iam-btn-primary">
  <a href="#" onclick="iamUserMgr.UserSetForm()">
     New User
  </a>
</li>
</script>
