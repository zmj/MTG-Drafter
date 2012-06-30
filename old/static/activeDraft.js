;

//BEGIN TEMPLATE DATA
var draftId = '12345'
//END TEMPLATE DATA

var ws;

function wsSend(msg) {
    ws.send(JSON.stringify(msg)); 
}

function wsMessage(e) {
    var msg = JSON.parse(e.data);
    alert("received: "+msg.Msg);
    //types of messages
        //error
        //draft full
        //player count update
        //draft starting
        //pack
        //forced pick
        //draft over
}

function wsOpen(e) {
    wsSend({ Msg:'join', DraftId:draftId });
}

function wsError(e) {
    alert("error: "+e);
}

function wsClose(e) {
    alert("closed: "+e.code+" "+e.reason);
}

$(function() {
    ws = new WebSocket("ws://localhost:80/ws");
    ws.onopen = wsOpen;
    ws.onclose = wsClose;
    ws.onerror = wsError;
    ws.onmessage = wsMessage;
});