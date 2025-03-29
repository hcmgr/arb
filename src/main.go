package main

import (
	"fmt"
)

type Sport struct {
	SportKey     string `json:"key"`
	HasOutrights bool   `json:"has_outrights"`
}

// Together, Match, Bookmaker, Market and Outcome
// objects define the structure of a Match returned by
// the odds API.
//
// JSON parses them directly.
type Match struct {
	MatchId      string      `json:"id"`
	SportKey     string      `json:"sport_key"`
	SportTitle   string      `json:"sport_title"`
	CommenceTime string      `json:"commence_time"`
	HomeTeam     string      `json:"home_team"`
	AwayTeam     string      `json:"away_team"`
	Bookmakers   []Bookmaker `json:"bookmakers"`
}

type Bookmaker struct {
	BookmakerKey string   `json:"key"`
	Markets      []Market `json:"markets"`
}

type Market struct {
	Outcomes []Outcome `json:"outcomes"`
}

type Outcome struct {
	BookmakerKey string  `json:"-"` // NOTE: inserted in post-processing (i.e. not read from raw json)
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
}

// Represents an arbitrage opportunity
// i.e. sum of 1/o guaranteed to be < 1
type Arb struct {
	// match metadata
	MatchId      string
	SportKey     string
	SportTitle   string
	CommenceTime string
	HomeTeam     string
	AwayTeam     string

	// arb info
	Outcomes []*Outcome
	R        float64
}

func (arb *Arb) toString() {
	fmt.Printf("arb: %s %.5f (%.5f)", arb.SportKey, arb.R, (1-arb.R)*100)
	for _, o := range arb.Outcomes {
		fmt.Print(" ", o.BookmakerKey, " ", o.Name, " ", o.Price)
	}
	fmt.Println()
}

func findMatchArbs(match *Match, arbs *[]Arb) {
	// calculate most frequent number of outcomes
	numOutcomesFreqs := make(map[int]int)
	for i := range match.Bookmakers {
		bookmaker := &match.Bookmakers[i]
		market := &bookmaker.Markets[0] // NOTE: always take first market, rarely more than one
		numOutcomesFreqs[len(market.Outcomes)]++
	}
	mostFrequentNumOutcomes := findMaxKey(numOutcomesFreqs)

	switch mostFrequentNumOutcomes {
	case 0:
		// fmt.Println("No odds")
	case 2:
		findTwoWayMatchArbs(match, arbs)
		break
	case 3:
		findThreeWayMatchArbs(match, arbs)
		break
	default:
		fmt.Println("Only support 2 and 3-outcome matches:", mostFrequentNumOutcomes)
		return
	}
}

// global config
var config *Config

// global db
var db *Database

// sources
const SOURCE_DB int = 0
const SOURCE_API int = 1
const SOURCE_FILE int = 2

func findArbs() []Arb {
	// get sports list
	var sports []Sport

	sports = getSports()

	arbs := make([]Arb, 0)
	cnt := 0
	for _, sport := range sports {
		if sport.HasOutrights {
			continue
		}
		if cnt > 100 {
			break
		}
		cnt++

		sportKey := sport.SportKey

		// get odds for this sportkey and write to db
		matches := getSportMatches(sportKey)

		for i := range matches {
			match := &matches[i]
			findMatchArbs(match, &arbs)
		}
	}

	// write arbs to db
	db.writeArbs(arbs)

	return arbs
}

func getArbs() []Arb {
	arbs, err := db.readArbs()
	if err != nil {
		panic(err)
	}
	return arbs
}

func showArbs(arbs []Arb) {
	for _, arb := range arbs {
		arb.toString()
	}
}

func main() {
	initConfig()
	initDb()

	arbs := findArbs()

	// arbs := getArbs()
	showArbs(arbs)
}
