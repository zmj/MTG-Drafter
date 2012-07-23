package main

type MessageIn struct {
	Msg string
	Pick string
}


type MessageOut struct {
	Msg string
	Pack MessagePack
}

type MessagePack struct {
	Pack int
	Pick int
	Cards []*Card
}