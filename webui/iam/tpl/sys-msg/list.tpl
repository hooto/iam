<div class="iam-div-light">

<table class="table table-hover valign-middle">
<thead>
  <tr>
    <th>{[=l4i.T("User")]}</th>
    <th>{[=l4i.T("Email")]}</th>
    <th>{[=l4i.T("Subject")]}</th>
    <th>Action</th>
    <th>{[=l4i.T("Posted")]}</th>
    <th width="30px"></th>
  </tr>
</thead>

<tbody id="iam-msglist">
{[~it.items :v]}
<tr class="hover" onclick="iamSysMsg.Info('{[=v.id]}')">
  <td>{[=v.to_user]}</td>
  <td>{[=v.to_email]}</td>
  <td>{[=v.title]}</td>
  <td>
    {[~it._actionls :sv]}
    {[? v.action == sv.action]}{[=sv.title]}{[?]}
    {[~]}
  </td>
  <td>{[=l4i.UnixTimeFormat(v.posted, "Y-m-d H:i:s")]}</td>
  <td align="right">
    <span class="fa fa-chevron-right"></span>
  </td>
</tr>
{[~]}
</tbody>
</table>

</div>


