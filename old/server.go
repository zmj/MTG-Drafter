package main

import (
    "fmt"
    "net/http"
    "net/url"
    "regexp"
    "io"
    "code.google.com/p/go.net/websocket"
)

func main() {
    serv := new(Server)
    serv.Serve()
}

type Server struct {
    publicQueues map[string] *Draft //format -> Draft
    activeDrafts map[string] *Draft //id -> Draft
    newDraftRequests chan *NewDraftRequest
    draftFull chan *Draft
    draftFinished chan *Draft
    archive DraftArchive
}

func (serv *Server) Serve() {
    serv.archive = new(DraftFileArchive)

    //initialize collections
    serv.publicQueues = make(map [string] *Draft)
    serv.activeDrafts = make(map [string] *Draft)

    //initialize channels
    //buffer these
    serv.newDraftRequests = make(chan *NewDraftRequest)
    serv.draftFull = make(chan *Draft)
    serv.draftFinished = make(chan *Draft)

    //register http handlers
    http.HandleFunc("/draft/", func(w http.ResponseWriter, r *http.Request) { 
        serv.HandlePageRequestById(w, r) })
    http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
        serv.HandlePageRequestByFormat(w, r) })    
    http.Handle("/", http.FileServer(http.Dir("./static"))) //does this cache?
    http.Handle("/ws", websocket.Handler(func (ws *websocket.Conn) { 
        serv.HandleWS(ws) }))

    //serve it up
    go serv.ManageDraftLifetimes()
    err := http.ListenAndServe(":80", nil)
    if err != nil {
        panic("Serve failed: " + err.Error())
    }
}

func (serv *Server) HandleWS(ws *websocket.Conn) {
    var msg WsMessageIn
    recErr := websocket.JSON.Receive(ws, &msg)
    if recErr != nil {
        fmt.Println(recErr.Error())
        return
    } 
    fmt.Println("message: ", msg)            
    d, ok := serv.activeDrafts[msg.DraftId]
    if !ok {
        websocket.JSON.Send(ws, WsMessageOut{ "full" })
        return
    }
    d.AddPlayer(ws)
}

func (serv *Server) ManageDraftLifetimes() {
    for {
        var d *Draft
        var req *NewDraftRequest
        select {
        case req = <-serv.newDraftRequests:            
            d = NewDraft(req.format, req.private)
            serv.activeDrafts[d.Id] = d
            if !d.Private {
                serv.publicQueues[d.Format] = d
            }
            req.response <- d
        
        case d = <-serv.draftFull:            
            delete(serv.publicQueues, d.Format)

        case d = <-serv.draftFinished:
            delete(serv.activeDrafts, d.Id)
        }
    }
}

type NewDraftRequest struct {
    format string
    private bool
    response chan *Draft
}

func (serv *Server) GetPublicDraft(format string) *Draft {
    d, ok := serv.publicQueues[format]
    if !ok {
        return serv.CreateDraft(format, false)
    }
    return d
}

func (serv *Server) CreateDraft(format string, private bool) *Draft {
    req := NewDraftRequest{format, private, make(chan *Draft)}
    serv.newDraftRequests <- &req
    return <-req.response
}

func (serv *Server) HandlePageRequestByFormat(w http.ResponseWriter, r *http.Request) {
    query, err := url.ParseQuery(r.URL.RawQuery)
    if err != nil {
        http.Error(w, err.Error(), 404)
    }
    format := query.Get("format")
    private := query.Get("private") == "true"
    //validate format
    fmt.Println("req by format", format, private)
    var draftId string
    if private {
        draftId = serv.CreateDraft(format, true).Id
    } else {
        draftId = serv.GetPublicDraft(format).Id
    }
    http.Redirect(w, r, "/draft/"+draftId, 302)
}

var draftPageReg, _ = regexp.Compile(`^/draft/(\d+)$`) //change this to StripPrefix
func (serv *Server) HandlePageRequestById(w http.ResponseWriter, r *http.Request) {
    m := draftPageReg.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.Error(w, "Not found", 404)
        return
    }
    id := m[1]
    fmt.Println("req by id", id)

    var d io.WriterTo
    var ok bool
    d, ok = serv.activeDrafts[id]
    if !ok {
        d, ok = serv.archive.GetDraft(id)
    }
    if !ok {
        //404
    }
    _, err := d.WriteTo(w)
    if err != nil {
        //server error
    }
}
