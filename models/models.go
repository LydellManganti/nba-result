package models

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
	TimesTied   string
	LeadChanges string
	VTeam       TeamStat
	HTeam       TeamStat
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
