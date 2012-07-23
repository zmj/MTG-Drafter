function startDraft() {
	$('#waiting').hide();
	$('#pack').show();
	//initialize deck section
}

function showPack(cards) {
    $('#cards').empty();
	var images = $.map(cards, function(card, i) { 
		return '<img class="cardImage" data-card="' + card.Name + '" src="' + card.ImageURL + '"/>'; 
	});
	$('#cards').append(images.join(''));
	$('#cards').one("click", ".cardImage", function(e) { 
		pick($(e.target).data('card'));
	});
}

function pick(card) {
	wsSend({ "Msg":"Pick", "Pick":card })
}

var ws;

function wsSend(msg) {
    ws.send(JSON.stringify(msg)); 
}

function wsMessage(e) {
    var msg = JSON.parse(e.data);
    //alert("received: "+msg.Msg);
    //types of messages
        //error
        //draft full
        //player count update
        //draft starting
        //pack
        //forced pick
        //draft over
    if(msg.Msg == "Pack") {
    	showPack(msg.Pack.Cards);
    }
}

function wsOpen(e) {
    //wsSend({ Msg:'join', DraftId:draftId });
}

function wsError(e) {
    alert("error: "+e);
}

function wsClose(e) {
    alert("closed: "+e.code+" "+e.reason);
}

$(function () {	
	var wsUrl = 'ws://' + document.domain + '/ws/' + draftId;
	ws = new WebSocket(wsUrl);
    ws.onopen = wsOpen;
    ws.onclose = wsClose;
    ws.onerror = wsError;
    ws.onmessage = wsMessage;

	startDraft();
	//showPack(testPack);
});