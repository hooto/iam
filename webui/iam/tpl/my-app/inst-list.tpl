<style type="text/css">
.pagination {
  margin: 10px 0;
}
</style>

<div id="iam-myapp-insts-alert" class="hide" style="margin:20px 0;"></div>

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
    <tbody id="iam-myapp-insts"></tbody>
  </table>

  <div id="iam-myapp-insts-pager"></div>
</div>

<script id="iam-myapp-insts-tpl" type="text/html">
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
  <td>{[=l4i.TimeParseFormat(v.meta.created, "Y-m-d")]}</td>
  <td>{[=l4i.TimeParseFormat(v.meta.updated, "Y-m-d")]}</td>
  <td align="right">
    <button class="pure-button button-xsmall"
      onclick="iamMyApp.InstSetForm('{[=v.meta.id]}')">
      Setting
    </button>
  </td>
</tr>
{[~]}
</script>


<script id="iam-myapp-insts-pager-tpl" type="text/html">
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
