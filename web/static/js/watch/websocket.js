const host = "localhost"
let socket = new WebSocket("http://localhost:8080/ws");

socket.onopen = function(e) {
  console.log("connected");

      };
      
      socket.onmessage = function(event) {
  console.log("message received");
      };
      
      socket.onclose = function(event) {
  console.log("connection closed");
      };
      
      socket.onerror = function(error) {
  console.log("error: " + error);
      };