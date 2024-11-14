function recover() {
  window.conn = new WebSocket("ws://" + document.location.host + "/ws");
  window.conn.onopen = function (evt) {
    window.conn.onmessage = function (evt) {
      var data = JSON.parse(evt.data);
      if (data.type === "settingschanged") {
        if (data.payload.online) {
          window.conn.close();
          window.location.reload();
        }
      }
    };
  };
}

if (document.readyState !== "loading") {
  recover();
}

document.addEventListener("DOMContentLoaded", function () {
  recover();
});
