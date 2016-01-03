package main

import (
)

type Pack struct {
	Cards []Card
	Id string
}

type Card struct {
	Name string
	CastingCost string
	Set string
	Rarity string
	ImageUrl string
}

func (p *Pack) RemoveCard(card Card) {
	for i, c := range p.Cards {		
		if c.Name == card.Name {
			p.Cards = append(p.Cards[:i], p.Cards[i+1:]...)
			return
		} 
	}
	panic("not found")
}
