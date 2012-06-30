package main

import (
    "io"
)


type DraftArchive interface {
    GetDraft(id string) (io.WriterTo, bool)
}

type DraftFileArchive struct {
    
}

func (arch *DraftFileArchive) GetDraft(id string) (io.WriterTo, bool) {
    //try to pull it from a dict
    //if not there, try to read from disk and load to dict
    //if that fails, return false
    return nil, false
}