package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseURL  = "http://data.nba.net/10s"
	todayURL = "/prod/v1/today.json"
)

var gameStatus = map[int]string{
	1: "Not Started",
	2: "In Progress",
	3: "Finished",
}

type Nba struct {
	Links Links
}

type Links struct {
	AnchorDate      string
	TodayScoreboard string
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

type ScoreBoard struct {
	Games    []Game
	NumGames int
}

type Game struct {
	VTeam        Team
	HTeam        Team
	Arena        Arena
	GameId       string
	Nugget       Nugget
	Period       Period
	StatusNum    int
	Clock        string
	StartTimeUTC string
}

type Team struct {
	TeamID    string
	TriCode   string
	Win       string
	Loss      string
	Score     string
	LineScore []Score
}

type Score struct {
	Score string
}

type Arena struct {
	Name      string
	City      string
	StateAbbr string
}

type Nugget struct {
	Text string
}

type Period struct {
	Current       int
	IsHalfTime    bool
	IsEndOfPeriod bool
}

type Boxscore struct {
	BasicGameData BasicGameData
	Stats         Stats
}

type BasicGameData struct {
	GameId string
}

type Stats struct {
	TimesTied     string
	LeadChanges   string
	VTeam         TeamStat
	HTeam         TeamStat
	ActivePlayers []ActivePlayer
}

type ActivePlayer struct {
	FirstName string
	LastName  string
	Jersey    string
	TeamID    string
	IsOnCourt bool
	Points    string
	Min       string
	FGM       string
	FGA       string
	FGP       string
	FTM       string
	FTA       string
	FTP       string
	TPM       string
	TPA       string
	TPP       string
	OffReb    string
	DefReb    string
	TotReb    string
	Assists   string
	PFouls    string
	Steals    string
	Turnovers string
	Blocks    string
	PlusMinus string
	DNP       string
}

type TeamStat struct {
	FastBreakPoints    string
	PointsInPaint      string
	BiggestLead        string
	SecondChancePoints string
	PointsOffTurnovers string
	LongestRun         string
	Totals             TeamTotals
}

type TeamTotals struct {
	Points    string
	FGM       string
	FGA       string
	FGP       string
	FTM       string
	FTA       string
	FTP       string
	TPM       string
	TPA       string
	TPP       string
	OffReb    string
	DefReb    string
	TotReb    string
	Assists   string
	PFouls    string
	Steals    string
	Turnovers string
	Blocks    string
}

type NbaData interface {
	GetNBAData() Nba
	GetTodaysSchedule(nba Nba) ([]string, ScoreBoard)
	GetTeamStats(game Game, anchorDate string) Boxscore
}

type NbaService interface {
	NbaData
}

type nbaService struct {
	NbaData
	BaseURL    string
	TodayURL   string
	GameStatus map[int]string
}

func NewNbaService() NbaService {
	return &nbaService{
		BaseURL:    baseURL,
		TodayURL:   todayURL,
		GameStatus: gameStatus,
	}
}

func (ns *nbaService) GetNBAData() Nba {
	var data []byte
	response, err := http.Get(ns.BaseURL + ns.TodayURL)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ = ioutil.ReadAll(response.Body)
	}
	var nba Nba
	json.Unmarshal([]byte(data), &nba)
	return nba
}

func (ns *nbaService) GetTodaysSchedule(nba Nba) ([]string, ScoreBoard) {
	var data []byte
	response, err := http.Get(ns.BaseURL + nba.Links.TodayScoreboard)
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

func (ns *nbaService) GetTeamStats(game Game, anchorDate string) Boxscore {
	var data []byte
	response, err := http.Get(ns.BaseURL + fmt.Sprintf("/prod/v1/%s/%s_boxscore.json", anchorDate, game.GameId))
	if err != nil {
		fmt.Printf("The HTTP Request 101 failed with error %s\n", err)
	} else {
		data, _ = ioutil.ReadAll(response.Body)
	}

	var boxscore Boxscore
	json.Unmarshal([]byte(data), &boxscore)
	return boxscore
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
