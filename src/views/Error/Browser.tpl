
<style>
body {
    width: 680px;
    margin: 0 auto;
    background-color: #eee;
}
.err-brw {
    margin: 40px;
    padding: 20px;
    width: 600px;
    border: 1px solid #ccc;
    background-color: #fff;
    border-radius: 4px;
}
.err-brw td {
    padding: 10px 20px 10px 0;
}
.err-brw .imgs1 {
    width: 32px; height: 32px; 
}
.err-footer {
    width: 600px;
    text-align: center;
}
</style>

<div class="err-brw">
  <div class="">
    <div class="alert alert-danger">{{T . "browser-reject-desc"}}</div>
    
    <p>{{T . "browser-reject-advice-desc"}}</p>
    <table width="100%">
      <tr>
        <td><img src="/ids/~/lessui/master/img/browser/chrome.png" class="imgs1" /></td>
        <td><strong>Google Chrome</strong></td>
        <td><a href="http://www.google.com/chrome/" target="_blank">http://www.google.com/chrome/</a></td>
      </tr>
      <tr>
        <td><img src="/ids/~/lessui/master/img/browser/safari.png" class="imgs1" /></td>
        <td><strong>Apple Safari</strong></td>
        <td><a href="http://www.apple.com/safari/" target="_blank">http://www.apple.com/safari/</a></td>
      </tr>
      <tr>
        <td><img src="/ids/~/lessui/master/img/browser/firefox.png" class="imgs1" /></td>
        <td><strong>Mozilla Firefox</strong></td>
        <td><a href="http://www.mozilla.org/" target="_blank">http://www.mozilla.org/</a></td>
      </tr>
    </table>
    
  </div>

</div>

<div class="err-footer">
    &copy; 2014 <a href="http://lesscompute.com" target="_blank">lessCompute.com</a>
</div>
