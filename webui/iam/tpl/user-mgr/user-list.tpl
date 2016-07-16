<style type="text/css">
.pagination {
  margin: 10px 0;
}
</style>

<table width="100%">
  <tr>
    <td>
      <form id="iam-usermgr-list-query" action="#" class="form-inlines">
        <input id="query_text" type="text"
          class="form-control iam-input-query" 
          placeholder="Enter to search" 
          value="">
      </form>
    </td>
    <td align="right">
      <button type="button" 
        class="btn btn-primary btn-sm" 
        onclick="iamusrmgr.UserSetForm()">
        New User
      </button>
    </td>
  </tr>
</table>

<div id="iam-usermgr-list-alert" class="hide" style="margin:20px 0;"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Username</th>
      <th>Name</th>
      <th>Email</th>
      <th>Status</th>
      <th>Roles</th>
      <th>Registered</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="iam-usermgr-list"></tbody>
</table>
<div id="iam-usermgr-list-pager"></div>

<script id="iam-usermgr-list-tpl" type="text/html">
{[~it.items :v]}
    <tr>
      <td>{[=v.meta.name]}</td>
      <td>{[=v.name]}</td>
      <td>{[=v.email]}</td>
      <td>
        {[~it._statusls :sv]}
        {[ if (v.status == sv.status) { ]}{[=sv.title]}{[ } ]}
        {[~]}
      </td>
      <td>
        {[~v.roles :rv]}
        {[~it._roles.items :drv]}
        {[ if (drv.idxid == rv) { ]}
        <div>{[=drv.meta.name]}</div>
        {[ } ]}
        {[~]}
        {[~]}
      </td>
      <td>{[=l4i.TimeParseFormat(v.meta.created, "Y-m-d")]}</td>
      <td>{[=l4i.TimeParseFormat(v.meta.updated, "Y-m-d")]}</td>
      <td align="right">
        <a href="#user-mgr/user-set" onclick="iamusrmgr.UserSetForm('{[=v.meta.id]}')" 
          class="btn btn-default btn-xs">Setting</a>
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
