var sock = null;
var wsuri = "ws://localhost:8000/ws";

window.onload = function() {

    sock = new WebSocket(wsuri);

    sock.onopen = function() {
        console.log("connected to " + wsuri);
    }

    sock.onclose = function(e) {
        console.log("connection closed (" + e.code + ")");
    }

    sock.onmessage = function(e) {
        result = JSON.parse(e.data);
        console.log(result.usedPercent);
        $('#memory').html(result.usedPercent);
        //$('<li>' + e.data '</li>').hide().prependTo('#messages').fadeIn(1000);
    }
};