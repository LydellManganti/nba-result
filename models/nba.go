package models

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

func (nba Nba) GetTodaysSchedule() ([]string, ScoreBoard) {
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
		schedule = append(schedule, "    "+game.HTeam.TriCode+" vs "+game.VTeam.TriCode+"    ")
	}
	return schedule, scoreBoard
}

func (game Game) GetTeamStats(anchorDate string) Boxscore {
	var data []byte
	response, err := http.Get(baseURL + fmt.Sprintf("/prod/v1/%s/%s_boxscore.json", anchorDate, game.GameId))
	if err != nil {
		fmt.Printf("The HTTP Request 101 failed with error %s\n", err)
	} else {
		data, _ = ioutil.ReadAll(response.Body)
	}

	var boxscore Boxscore
	json.Unmarshal([]byte(data), &boxscore)
	return boxscore
}

// GetDisplayHighlight Retrieve Highlight Information to be Displayed
func (game Game) GetDisplayHighlight() DisplayHighlight {
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
	} else if gameStatus[game.StatusNum] == "In Progress" {
		if hScore == vScore {
			displayHighlight.Result = fmt.Sprintf("  Game Tied : %s - %s\n", game.HTeam.Score, game.VTeam.Score)
		} else if hScore > vScore {
			displayHighlight.Result = fmt.Sprintf("  %s leads : %s - %s\n", game.HTeam.TriCode, game.HTeam.Score, game.VTeam.Score)
		} else {
			displayHighlight.Result = fmt.Sprintf("  %s leads : %s - %s\n", game.VTeam.TriCode, game.VTeam.Score, game.HTeam.Score)
		}
	}
	displayHighlight.Highlight = fmt.Sprintf("  Highlight : %s\n", game.Nugget.Text)
	return displayHighlight
}

func (game Game) GetDisplayStandings() DisplayStandings {
	var displayStandings DisplayStandings
	displayStandings.Header = "Team Win Loss\n"
	displayStandings.HomeTeam = fmt.Sprintf("%s   %s   %s\n", game.HTeam.TriCode, game.HTeam.Win, game.HTeam.Loss)
	displayStandings.VisitorTeam = fmt.Sprintf("%s   %s   %s", game.VTeam.TriCode, game.VTeam.Win, game.VTeam.Loss)
	return displayStandings
}

func (game Game) IsBoxscoreDataRetrieved(boxscores []Boxscore) bool {
	for _, boxscore := range boxscores {
		if boxscore.BasicGameData.GameId == game.GameId {
			return true
		}
	}
	return false
}

func (game Game) RetrieveBoxscoreData(boxscores []Boxscore) Boxscore {
	for _, boxscore := range boxscores {
		if boxscore.BasicGameData.GameId == game.GameId {
			return boxscore
		}
	}
	var emptyBoxscore Boxscore
	return emptyBoxscore
}
