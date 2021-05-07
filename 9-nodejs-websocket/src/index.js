var WebSocketServer = require("ws").Server;
var wss = new WebSocketServer({ port: 8080 });

wss.on("connection", function (ws) {
    ws.send("Hello! I am a server.");
    ws.on("message", function (message) {
        console.log("Received: %s", message);
    });
});
