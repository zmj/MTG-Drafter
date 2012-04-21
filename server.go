package main

import (
    "fmt"
)

func main() {
    serv := new(Server)
    go serv.Serve()
}

type Server struct {
    var publicDrafts map[string] *Draft
    var activeDrafts map[string] *Draft
}

func (serv *Server) Serve() {
    for {
        var reqById *DraftPageByIdRequest
        var reqByFormat *DraftPageByFormatRequest
        select {
        case reqById = <-requestsById:
            //http draft page request (id) -> (html string)
            d, ok := serv.activeDrafts[reqById.id]
            if ok {
                reqById.response <- d.page
            } else {
                archive.requests <- reqById
            }
        case reqByFormat = <-requestsByFormat:
            //http join draft request (format) -> (html string)
            //validate format string
            var page *DraftPage            
            if reqByFormat.private || serv.publicDrafts[reqByFormat.format] == nil {
                d, err := serv.CreateDraft(reqByFormat)
                if err != nil {
                    page = DraftPage{nil, err}
                } else {
                    page = d.page
                }
            }
            reqByFormat.response <- page
        //ws join draft request (id) -> ?
        //draft ending (draft obj) which thread should be responsible for writing finished draft data?
        //draft full (draft obj)
        }
    }
}

type DraftPageByIdRequest struct {
    id string
    response chan *DraftPage
}

type DraftPageByFormatRequest struct {
    format string
    private bool
    response chan *DraftPage
}

func (serv *Server) CreateDraft(req *DraftPageByFormatRequest) *Draft, error {
    d, err := NewDraft(req.format, req.private)
    if err != nil {
        return nil, err
    }
    serv.activeDrafts[d.Id] = d
    if !d.Private {
        serv.publicDrafts[d.Format] = d
    }
    return d, nil
}

type DraftPage struct {
    *fmt.Stringer
    err error
}