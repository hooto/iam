<style type="text/css">
.pagination {
  margin: 10px 0;
}
</style>

<table width="100%">
  <tr>
    <td>
      <form id="gir0c7" action="#" class="form-inlines">
        <input id="query_text" type="text"
          class="form-control iam-input-query" 
          placeholder="Enter to search" 
          value="">
      </form>
    </td>
    <td align="right">
      <!-- <button type="button" 
        class="btn btn-primary btn-sm" 
        onclick="iamWorkLoader('user-mgr/auth-edit')">
        New Instance
      </button> -->
    </td>
  </tr>
</table>

<div id="iam-appmgr-insts-alert" class="hide" style="margin:20px 0;"></div>

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
  <tbody id="iam-appmgr-insts"></tbody>
</table>

<div id="iam-appmgr-insts-pager"></div>


<script id="iam-appmgr-insts-tpl" type="text/html">
{[~it.items :v]}
    <tr>
      <td>{[=v.meta.id]}</td>
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
        <a href="#app-mgr/inst-set" onclick="iamappmgr.InstSetForm('{[=v.meta.id]}')" class="btn btn-default btn-xs">Setting</a>
      </td>
    </tr>
{[~]}
</script>


<script id="iam-appmgr-insts-pager-tpl" type="text/html">
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
