package controllers

import (
	"log"
	"nba-result/models"

	ui "github.com/gizak/termui/v3"
)

func Render(l *models.Layout, ls models.LayoutService, ns models.NbaService) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	boxscores := []models.Boxscore{}
	sb, anchorDate := getScheduleData(l, ls, ns)
	g, b := getGamesData(sb, ns, anchorDate, 0, boxscores)
	setScreenData(l, ls, g, b)
	ui.Render(l.Schedule, l.Highlight, l.Standings, l.QuarterScores, l.GameTime, l.TeamStats, l.BoxScore)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.Schedule.ScrollDown()
			g, b := getGamesData(sb, ns, anchorDate, l.Schedule.SelectedRow, boxscores)
			setScreenData(l, ls, g, b)
		case "k", "<Up>":
			l.Schedule.ScrollUp()
			g, b := getGamesData(sb, ns, anchorDate, l.Schedule.SelectedRow, boxscores)
			setScreenData(l, ls, g, b)
		case "r":
			sb, anchorDate := getScheduleData(l, ls, ns)
			g, b := getGamesData(sb, ns, anchorDate, l.Schedule.SelectedRow, boxscores)
			setScreenData(l, ls, g, b)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}
		ui.Render(l.Schedule, l.Highlight, l.Standings, l.QuarterScores, l.GameTime, l.TeamStats, l.BoxScore)
	}
}

func getScheduleData(l *models.Layout, ls models.LayoutService, ns models.NbaService) (models.ScoreBoard, string) {
	nba := ns.GetNBAData()
	schedule, scoreBoard := ns.GetTodaysSchedule(nba)
	ls.UpdateScheduleWidget(l.Schedule, schedule)
	return scoreBoard, nba.Links.AnchorDate
}

func getGamesData(scoreBoard models.ScoreBoard, ns models.NbaService, anchorDate string, row int, boxscores []models.Boxscore) (models.Game, models.Boxscore) {
	game := scoreBoard.Games[row]
	boxscore := models.Boxscore{}
	if game.IsBoxscoreDataRetrieved(boxscores) {
		boxscore = game.RetrieveBoxscoreData(boxscores)
	} else {
		boxscore = ns.GetTeamStats(game, anchorDate)
		boxscores = append(boxscores, boxscore)
	}
	return game, boxscore
}

func setScreenData(l *models.Layout, ls models.LayoutService, game models.Game, boxscore models.Boxscore) {
	highlightInformation := ls.GetDisplayHighlight(game)
	standingsInformation := ls.GetDisplayStandings(game)
	ls.UpdateHighlightWidget(l.Highlight, highlightInformation)
	ls.UpdateStandingsWidget(l.Standings, standingsInformation)
	ls.UpdateQuarterScoresWidget(l.QuarterScores, game)
	ls.UpdateGameTimeWidget(l.GameTime, game)
	ls.UpdateTeamStatsWidget(l.TeamStats, game, boxscore)
	ls.UpdateBoxScoreWidget(l.BoxScore, game, boxscore)
}
