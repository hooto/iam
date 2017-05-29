<style type="text/css">
.pagination {
  margin: 10px 0;
}
</style>

<div id="iam-appmgr-instls-alert" class="hide" style="margin:20px 0;"></div>

<div class="iam-div-light">
  <table class="table table-hover">
    <thead>
      <tr>
        <th>ID</th>
        <th>App ID</th>
        <th>App Name</th>
        <th>Version</th>
        <th>Status</th>
        <th>Created</th>
        <th>Updated</th>
        <th></th>
      </tr>
    </thead>
    <tbody id="iam-appmgr-instls"></tbody>
  </table>

  <div id="iam-appmgr-instls-pager"></div>
</div>

<script id="iam-appmgr-instls-tpl" type="text/html">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">{[=v.meta.id]}</td>
  <td>{[=v.app_id]}</td>
  <td>{[=v.app_title]}</td>
  <td>{[=v.version]}</td>
  <td>
    {[~it._statusls :sv]}
    {[ if (v.status == sv.status) { ]}{[=sv.title]}{[ } ]}
    {[~]}
  </td>
  <td>{[=l4i.MetaTimeParseFormat(v.meta.created, "Y-m-d")]}</td>
  <td>{[=l4i.MetaTimeParseFormat(v.meta.updated, "Y-m-d")]}</td>
  <td align="right">
  <button class="pure-button button-xsmall"
      onclick="iamAppMgr.InstSetForm('{[=v.meta.id]}')">
      Setting
    </button>
  </td>
</tr>
{[~]}
</script>


<script id="iam-appmgr-instls-pager-tpl" type="text/html">
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

<script type="text/html" id="iam-appmgr-instls-optools">
<li>
  <form action="#" onsubmit="iamAppMgr.InstList()" class="form-inlines">
    <input id="iam_appmgr_instls_qry_text" type="text"
      class="form-control iam-query-input"
      placeholder="Press Enter to Search"
      value="">
  </form>
</li>
</script>

