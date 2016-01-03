package main

import (
	"math/rand"
	"io/ioutil"
	"encoding/json"
	"path"
	"time"
)

var (
	FormatPath = "./formats"
	SetPath = "./formats/sets"
	CubePath = "./formats/cubes"
)

type Format interface {
	GeneratePacks(int) [][]Pack
}

type BoosterFormat struct {
	Name string
	Sets []BoosterSet
}

type BoosterSet struct {
	Name string
	Collation []BoosterCollation
	Cards []Card

	cardsByRarity map[string] []Card //rarity -> cards
}

type BoosterCollation map[string] int

type CubeFormat struct {}

func (format *BoosterFormat) GeneratePacks(numPlayers int) [][]Pack {
	draftPacks := make([][]Pack, len(format.Sets))
	for i, set := range format.Sets {
		packs := make([]Pack, numPlayers)
		for j:=0; j<numPlayers; j++ {
			packs[j] = set.GeneratePack()
		}
		draftPacks[i] = packs
	}
	return draftPacks
}

func (set *BoosterSet) GeneratePack() Pack {
	pack := new(Pack)
	pack.Cards = make([]Card, len(set.Collation))
	exclude := make(map[string] bool)
	for i, slot := range set.Collation {
		card := set.ChooseCard(slot, exclude)
		pack.Cards[i] = card
		exclude[card.Name] = true
	}
	return *pack
}

func (set *BoosterSet) ChooseCard(slot BoosterCollation, exclude map[string] bool) Card {
	//choose a rarity	
	rarity := ""
	total := 0
	for r, n := range slot {
		total += n
		rarity = r
	}
	if total > 1 {
		rarity = ""
		roll := rand.Intn(total) + 1
		for r, n := range slot {
			if roll <= n {
				rarity = r
				break
			} else {
				roll -= n
			}
		}
	}

	if rarity == "" {
		panic("what")
	}

	for {
		cardRoll := rand.Intn(len(set.cardsByRarity[rarity]))
		info := set.cardsByRarity[rarity][cardRoll]

		if _, excluded := exclude[info.Name]; excluded {
			continue
		}

		return Card {
			Name: info.Name,
			CastingCost: info.CastingCost,
			Set: info.Set,
			Rarity: info.Rarity,
		}
	}
	panic("unreachable")
}

func LoadFormats() (map[string] Format) {	
	rand.Seed(time.Now().UTC().UnixNano())

	sets := LoadSets()

	formats := make(map[string] Format)	
	for formatFileContents := range ReadFiles(FormatPath) {		
		var formatInfo map[string] interface{}
		parseFormatErr := json.Unmarshal(formatFileContents, &formatInfo)
		if parseFormatErr != nil {
			panic(parseFormatErr.Error())
		}
		switch formatInfo["Type"] {
		case "Booster":
			format := new(BoosterFormat)
			format.Name = formatInfo["Name"].(string)
			setNames := formatInfo["Sets"].([]interface{})
			format.Sets = make([]BoosterSet, len(setNames))
			for i, setName := range setNames {
				format.Sets[i] = sets[setName.(string)]
			}
			formats[format.Name] = format
		}
	}

	return formats
}

func LoadSets() (map[string] BoosterSet) {
	sets := make(map[string] BoosterSet)

	for setFileContents := range ReadFiles(SetPath) {
		var set BoosterSet
		parseSetErr := json.Unmarshal(setFileContents, &set)
		if parseSetErr != nil {
			panic(parseSetErr.Error())
		}

		//sort cards by rarity to make rolling up packs easier
		rarity := make(map[string] []Card)
		for _, card := range set.Cards {
			if _, exists := rarity[card.Rarity]; !exists {
				rarity[card.Rarity] = make([]Card, 0)
			}
			rarity[card.Rarity] = append(rarity[card.Rarity], card)
		}
		set.cardsByRarity = rarity

		sets[set.Name] = set
	}

	return sets
}

func ReadFiles(dirPath string) (chan []byte) {
	results := make(chan []byte)
	go func(results chan []byte) {
		defer close(results)
		dirContents, readDirErr := ioutil.ReadDir(dirPath)
		if readDirErr != nil {
			panic(readDirErr.Error())
		}
		for _, file := range dirContents {
			if file.IsDir() {
				continue
			}
			fileContents, readFileErr := ioutil.ReadFile(path.Join(dirPath, file.Name()))
			if readFileErr != nil {
				panic(readFileErr.Error())
			}

			results <- fileContents
		}

	}(results)
	return results
}