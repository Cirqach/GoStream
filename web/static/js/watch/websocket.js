const host = "localhost"

const video = document.getElementById('myVideo');
const mediaSource = new MediaSource();
video.src = URL.createObjectURL(mediaSource);

let sourceBuffer;

let socket = new WebSocket("http://localhost:8080/ws");

socket.onopen = function(e) {
        alert("[open] Соединение установлено");
      };
      
      socket.onmessage = function(event) {
        alert(`[message] Данные получены с сервера: ${event.data}`);
      };
      
      socket.onclose = function(event) {
        if (event.wasClean) {
          alert(`[close] Соединение закрыто чисто, код=${event.code} причина=${event.reason}`);
        } else {
          // например, сервер убил процесс или сеть недоступна
          // обычно в этом случае event.code 1006
          alert('[close] Соединение прервано');
        }
      };
      
      socket.onerror = function(error) {
        alert(`[error]`);
      };