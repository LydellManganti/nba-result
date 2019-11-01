package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const baseUrl = "http://data.nba.net/10s"
const todayUrl = "/prod/v1/today.json"

func main() {
	var data []byte
	response, err := http.Get(baseUrl + todayUrl)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ = ioutil.ReadAll(response.Body)
	}

	var nba Nba
	json.Unmarshal([]byte(data), &nba)
	getTodaysResult(nba)
}

func getTodaysResult(nba Nba) {
	var data []byte
	fmt.Println("Todays Result", nba.Links.TodayScoreboard)
	response, err := http.Get(baseUrl + nba.Links.TodayScoreboard)
	if err != nil {
		fmt.Printf("The HTTP Request 101 failed with error %s\n", err)
	} else {
		data, _ = ioutil.ReadAll(response.Body)
	}

	var scoreBoard ScoreBoard
	json.Unmarshal([]byte(data), &scoreBoard)
	getScoreboard(scoreBoard)
}

func getScoreboard(scoreBoard ScoreBoard) {
	for _, game := range scoreBoard.Games {
		fmt.Printf("%s vs %s\n", game.VTeam.TriCode, game.HTeam.TriCode)
		hScore, _ := strconv.Atoi(game.HTeam.Score)
		vScore, _ := strconv.Atoi(game.VTeam.Score)
		if hScore > vScore {
			fmt.Printf("  %s win   : %s - %s\n", game.HTeam.TriCode, game.HTeam.Score, game.VTeam.Score)
		} else {
			fmt.Printf("  %s win   : %s - %s\n", game.VTeam.TriCode, game.VTeam.Score, game.HTeam.Score)
		}
		fmt.Printf("  Home      : %s\n", game.HTeam.TriCode)
		fmt.Printf("  Location  : %s, %s, %s\n", game.Arena.Name, game.Arena.City, game.Arena.StateAbbr)
		fmt.Printf("  Highlight : %s\n", game.Nugget.Text)
	}
}

type Nba struct {
	Links Links
}

type Links struct {
	TodayScoreboard string
}

type ScoreBoard struct {
	Games    []Game
	NumGames int
}

type Game struct {
	VTeam  Team
	HTeam  Team
	Arena  Arena
	Nugget Nugget
}

type Team struct {
	TriCode string
	Win     string
	Loss    string
	Score   string
}

type Arena struct {
	Name      string
	City      string
	StateAbbr string
}

type Nugget struct {
	Text string
}
