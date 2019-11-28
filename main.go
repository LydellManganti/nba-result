package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var fifteenSpace = strings.Repeat(" ", 15)
var tenSpace = strings.Repeat(" ", 10)

func main() {
	var boxscores []Boxscore
	nba := GetNBAData()
	schedule, scoreBoard := nba.GetTodaysSchedule()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	uiSchedule := widgets.NewList()
	uiSchedule.Title = "Today's Games"
	uiSchedule.TitleStyle.Fg = ui.ColorClear
	uiSchedule.Rows = schedule
	uiSchedule.TextStyle.Fg = ui.ColorWhite
	uiSchedule.SelectedRowStyle.Fg = ui.ColorRed
	uiSchedule.SelectedRowStyle.Bg = ui.ColorWhite
	uiSchedule.BorderStyle.Fg = ui.ColorCyan
	uiSchedule.SetRect(0, 0, 20, 40)

	uiHighlight := widgets.NewParagraph()
	uiHighlight.Title = "Highlights"
	uiHighlight.TitleStyle.Fg = ui.ColorClear
	uiHighlight.SetRect(21, 0, 90, 10)
	uiHighlight.TextStyle.Fg = ui.ColorWhite
	uiHighlight.BorderStyle.Fg = ui.ColorCyan

	uiStandings := widgets.NewParagraph()
	uiStandings.Title = "Standings"
	uiStandings.SetRect(91, 0, 110, 10)
	uiStandings.TextStyle.Fg = ui.ColorWhite
	uiStandings.BorderStyle.Fg = ui.ColorCyan

	uiQuarterScores := widgets.NewParagraph()
	uiQuarterScores.Title = "Quarter Scores"
	uiQuarterScores.SetRect(21, 10, 90, 17)
	uiQuarterScores.TextStyle.Fg = ui.ColorWhite
	uiQuarterScores.BorderStyle.Fg = ui.ColorCyan

	uiGameTime := widgets.NewParagraph()
	uiGameTime.Title = "Game Time"
	uiGameTime.SetRect(91, 10, 110, 17)
	uiGameTime.TextStyle.Fg = ui.ColorWhite
	uiGameTime.BorderStyle.Fg = ui.ColorCyan

	uiTeamStats := widgets.NewParagraph()
	uiTeamStats.Title = "Team Statistics"
	uiTeamStats.Text = "0"
	uiTeamStats.SetRect(21, 40, 110, 17)
	uiTeamStats.BorderStyle.Fg = ui.ColorCyan

	updateHighlightWidget := func(displayHighlight DisplayHighlight) {
		uiHighlight.Text = displayHighlight.Versus + "\n"
		uiHighlight.Text = fmt.Sprintf("%s%s", uiHighlight.Text, displayHighlight.Status)
		uiHighlight.Text = fmt.Sprintf("%s%s", uiHighlight.Text, displayHighlight.Home)
		uiHighlight.Text = fmt.Sprintf("%s%s", uiHighlight.Text, displayHighlight.Location)
		uiHighlight.Text = fmt.Sprintf("%s%s", uiHighlight.Text, displayHighlight.Result)
		uiHighlight.Text = fmt.Sprintf("%s%s", uiHighlight.Text, displayHighlight.Highlight)
	}

	updateStandingsWidget := func(displayStandings DisplayStandings) {
		uiStandings.Text = displayStandings.Header + "\n"
		uiStandings.Text = fmt.Sprintf("%s%s", uiStandings.Text, displayStandings.HomeTeam)
		uiStandings.Text = fmt.Sprintf("%s%s", uiStandings.Text, displayStandings.VisitorTeam)
	}

	updateQuarterScoresWidget := func(game Game) {
		uiQuarterScores.Text = "      1st  2nd  3rd  4th\n"
		uiQuarterScores.Text = fmt.Sprintf("%s %s", uiQuarterScores.Text, game.HTeam.TriCode)
		for _, score := range game.HTeam.LineScore {
			uiQuarterScores.Text = fmt.Sprintf("%s  %s", uiQuarterScores.Text, addPrefixChar(score.Score, 3))
		}
		uiQuarterScores.Text = fmt.Sprintf("%s\n %s", uiQuarterScores.Text, game.VTeam.TriCode)
		for _, score := range game.VTeam.LineScore {
			uiQuarterScores.Text = fmt.Sprintf("%s  %s", uiQuarterScores.Text, addPrefixChar(score.Score, 3))
		}
	}

	updateGameTime := func(game Game) {
		if gameStatus[game.StatusNum] == "In Progress" {
			if game.Period.IsHalfTime {
				uiGameTime.Text = fmt.Sprintf("\n Halftime!")
			} else {
				uiGameTime.Text = fmt.Sprintf("\n Period: %s\n", strconv.Itoa(game.Period.Current))
				uiGameTime.Text = fmt.Sprintf("%s Time  : %s", uiGameTime.Text, game.Clock)
			}
		} else if gameStatus[game.StatusNum] == "Not Started" {
			usGameTime, err := time.Parse(time.RFC3339, game.StartTimeUTC)
			if err != nil {
				uiGameTime.Text = fmt.Sprintf("\n%s", err)
			}
			localGameTime := usGameTime.In(time.Local)
			uiGameTime.Text = fmt.Sprintf("\n Starts at %s:%s\n", addPrefixChar(strconv.Itoa(localGameTime.Hour()), 2, "0"), addPrefixChar(strconv.Itoa(localGameTime.Minute()), 2, "0"))
		} else {
			uiGameTime.Text = "\n Finished!"
		}
	}

	updateTeamStatsWidget := func(game Game, boxscore Boxscore) {
		uiTeamStats.Text = fmt.Sprintf("\n%s%s%s      Team Stats      %s%s\n\n", tenSpace, addPrefixChar(game.HTeam.TriCode, 6), fifteenSpace, fifteenSpace, addSuffixSpace(game.VTeam.TriCode, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s\n", uiTeamStats.Text, tenSpace, strings.Repeat("=", 65))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s        Points        %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.Totals.Points, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.Points, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s       Rebounds       %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.Totals.TotReb, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.TotReb, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s       Assists        %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.Totals.Assists, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.Assists, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s        Blocks        %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.Totals.Blocks, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.Blocks, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s      Field Goal      %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.Totals.FGM+"/"+boxscore.Stats.HTeam.Totals.FGA, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.FGM+"/"+boxscore.Stats.VTeam.Totals.FGA, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s         FG %%         %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.Totals.FGP, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.FGP, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s   3 Pt Field Goal    %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.Totals.TPM+"/"+boxscore.Stats.HTeam.Totals.TPA, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.TPM+"/"+boxscore.Stats.VTeam.Totals.TPA, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s      3 Pt FG %%       %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.Totals.TPP, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.TPP, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s      Free Throw      %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.Totals.FTM+"/"+boxscore.Stats.HTeam.Totals.FTA, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.FTM+"/"+boxscore.Stats.VTeam.Totals.FTA, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s         FT %%         %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.Totals.FTP, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.FTP, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s  Fast Break Points   %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.FastBreakPoints, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.FastBreakPoints, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s   Points In Paint    %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.PointsInPaint, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.PointsInPaint, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s    Biggest Lead      %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.BiggestLead, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.BiggestLead, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s Second Chance Points %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.SecondChancePoints, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.SecondChancePoints, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s Points Off Turnovers %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.PointsOffTurnovers, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.PointsOffTurnovers, 6))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s     Longest Run      %s%s\n", uiTeamStats.Text, tenSpace, addPrefixChar(boxscore.Stats.HTeam.LongestRun, 6), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.LongestRun, 6))
	}

	game := scoreBoard.Games[0]
	boxscore := game.GetTeamStats(nba.Links.AnchorDate)
	boxscores = append(boxscores, boxscore)
	highlightInformation := game.GetDisplayHighlight()
	standingsInformation := game.GetDisplayStandings()
	updateHighlightWidget(highlightInformation)
	updateStandingsWidget(standingsInformation)
	updateQuarterScoresWidget(game)
	updateGameTime(game)
	updateTeamStatsWidget(game, boxscore)
	ui.Render(uiSchedule, uiHighlight, uiStandings, uiQuarterScores, uiGameTime, uiTeamStats)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			uiSchedule.ScrollDown()
			game := scoreBoard.Games[uiSchedule.SelectedRow]
			if game.isBoxscoreDataRetrieved(boxscores) {
				boxscore = game.retrieveBoxscoreData(boxscores)
			} else {
				boxscore = game.GetTeamStats(nba.Links.AnchorDate)
				boxscores = append(boxscores, boxscore)
			}
			highlightInformation := game.GetDisplayHighlight()
			standingsInformation := game.GetDisplayStandings()
			updateHighlightWidget(highlightInformation)
			updateStandingsWidget(standingsInformation)
			updateQuarterScoresWidget(game)
			updateGameTime(game)
			updateTeamStatsWidget(game, boxscore)
		case "k", "<Up>":
			uiSchedule.ScrollUp()
			game := scoreBoard.Games[uiSchedule.SelectedRow]
			if game.isBoxscoreDataRetrieved(boxscores) {
				boxscore = game.retrieveBoxscoreData(boxscores)
			} else {
				boxscore = game.GetTeamStats(nba.Links.AnchorDate)
				boxscores = append(boxscores, boxscore)
			}
			highlightInformation := game.GetDisplayHighlight()
			standingsInformation := game.GetDisplayStandings()
			updateHighlightWidget(highlightInformation)
			updateStandingsWidget(standingsInformation)
			updateQuarterScoresWidget(game)
			updateGameTime(game)
			updateTeamStatsWidget(game, boxscore)
		case "r":
			schedule, scoreBoard = nba.GetTodaysSchedule()
			game := scoreBoard.Games[uiSchedule.SelectedRow]
			var emptyBoxscores []Boxscore
			boxscore = game.GetTeamStats(nba.Links.AnchorDate)
			boxscores = append(emptyBoxscores, boxscore)
			highlightInformation := game.GetDisplayHighlight()
			standingsInformation := game.GetDisplayStandings()
			updateHighlightWidget(highlightInformation)
			updateStandingsWidget(standingsInformation)
			updateQuarterScoresWidget(game)
			updateGameTime(game)
			updateTeamStatsWidget(game, boxscore)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(uiSchedule, uiHighlight, uiStandings, uiQuarterScores, uiGameTime, uiTeamStats)
	}
}

func (game Game) isBoxscoreDataRetrieved(boxscores []Boxscore) bool {
	for _, boxscore := range boxscores {
		if boxscore.BasicGameData.GameId == game.GameId {
			return true
		}
	}
	return false
}

func (game Game) retrieveBoxscoreData(boxscores []Boxscore) Boxscore {
	for _, boxscore := range boxscores {
		if boxscore.BasicGameData.GameId == game.GameId {
			return boxscore
		}
	}
	var emptyBoxscore Boxscore
	return emptyBoxscore
}

func addPrefixChar(s string, length int, charOptional ...string) string {
	var char string
	if len(charOptional) == 0 {
		char = " "
	} else {
		char = charOptional[0]
	}
	formattedString := strings.Repeat(char, length) + s
	return formattedString[len(formattedString)-length:]
}

func addSuffixSpace(s string, length int) string {
	formattedString := s + strings.Repeat(" ", length)
	return formattedString[:len(formattedString)-length]
}
