package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"
	"path"
)

var cardReg *regexp.Regexp = regexp.MustCompile(`<td align=[^>]*>(\d+)</td>\s+<td><a [^>]*>([^<]*)</a></td>\s+<td>([^<]*)</td>\s+<td>([^<]*)</td>\s+<td>([^<]*)</td>`)

func ScrapeSet(setName string) {
	url := "http://magiccards.info/" + strings.ToLower(setName) + "/en.html"
	response, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}

	set := new(BoosterSet)
	set.Name = setName
	set.Cards = make([]Card, 0)
	for card := range ParseCards(string(content)) {
		card.Set = setName
		set.Cards = append(set.Cards, card)
	}

	WriteSet(set, setName)
}

func ParseCards(body string) (chan Card) {
	results := make(chan Card)
	go func() {
		defer close(results)

		for _, cardInfo := range cardReg.FindAllStringSubmatch(body, -1) {
			card := new(Card)
			card.Name = cardInfo[2]
			card.Rarity = cardInfo[5]
			card.CastingCost = cardInfo[4]
			results <- *card
		}
	}()
	return results
}

func WriteSet(set *BoosterSet, fileName string) {
	setData, err := json.MarshalIndent(set, "", "\t")
	if err != nil {
		panic(err.Error())
	}
	err = ioutil.WriteFile(path.Join(SetPath, fileName), setData, 0644)
	if err != nil {
		panic(err.Error())
	}
}