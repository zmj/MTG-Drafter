package main

import (
	"fmt"
	"net/http"
	"html/template"
    "code.google.com/p/go.net/websocket"
)

func main() {
	handle("/static/", http.FileServer(http.Dir("./static")))
	handle("/draft/", http.HandlerFunc(getDraftPage))
	handle("/ws/", websocket.Handler(wsJoin))

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err.Error())
	}
}

func handle(prefix string, h http.Handler) {
	http.Handle(prefix, http.StripPrefix(prefix, h))
}

type WsConnection struct {
	Receive chan MessageIn
	Send chan MessageOut
}

var draft = new(DraftTable)

func wsJoin(ws *websocket.Conn) {
	//try to get the draft
	//if it doesn't exist, close

	conn := WsConnection{ make(chan MessageIn), make(chan MessageOut) }
	//send to draft, handle error

	//test code
	draft.Connections = []WsConnection{ conn }
	go draft.RunDraft()

	go func() {
		defer close(conn.Send)
		for msgOut := range conn.Send {
			websocket.JSON.Send(ws, msgOut)
			//handle error
		}
	}()
	defer close(conn.Receive)
	var msgIn MessageIn
	for {
		websocket.JSON.Receive(ws, &msgIn)
		//handle error
		conn.Receive <- msgIn
	}
}

var draftTmpl, _ = template.ParseFiles("draft.tmpl")

func getDraftPage(wr http.ResponseWriter, req *http.Request) {
	err := draftTmpl.Execute(wr, nil)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(wr, err.Error(), http.StatusInternalServerError)
	}
}