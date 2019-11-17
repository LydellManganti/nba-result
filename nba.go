package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const baseURL = "http://data.nba.net/10s"
const todayURL = "/prod/v1/today.json"

var gameStatus = map[int]string{
	1: "Not Started",
	2: "In Progress",
	3: "Finished",
}

func GetNBAData() Nba {
	var data []byte
	response, err := http.Get(baseURL + todayURL)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ = ioutil.ReadAll(response.Body)
	}

	var nba Nba
	json.Unmarshal([]byte(data), &nba)
	return nba
}

func GetTodaysSchedule(nba Nba) ([]string, ScoreBoard) {
	var data []byte
	response, err := http.Get(baseURL + nba.Links.TodayScoreboard)
	if err != nil {
		fmt.Printf("The HTTP Request 101 failed with error %s\n", err)
	} else {
		data, _ = ioutil.ReadAll(response.Body)
	}

	var scoreBoard ScoreBoard
	json.Unmarshal([]byte(data), &scoreBoard)
	var schedule []string
	for _, game := range scoreBoard.Games {
		schedule = append(schedule, "   ["+game.HTeam.TriCode+" vs "+game.VTeam.TriCode+"](fg:cyan)")
	}
	return schedule, scoreBoard
}

// GetDisplayHighlight Retrieve Highlight Information to be Displayed
func GetDisplayHighlight(game Game) DisplayHighlight {
	var displayHighlight DisplayHighlight
	displayHighlight.Versus = fmt.Sprintf("%s vs %s\n", game.HTeam.TriCode, game.VTeam.TriCode)
	displayHighlight.Status = fmt.Sprintf("  Status    : %s\n", gameStatus[game.StatusNum])
	displayHighlight.Home = fmt.Sprintf("  Home      : %s\n", game.HTeam.TriCode)
	displayHighlight.Location = fmt.Sprintf("  Location  : %s, %s, %s\n", game.Arena.Name, game.Arena.City, game.Arena.StateAbbr)
	hScore, _ := strconv.Atoi(game.HTeam.Score)
	vScore, _ := strconv.Atoi(game.VTeam.Score)
	if gameStatus[game.StatusNum] == "Finished" {
		if hScore > vScore {
			displayHighlight.Result = fmt.Sprintf("  %s win   : %s - %s\n", game.HTeam.TriCode, game.HTeam.Score, game.VTeam.Score)
		} else {
			displayHighlight.Result = fmt.Sprintf("  %s win   : %s - %s\n", game.VTeam.TriCode, game.VTeam.Score, game.HTeam.Score)
		}
		displayHighlight.Highlight = fmt.Sprintf("  Highlight : %s\n", game.Nugget.Text)
	} else if gameStatus[game.StatusNum] == "In Progress" {
		if hScore > vScore {
			displayHighlight.Result = fmt.Sprintf("  %s leads : %s - %s\n", game.HTeam.TriCode, game.HTeam.Score, game.VTeam.Score)
		} else {
			displayHighlight.Result = fmt.Sprintf("  %s leads : %s - %s\n", game.VTeam.TriCode, game.VTeam.Score, game.HTeam.Score)
		}
	}
	return displayHighlight
}

func GetDisplayStandings(game Game) DisplayStandings {
	var displayStandings DisplayStandings
	displayStandings.Header = "Team Win Loss\n"
	displayStandings.HomeTeam = fmt.Sprintf("%s   %s   %s\n", game.HTeam.TriCode, game.HTeam.Win, game.HTeam.Loss)
	displayStandings.VisitorTeam = fmt.Sprintf("%s   %s   %s", game.VTeam.TriCode, game.VTeam.Win, game.VTeam.Loss)
	return displayStandings
}

type DisplayHighlight struct {
	Versus    string
	Status    string
	Home      string
	Location  string
	Result    string
	Highlight string
}

type DisplayStandings struct {
	Header      string
	HomeTeam    string
	VisitorTeam string
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
	VTeam     Team
	HTeam     Team
	Arena     Arena
	Nugget    Nugget
	StatusNum int
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
