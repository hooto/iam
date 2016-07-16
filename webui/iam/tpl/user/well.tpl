<div class="iam-container" style="padding:10px 0;">
<table id="iam-user-header">
  <tr>
    <td align="left" width="220px">
      <div class="iuh-brand">
        <img class="iuh-brand-logo" src="/iam/~/iam/img/iam-s2-32.png"> 
        <span class="iuh-brand-title">{[=it.webui_banner_title]}</span>
      </div>
    </td>
    <td>
      <div id="iam-uh-topnav">
        <a class="l4i-nav-item active" href="#user/overview">My Account</a>
        {[~it.topnav :v]}
        <a href="{[=v.path]}" class="l4i-nav-item">{[=v.title]}</a>
        {[~]}
      </div>
    </td>
    <td align="right">
      <div id="iam-userbox">
        <span class="btn btn-default btn-xs iam-userbox-signout" onclick="iamuser.SignOut()">Sign Out</span>
        <img class="iam-userbox-ico" src="/iam/~/iam/img/user-default.png">
      </div>
    </td>
  </tr>
</table>
</div>

<div id="com-content" class="iam-container"></div>

