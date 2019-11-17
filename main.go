package main

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	nba := GetNBAData()
	schedule, scoreBoard := GetTodaysSchedule(nba)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	scheduleList := widgets.NewList()
	scheduleList.Title = "Today's Games"
	scheduleList.Rows = schedule
	scheduleList.TextStyle = ui.NewStyle(ui.ColorWhite)
	scheduleList.WrapText = false
	scheduleList.BorderStyle.Fg = ui.ColorCyan
	scheduleList.SetRect(0, 0, 20, 10)

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

	boxScore := widgets.NewParagraph()
	boxScore.Title = "Box Score"
	boxScore.Text = "0"
	boxScore.SetRect(0, 40, 90, 10)
	boxScore.BorderStyle.Fg = ui.ColorCyan

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

	game := scoreBoard.Games[0]
	highlightInformation := GetDisplayHighlight(game)
	standingsInformation := GetDisplayStandings(game)
	updateHighlightWidget(highlightInformation)
	updateStandingsWidget(standingsInformation)
	ui.Render(scheduleList, uiHighlight, uiStandings, boxScore)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			scheduleList.ScrollDown()
			game := scoreBoard.Games[scheduleList.SelectedRow]
			highlightInformation := GetDisplayHighlight(game)
			standingsInformation := GetDisplayStandings(game)
			updateHighlightWidget(highlightInformation)
			updateStandingsWidget(standingsInformation)
		case "k", "<Up>":
			scheduleList.ScrollUp()
			game := scoreBoard.Games[scheduleList.SelectedRow]
			highlightInformation := GetDisplayHighlight(game)
			standingsInformation := GetDisplayStandings(game)
			updateHighlightWidget(highlightInformation)
			updateStandingsWidget(standingsInformation)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(scheduleList, uiHighlight, uiStandings, boxScore)
	}
}
