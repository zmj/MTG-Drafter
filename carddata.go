package main

import (
	"fmt"
	"strings"
)

type Card struct {
	Name string
	Set string
	Rarity string
	Cost string	
	ImageURL string //computed
}

type Format interface {
	BuildDraftPacks(players int) [][]*Card
	GetName() string
}

func (c *Card) SetImageUrl() {
	escapedName := strings.Map(escapeWizardsUrl, c.Name)
	c.ImageURL = fmt.Sprintf("http://www.wizards.com/global/images/magic/%s/%s.jpg", c.Set, escapedName)
}

func escapeWizardsUrl(c rune) rune {
	switch c {
	case ' ':
		return '_'
	}
	return c
}


//test stuff
type TestFormat struct {

}

func (t *TestFormat) BuildDraftPacks(players int) [][]*Card {
	pack := []*Card{
		&Card{
			Name: "Vedalken Entrancer",
			Cost: "3U",		
			Set: "M13",
			Rarity: "C",
		},
		&Card{
			Name: "Wild Guess",
			Cost: "RR",
			Set: "M13",
			Rarity: "C",
		},
		&Card{
			Name: "Dark Favor",
			Cost: "1B",
			Set: "M13",
			Rarity: "C",
		},
		&Card{
			Name: "Vile Rebirth",
			Cost: "B",
			Set: "M13",
			Rarity: "C",
		},
		&Card{
			Name: "Canyon Minotaur",
			Cost: "3R",			
			Set: "M13",
			Rarity: "C",
		},
		&Card{
			Name: "Aven Squire",
			Cost: "1W",
			Set: "M13",
			Rarity: "C",			
		},
		&Card{
			Name: "Wind Drake",
			Cost: "2U",
			Set: "M13",
			Rarity: "C",			
		},
		&Card{
			Name: "Primal Huntbeast",
			Cost: "3G",
			Set: "M13",
			Rarity: "C",
		},
		&Card{
			Name: "Attended Knight",
			Cost: "2W",
			Set: "M13",
			Rarity: "C",
		},
		&Card{
			Name: "Ring of Valkas",
			Cost: "2",
			Set: "M13",
			Rarity: "U",
		},
		&Card{
			Name: "Courtly Provocateur",
			Cost: "2U",
			Set: "M13",
			Rarity: "U",
		},
		&Card{
			Name: "Prized Elephant",
			Cost: "3W",
			Set: "M13",
			Rarity: "U",
		},
		&Card{
			Name: "Void Stalker",
			Cost: "1U",
			Set: "M13",
			Rarity: "R",
		},
		&Card{
			Name: "Plains",
			Cost: "0",
			Set: "M13",
			Rarity: "Basic",
		},
		&Card{
			Name: "Ring of Evos Isle",
			Cost: "2",
			Set: "M13",
			Rarity: "U",
		},
	}
	for _,c := range pack {
		c.SetImageUrl()
	}
	return [][]*Card{ pack }
}

func (t *TestFormat) Getname() string {
	return "TST"
}