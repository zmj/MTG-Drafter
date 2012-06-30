package main

import (
    "io"    
    "code.google.com/p/go.net/websocket"
    "html/template"
    "math/rand"
    "fmt"
)

var (
    activeDraftTmpl, _ = template.ParseFiles("activeDraft.tmpl")
)

type Draft struct {
    Id string
    Format string
    Private bool
}

func GetDraftId() string {
    id := rand.Int63n(1e10)
    return fmt.Sprintf("%010d", id)
}

func NewDraft(format string, private bool) *Draft {
    id := GetDraftId()
    fmt.Println("New draft", id)
    return &Draft{ id, format, private }
}

func (d *Draft) WriteTo(w io.Writer) (int64, error) {
    err := activeDraftTmpl.Execute(w, d)
    if err != nil {
        return 0, err
    }
    return 0, nil
}

type WsMessageIn struct {
    Msg string
    DraftId string //msg == 'join'
}

type WsMessageOut struct {
    Msg string
}

func (d *Draft) AddPlayer(ws *websocket.Conn) {
    websocket.JSON.Send(ws, WsMessageOut{ "greetings, earthling" } )
    for {

    }
}