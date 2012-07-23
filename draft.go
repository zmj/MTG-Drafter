package main

import (
	"fmt"
)

type DraftTable struct {
	Connections []WsConnection
}

var draftFmt = new(TestFormat)

func (d *DraftTable) RunDraft() {
	pack := draftFmt.BuildDraftPacks(1)[0]

	for len(pack) > 0 {
		msg := MessageOut{
			Msg: "Pack",
			Pack: MessagePack{ 1, 1, pack },
		}
		d.Connections[0].Send <- msg
		
		pick := <-d.Connections[0].Receive
		fmt.Println(pick.Pick)

		for i, card := range pack {
			if card.Name == pick.Pick {
				pack = append(pack[:i], pack[i+1:]...)
				break
			}
		}
	}
}