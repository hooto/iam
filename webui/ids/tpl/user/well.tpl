<div class="ids-container" style="padding:10px 0;">
<table id="ids-user-header">
  <tr>
    <td align="left" width="220px">
      <div class="iuh-brand">
        <img class="iuh-brand-logo" src="/ids/~/ids/img/ids-s2-32.png"> 
        <span class="iuh-brand-title">{[=it.webui_banner_title]}</span>
      </div>
    </td>
    <td>
      <div id="ids-uh-topnav">
        <a class="l4i-nav-item active" href="#user/overview">My Account</a>
        {[~it.topnav :v]}
        <a href="{[=v.path]}" class="l4i-nav-item">{[=v.title]}</a>
        {[~]}
      </div>
    </td>
    <td align="right">
      <div id="ids-userbox">
        <span class="btn btn-default btn-xs ids-userbox-signout" onclick="idsuser.SignOut()">Sign Out</span>
        <img class="ids-userbox-ico" src="/ids/~/ids/img/user-default.png">
      </div>
    </td>
  </tr>
</table>
</div>

<div id="com-content" class="ids-container"></div>

