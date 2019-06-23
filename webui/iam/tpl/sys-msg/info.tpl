<table class="iam-formtable valign-middle">
<tbody>
  <tr>
    <td width="200px">ID</td>
    <td class="iam-monofont">
      {[=it.data.id]}
    </td>
  </tr>

  <tr>
    <td>{[=l4i.T("Status")]}</td>
    <td>
      {[~it._actionls :v]}
      {[? it.data.action == v.action]}{[=v.title]}{[?]}
      {[~]}
    </td>
  </tr>

  <tr>
    <td>{[=l4i.T("Recipient")]}</td>
    <td>
      {[=it.data.to_user]} {[? it.data.to_email]}({[=it.data.to_email]}){[?]}
    </td>
  </tr>

  {[? it.data.title]}
  <tr>
    <td>{[=l4i.T("Subject")]}</td>
    <td>
      {[=it.data.title]}
    </td>
  </tr>
  {[?]}

  {[? it.data.body]}
  <tr>
    <td>{[=l4i.T("Content")]}</td>
    <td>
      <pre class="pre-wrap">{[=it.data.body]}</pre>
    </td>
  </tr>
  {[?]}

  {[? it.data.retry && it.data.retry > 0]}
  <tr>
    <td>{[=l4i.T("Retry")]}</td>
    <td>
      {[=it.data.retry]}
    </td>
  </tr>
  {[?]}

  <tr>
    <td>{[=l4i.T("Created")]}</td>
    <td>{[=l4i.UnixTimeFormat(it.data.created, "Y-m-d H:i:s")]}</td>
  </tr>

  <tr>
    <td>{[=l4i.T("Posted")]}</td>
    <td>{[=l4i.UnixTimeFormat(it.data.posted, "Y-m-d H:i:s")]}</td>
  </tr>
</tbody>
</table>

