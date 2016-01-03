package main

import (
	"fmt"
	"sync"
	"strings"
)

type DraftTable struct {
	Packs [][]Pack
	Seats []DraftSeat
}

type DraftSeat struct {
	Number int
	Drafter Drafter
	Left chan *Pack
	Right chan *Pack
}

type Drafter interface {
	Pick(Pack) Card
}

func NewTable(drafters []Drafter, format Format) *DraftTable {
	table := new(DraftTable)

	table.Packs = format.GeneratePacks(len(drafters))

	channels := make([]chan *Pack, len(drafters))
	for i, _ := range channels {
		channels[i] = make(chan *Pack, len(drafters))
	}
	table.Seats = make([]DraftSeat, len(drafters))
	for i, drafter := range drafters {
		seat := new(DraftSeat)		
		seat.Drafter = drafter
		seat.Number = i
		seat.Left = channels[i % len(channels)]
		seat.Right = channels[(i+1) % len(channels)]
		table.Seats[i] = *seat
	}

	return table
}

func (table *DraftTable) Draft() {
	for roundNum, roundPacks := range table.Packs {
		passLeft := roundNum % 2 == 0
		wg := new(sync.WaitGroup)
		for i, seat := range table.Seats {
			wg.Add(1)
			go func(seat DraftSeat, pack Pack, passLeft bool) {
				seat.Draft(pack, passLeft)
				wg.Done()				
			}(seat, roundPacks[i], passLeft)
		}
		wg.Wait()
	}
}

func (seat *DraftSeat) Draft(initialPack Pack, passLeft bool) {	
	var passTo chan *Pack
	var receiveFrom chan *Pack
	if passLeft {
		passTo = seat.Left
		receiveFrom = seat.Right
	} else {
		passTo = seat.Right
		receiveFrom = seat.Left
	}
	
	for pack := &initialPack; len(pack.Cards)>0; pack = <-receiveFrom {
		//fmt.Printf("Seat %d drafting from pack %s\n", seat.Number, pack.Id)
		//add timeouts
		pick := seat.Drafter.Pick(*pack)
		//validate?
		pack.RemoveCard(pick)
		//fmt.Printf("Seat %d passing pack %s\n", seat.Number, pack.Id)
		passTo <- pack
	}
}

func main() {
	fmt.Println("hi")
	drafters := []Drafter { new(RandomDrafter), new(RandomDrafter), new(ConsoleDrafter) }
	format := LoadFormats()["RTR Block"]
	table := NewTable(drafters, format)
	table.Draft()
	fmt.Println("done")
}

type RandomDrafter struct {}

func (r *RandomDrafter) Pick(p Pack) Card {
	return p.Cards[0]
}

type TestFormat struct {}

func (f *TestFormat) GeneratePacks(n int) [][]Pack {
	result := make([][]Pack, 3)
	for round:=0; round<3; round++ {
		packs := make([]Pack, n)
		for i:=0; i<n; i++ {
			packs[i] = MakeTestPack(round, i)
		}
		result[round] = packs
	}
	return result
}

func MakeTestPack(round, seat int) Pack {
	p := new(Pack)
	p.Id = fmt.Sprintf("%d.%d", round, seat)
	p.Cards = make([]Card, 3)
	for i:=0; i<3; i++ {
		c := new(Card)
		c.Name = fmt.Sprintf("%d.%d.%d", round, seat, i)
		p.Cards[i] = *c
	}
	return *p
}

type ConsoleDrafter struct {}

func (cd *ConsoleDrafter) Pick(p Pack) Card {
	for _, c := range p.Cards {
		fmt.Println("> ", c.Name)
	}
	var name string
	_, _ = fmt.Scanln(&name)
	return Card { Name:strings.Replace(name, "^", " ", -1) }
}