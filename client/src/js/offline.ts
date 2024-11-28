function recover() {
  const host = document.location.host;
  const wsProtocol =
    document.location.protocol === "https:" ? "wss://" : "ws://";
  window.conn = new WebSocket(wsProtocol + host + "/ws");
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
