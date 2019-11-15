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

	uiHiglight := widgets.NewParagraph()
	uiHiglight.Title = "Highlights"
	uiHiglight.SetRect(21, 0, 85, 10)
	uiHiglight.TextStyle.Fg = ui.ColorWhite
	uiHiglight.BorderStyle.Fg = ui.ColorCyan

	boxScore := widgets.NewParagraph()
	boxScore.Title = "Box Score"
	boxScore.Text = "0"
	boxScore.SetRect(0, 40, 85, 10)
	boxScore.BorderStyle.Fg = ui.ColorCyan

	updateHighlightWidget := func(displayHighlight DisplayHighlight) {
		uiHiglight.Text = displayHighlight.Versus + "\n"
		uiHiglight.Text = fmt.Sprintf("%s%s", uiHiglight.Text, displayHighlight.Status)
		uiHiglight.Text = fmt.Sprintf("%s%s", uiHiglight.Text, displayHighlight.Home)
		uiHiglight.Text = fmt.Sprintf("%s%s", uiHiglight.Text, displayHighlight.Location)
		uiHiglight.Text = fmt.Sprintf("%s%s", uiHiglight.Text, displayHighlight.Result)
		uiHiglight.Text = fmt.Sprintf("%s%s", uiHiglight.Text, displayHighlight.Highlight)
	}

	game := scoreBoard.Games[0]
	highlightInformation := GetDisplayHighlight(game)
	updateHighlightWidget(highlightInformation)
	ui.Render(scheduleList, uiHiglight, boxScore)

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
			updateHighlightWidget(highlightInformation)
		case "k", "<Up>":
			scheduleList.ScrollUp()
			game := scoreBoard.Games[scheduleList.SelectedRow]
			highlightInformation := GetDisplayHighlight(game)
			updateHighlightWidget(highlightInformation)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(scheduleList, uiHiglight, boxScore)
	}
}
