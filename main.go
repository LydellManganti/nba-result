package main

import (
	"fmt"
	"log"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var fifteenSpace = strings.Repeat(" ", 15)
var tenSpace = strings.Repeat(" ", 10)

func main() {
	var boxscores []Boxscore
	nba := GetNBAData()
	schedule, scoreBoard := nba.GetTodaysSchedule()
	// for _, game := range scoreBoard.Games {
	// 	game.GetTeamStats(nba.Links.AnchorDate)
	// }

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	uiSchedule := widgets.NewList()
	uiSchedule.Title = "Today's Games"
	uiSchedule.Rows = schedule
	uiSchedule.TextStyle = ui.NewStyle(ui.ColorWhite)
	uiSchedule.WrapText = false
	uiSchedule.BorderStyle.Fg = ui.ColorCyan
	uiSchedule.SetRect(0, 0, 20, 40)

	uiHighlight := widgets.NewParagraph()
	uiHighlight.Title = "Highlights"
	uiHighlight.SetRect(21, 0, 90, 10)
	uiHighlight.TextStyle.Fg = ui.ColorWhite
	uiHighlight.BorderStyle.Fg = ui.ColorCyan

	uiStandings := widgets.NewParagraph()
	uiStandings.Title = "Standings"
	uiStandings.SetRect(91, 0, 110, 10)
	uiStandings.TextStyle.Fg = ui.ColorWhite
	uiStandings.BorderStyle.Fg = ui.ColorCyan

	uiTeamStats := widgets.NewParagraph()
	uiTeamStats.Title = "Team Statistics"
	uiTeamStats.Text = "0"
	uiTeamStats.SetRect(21, 40, 110, 10)
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

	updateTeamStatsWidget := func(game Game, boxscore Boxscore) {
		uiTeamStats.Text = fmt.Sprintf("\n%s%s%s      Team Stats      %s%s\n\n", tenSpace, addPrefixSpace(game.HTeam.TriCode), fifteenSpace, fifteenSpace, addSuffixSpace(game.VTeam.TriCode))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s\n", uiTeamStats.Text, tenSpace, strings.Repeat("=", 65))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s        Points        %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.Totals.Points), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.Points))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s       Rebounds       %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.Totals.TotReb), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.TotReb))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s       Assists        %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.Totals.Assists), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.Assists))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s        Blocks        %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.Totals.Blocks), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.Blocks))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s      Field Goal      %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.Totals.FGM+"/"+boxscore.Stats.HTeam.Totals.FGA), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.FGM+"/"+boxscore.Stats.VTeam.Totals.FGA))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s         FG %%         %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.Totals.FGP), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.FGP))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s   3 Pt Field Goal    %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.Totals.TPM+"/"+boxscore.Stats.HTeam.Totals.TPA), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.TPM+"/"+boxscore.Stats.VTeam.Totals.TPA))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s      3 Pt FG %%       %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.Totals.TPP), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.TPP))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s      Free Throw      %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.Totals.FTM+"/"+boxscore.Stats.HTeam.Totals.FTA), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.FTM+"/"+boxscore.Stats.VTeam.Totals.FTA))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s         FT %%         %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.Totals.FTP), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.Totals.FTP))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s  Fast Break Points   %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.FastBreakPoints), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.FastBreakPoints))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s   Points In Paint    %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.PointsInPaint), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.PointsInPaint))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s    Biggest Lead      %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.BiggestLead), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.BiggestLead))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s Second Chance Points %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.SecondChancePoints), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.SecondChancePoints))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s Points Off Turnovers %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.PointsOffTurnovers), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.PointsOffTurnovers))
		uiTeamStats.Text = fmt.Sprintf("%s%s%s%s     Longest Run      %s%s\n", uiTeamStats.Text, tenSpace, addPrefixSpace(boxscore.Stats.HTeam.LongestRun), fifteenSpace, fifteenSpace, addSuffixSpace(boxscore.Stats.VTeam.LongestRun))
	}

	game := scoreBoard.Games[0]
	boxscore := game.GetTeamStats(nba.Links.AnchorDate)
	boxscores = append(boxscores, boxscore)
	highlightInformation := game.GetDisplayHighlight()
	standingsInformation := game.GetDisplayStandings()
	updateHighlightWidget(highlightInformation)
	updateStandingsWidget(standingsInformation)
	updateTeamStatsWidget(game, boxscore)
	ui.Render(uiSchedule, uiHighlight, uiStandings, uiTeamStats)

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
			updateTeamStatsWidget(game, boxscore)
		case "k", "<Up>":
			uiSchedule.ScrollUp()
			game := scoreBoard.Games[uiSchedule.SelectedRow]
			boxscore = game.retrieveBoxscoreData(boxscores)
			highlightInformation := game.GetDisplayHighlight()
			standingsInformation := game.GetDisplayStandings()
			updateHighlightWidget(highlightInformation)
			updateStandingsWidget(standingsInformation)
			updateTeamStatsWidget(game, boxscore)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(uiSchedule, uiHighlight, uiStandings, uiTeamStats)
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

func addPrefixSpace(s string) string {
	formattedString := strings.Repeat(" ", 6) + s
	return formattedString[len(formattedString)-6:]
}

func addSuffixSpace(s string) string {
	formattedString := s + strings.Repeat(" ", 6)
	return formattedString[:len(formattedString)-6]
}
