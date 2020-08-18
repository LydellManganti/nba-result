package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var fifteenSpace = strings.Repeat(" ", 15)
var tenSpace = strings.Repeat(" ", 10)

type Layout struct {
	Schedule      *widgets.List
	Highlight     *widgets.Paragraph
	Standings     *widgets.Paragraph
	QuarterScores *widgets.Paragraph
	GameTime      *widgets.Paragraph
	TeamStats     *widgets.Paragraph
	BoxScore      *widgets.Paragraph
}

type NbaLayout interface {
	GetDisplayHighlight(Game) DisplayHighlight
	GetDisplayStandings(Game) DisplayStandings
	Initialise() *Layout
	SetSchedule(*widgets.List)
	SetHighlight(*widgets.Paragraph)
	SetStandings(*widgets.Paragraph)
	SetQuarterScores(*widgets.Paragraph)
	SetGameTime(*widgets.Paragraph)
	SetTeamStats(*widgets.Paragraph)
	SetBoxScore(*widgets.Paragraph)
	UpdateScheduleWidget(*widgets.List, []string)
	UpdateHighlightWidget(*widgets.Paragraph, DisplayHighlight)
	UpdateStandingsWidget(*widgets.Paragraph, DisplayStandings)
	UpdateQuarterScoresWidget(*widgets.Paragraph, Game)
	UpdateGameTimeWidget(*widgets.Paragraph, Game)
	UpdateTeamStatsWidget(*widgets.Paragraph, Game, Boxscore)
	UpdateBoxScoreWidget(*widgets.Paragraph, Game, Boxscore)
}

type LayoutService interface {
	NbaLayout
}

type layoutService struct {
	NbaLayout
}

func NewLayoutService() LayoutService {
	return &layoutService{}
}

func (ls *layoutService) Initialise() *Layout {
	l := &Layout{
		Schedule:      widgets.NewList(),
		Highlight:     widgets.NewParagraph(),
		Standings:     widgets.NewParagraph(),
		QuarterScores: widgets.NewParagraph(),
		GameTime:      widgets.NewParagraph(),
		TeamStats:     widgets.NewParagraph(),
		BoxScore:      widgets.NewParagraph(),
	}
	ls.SetSchedule(l.Schedule)
	ls.SetHighlight(l.Highlight)
	ls.SetStandings(l.Standings)
	ls.SetQuarterScores(l.QuarterScores)
	ls.SetGameTime(l.GameTime)
	ls.SetTeamStats(l.TeamStats)
	ls.SetBoxScore(l.BoxScore)
	return l
}

func (ls *layoutService) GetDisplayHighlight(g Game) DisplayHighlight {
	var displayHighlight DisplayHighlight
	displayHighlight.Versus = fmt.Sprintf("%s vs %s\n", g.HTeam.TriCode, g.VTeam.TriCode)
	displayHighlight.Status = fmt.Sprintf(" Status    : %s\n", gameStatus[g.StatusNum])
	displayHighlight.Home = fmt.Sprintf(" Home      : %s\n", g.HTeam.TriCode)
	displayHighlight.Location = fmt.Sprintf(" Location  : %s, %s, %s\n", g.Arena.Name, g.Arena.City, g.Arena.StateAbbr)
	hScore, _ := strconv.Atoi(g.HTeam.Score)
	vScore, _ := strconv.Atoi(g.VTeam.Score)
	if gameStatus[g.StatusNum] == "Finished" {
		if hScore > vScore {
			displayHighlight.Result = fmt.Sprintf(" %s win   : %s - %s\n", g.HTeam.TriCode, g.HTeam.Score, g.VTeam.Score)
		} else {
			displayHighlight.Result = fmt.Sprintf(" %s win   : %s - %s\n", g.VTeam.TriCode, g.VTeam.Score, g.HTeam.Score)
		}
	} else if gameStatus[g.StatusNum] == "In Progress" {
		if hScore == vScore {
			displayHighlight.Result = fmt.Sprintf(" Game Tied : %s - %s\n", g.HTeam.Score, g.VTeam.Score)
		} else if hScore > vScore {
			displayHighlight.Result = fmt.Sprintf(" %s leads : %s - %s\n", g.HTeam.TriCode, g.HTeam.Score, g.VTeam.Score)
		} else {
			displayHighlight.Result = fmt.Sprintf(" %s leads : %s - %s\n", g.VTeam.TriCode, g.VTeam.Score, g.HTeam.Score)
		}
	}
	highlights := strings.Split(g.Nugget.Text, "|")
	var h string
	for i, each := range highlights {
		if i == 0 {
			h = fmt.Sprintf(" Highlight : %s\n", strings.TrimSpace(each))
		} else {
			h = fmt.Sprintf("%s%s%s\n", h, strings.Repeat(" ", 12), each)
		}
	}

	displayHighlight.Highlight = fmt.Sprintf(h)
	return displayHighlight
}

func (ls *layoutService) GetDisplayStandings(g Game) DisplayStandings {
	var displayStandings DisplayStandings
	displayStandings.Header = "Team  Win  Loss\n"
	displayStandings.HomeTeam = fmt.Sprintf(" %s   %s   %s\n", g.HTeam.TriCode, g.HTeam.Win, g.HTeam.Loss)
	displayStandings.VisitorTeam = fmt.Sprintf(" %s   %s   %s", g.VTeam.TriCode, g.VTeam.Win, g.VTeam.Loss)
	return displayStandings
}

func (ls *layoutService) SetSchedule(s *widgets.List) {
	s.Title = "Today's Games"
	s.TitleStyle.Fg = ui.ColorClear
	s.TextStyle.Fg = ui.ColorWhite
	s.SelectedRowStyle.Fg = ui.ColorRed
	s.SelectedRowStyle.Bg = ui.ColorWhite
	s.BorderStyle.Fg = ui.ColorCyan
	s.SetRect(0, 0, 20, 40)
}

func (ls *layoutService) SetHighlight(h *widgets.Paragraph) {
	h.Title = "Highlights"
	h.TitleStyle.Fg = ui.ColorClear
	h.TextStyle.Fg = ui.ColorWhite
	h.BorderStyle.Fg = ui.ColorCyan
	h.SetRect(21, 0, 90, 11)
}

func (ls *layoutService) SetStandings(s *widgets.Paragraph) {
	s.Title = "Standings"
	s.TextStyle.Fg = ui.ColorWhite
	s.BorderStyle.Fg = ui.ColorCyan
	s.SetRect(91, 0, 110, 11)
}

func (ls *layoutService) SetQuarterScores(q *widgets.Paragraph) {
	q.Title = "Quarter Scores"
	q.TextStyle.Fg = ui.ColorWhite
	q.BorderStyle.Fg = ui.ColorCyan
	q.SetRect(21, 11, 90, 17)
}

func (ls *layoutService) SetGameTime(g *widgets.Paragraph) {
	g.Title = "Game Time"
	g.TextStyle.Fg = ui.ColorWhite
	g.BorderStyle.Fg = ui.ColorCyan
	g.SetRect(91, 11, 110, 17)
}

func (ls *layoutService) SetTeamStats(t *widgets.Paragraph) {
	t.Title = "Team Statistics"
	t.BorderStyle.Fg = ui.ColorCyan
	t.SetRect(21, 40, 110, 17)
}

func (ls *layoutService) SetBoxScore(b *widgets.Paragraph) {
	b.Title = "Boxscore"
	b.BorderStyle.Fg = ui.ColorCyan
	b.SetRect(111, 0, 220, 40)
}

func (ls *layoutService) UpdateScheduleWidget(sc *widgets.List, s []string) {
	sc.Rows = s
}

func (ls *layoutService) UpdateHighlightWidget(h *widgets.Paragraph, d DisplayHighlight) {
	h.Text = d.Versus + "\n"
	h.Text = fmt.Sprintf("%s%s", h.Text, d.Status)
	h.Text = fmt.Sprintf("%s%s", h.Text, d.Home)
	h.Text = fmt.Sprintf("%s%s", h.Text, d.Location)
	h.Text = fmt.Sprintf("%s%s", h.Text, d.Result)
	h.Text = fmt.Sprintf("%s%s", h.Text, d.Highlight)
}

func (ls *layoutService) UpdateStandingsWidget(s *widgets.Paragraph, d DisplayStandings) {
	s.Text = d.Header + "\n"
	s.Text = fmt.Sprintf("%s%s", s.Text, d.HomeTeam)
	s.Text = fmt.Sprintf("%s%s", s.Text, d.VisitorTeam)
}

func (ls *layoutService) UpdateQuarterScoresWidget(q *widgets.Paragraph, g Game) {
	q.Text = "      1st  2nd  3rd  4th\n"
	q.Text = fmt.Sprintf("%s %s", q.Text, g.HTeam.TriCode)
	for _, score := range g.HTeam.LineScore {
		q.Text = fmt.Sprintf("%s  %s", q.Text, addPrefixChar(score.Score, 3))
	}
	q.Text = fmt.Sprintf("%s\n %s", q.Text, g.VTeam.TriCode)
	for _, score := range g.VTeam.LineScore {
		q.Text = fmt.Sprintf("%s  %s", q.Text, addPrefixChar(score.Score, 3))
	}
}

func (ls *layoutService) UpdateGameTimeWidget(gt *widgets.Paragraph, g Game) {
	if gameStatus[g.StatusNum] == "In Progress" {
		if g.Period.IsHalfTime {
			gt.Text = fmt.Sprintf("\n Halftime!")
		} else {
			gt.Text = fmt.Sprintf("\n Period: %s\n", strconv.Itoa(g.Period.Current))
			gt.Text = fmt.Sprintf("%s Time  : %s", gt.Text, g.Clock)
		}
	} else if gameStatus[g.StatusNum] == "Not Started" {
		usGameTime, err := time.Parse(time.RFC3339, g.StartTimeUTC)
		if err != nil {
			gt.Text = fmt.Sprintf("\n%s", err)
		}
		localGameTime := usGameTime.In(time.Local)
		gt.Text = fmt.Sprintf("\n Starts at %s:%s\n", addPrefixChar(strconv.Itoa(localGameTime.Hour()), 2, "0"), addPrefixChar(strconv.Itoa(localGameTime.Minute()), 2, "0"))
	} else {
		gt.Text = "\n Finished!"
	}
}

func (ls *layoutService) UpdateTeamStatsWidget(ts *widgets.Paragraph, g Game, b Boxscore) {
	ts.Text = fmt.Sprintf("\n%s%s%s      Team Stats      %s%s\n\n", tenSpace, addPrefixChar(g.HTeam.TriCode, 6), fifteenSpace, fifteenSpace, g.VTeam.TriCode)
	ts.Text = fmt.Sprintf("%s%s%s\n", ts.Text, tenSpace, strings.Repeat("=", 65))
	ts.Text = fmt.Sprintf("%s%s%s%s        Points        %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.Totals.Points, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.Totals.Points, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s       Rebounds       %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.Totals.TotReb, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.Totals.TotReb, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s       Assists        %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.Totals.Assists, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.Totals.Assists, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s        Blocks        %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.Totals.Blocks, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.Totals.Blocks, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s      Field Goal      %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.Totals.FGM+"/"+b.Stats.HTeam.Totals.FGA, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.Totals.FGM+"/"+b.Stats.VTeam.Totals.FGA, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s         FG %%         %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.Totals.FGP, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.Totals.FGP, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s   3 Pt Field Goal    %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.Totals.TPM+"/"+b.Stats.HTeam.Totals.TPA, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.Totals.TPM+"/"+b.Stats.VTeam.Totals.TPA, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s      3 Pt FG %%       %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.Totals.TPP, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.Totals.TPP, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s      Free Throw      %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.Totals.FTM+"/"+b.Stats.HTeam.Totals.FTA, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.Totals.FTM+"/"+b.Stats.VTeam.Totals.FTA, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s         FT %%         %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.Totals.FTP, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.Totals.FTP, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s  Fast Break Points   %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.FastBreakPoints, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.FastBreakPoints, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s   Points In Paint    %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.PointsInPaint, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.PointsInPaint, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s    Biggest Lead      %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.BiggestLead, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.BiggestLead, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s Second Chance Points %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.SecondChancePoints, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.SecondChancePoints, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s Points Off Turnovers %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.PointsOffTurnovers, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.PointsOffTurnovers, 6))
	ts.Text = fmt.Sprintf("%s%s%s%s     Longest Run      %s%s\n", ts.Text, tenSpace, addPrefixChar(b.Stats.HTeam.LongestRun, 6), fifteenSpace, fifteenSpace, addSuffixSpace(b.Stats.VTeam.LongestRun, 6))
}

func (ls *layoutService) UpdateBoxScoreWidget(bs *widgets.Paragraph, g Game, b Boxscore) {
	if len(b.Stats.ActivePlayers) == 0 {
		return
	}
	starterCnt := 0
	awayStats := ""
	homeStats := ""
	teamID := b.Stats.ActivePlayers[0].TeamID
	stats := getBoxScoreHeader(g, teamID)
	for _, p := range b.Stats.ActivePlayers {
		if starterCnt == 5 {
			stats = fmt.Sprintf("%s%s\n", stats, strings.Repeat("-", 107))
		}
		if teamID != p.TeamID {
			teamID = p.TeamID
			stats = stats + getBoxScoreTotal(g, b.Stats.VTeam.Totals)
			awayStats = stats
			stats = getBoxScoreHeader(g, teamID)
			starterCnt = 0
		}
		name := addSuffixSpace(fmt.Sprintf("%s %s", p.FirstName, p.LastName), 25)
		if p.DNP == "" {
			in := ""
			if p.IsOnCourt {
				in = "*"
			}
			in = addPrefixChar(in, 4)
			min := addPrefixChar(addPrefixChar(p.Min, 5, "0"), 7)
			fg := addPrefixChar(fmt.Sprintf("%s-%s", p.FGM, p.FGA), 6)
			tg := addPrefixChar(fmt.Sprintf("%s-%s", p.TPM, p.TPA), 6)
			ftg := addPrefixChar(fmt.Sprintf("%s-%s", p.FTM, p.FTA), 6)
			or := addPrefixChar(p.OffReb, 5)
			dr := addPrefixChar(p.DefReb, 5)
			tr := addPrefixChar(p.TotReb, 3)
			a := addPrefixChar(p.Assists, 3)
			s := addPrefixChar(p.Steals, 3)
			b := addPrefixChar(p.Blocks, 3)
			to := addPrefixChar(p.Turnovers, 3)
			pf := addPrefixChar(p.PFouls, 3)
			pm := addPrefixChar(p.PlusMinus, 3)
			pts := addPrefixChar(p.Points, 3)
			stats = fmt.Sprintf("%s%s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s\n", stats, name, in, min, fg, tg, ftg, or, dr, tr, a, s, b, to, pf, pm, pts)
		} else {
			stats = fmt.Sprintf("%s%s %s %s\n", stats, name, strings.Repeat(" ", 6), p.DNP)
		}
		starterCnt++
	}
	stats = stats + getBoxScoreTotal(g, b.Stats.HTeam.Totals)
	homeStats = stats
	bs.Text = fmt.Sprintf("%s\n\n%s", homeStats, awayStats)
}

func getBoxScoreHeader(g Game, teamID string) string {
	team := g.VTeam.TriCode
	if teamID == g.HTeam.TeamID {
		team = g.HTeam.TriCode
	}
	header := fmt.Sprintf("STARTERS - %s                In   MIN     FG    3PT     FT  OREB  DREB REB AST STL BLK  TO  PF +/- PTS\n%s\n", team, strings.Repeat("-", 107))
	return header
}

func getBoxScoreTotal(g Game, t TeamTotals) string {
	fg := addPrefixChar(fmt.Sprintf("%s-%s", t.FGM, t.FGA), 40)
	tg := addPrefixChar(fmt.Sprintf("%s-%s", t.TPM, t.TPA), 6)
	ftg := addPrefixChar(fmt.Sprintf("%s-%s", t.FTM, t.FTA), 6)
	or := addPrefixChar(t.OffReb, 5)
	dr := addPrefixChar(t.DefReb, 5)
	tr := addPrefixChar(t.TotReb, 3)
	a := addPrefixChar(t.Assists, 3)
	s := addPrefixChar(t.Steals, 3)
	b := addPrefixChar(t.Blocks, 3)
	to := addPrefixChar(t.Turnovers, 3)
	pf := addPrefixChar(t.PFouls, 3)
	pts := addPrefixChar(t.Points, 7)
	total := fmt.Sprintf("%s\n", strings.Repeat("=", 107))
	total = fmt.Sprintf("%s%s %s %s %s %s %s %s %s %s %s %s %s %s\n", total, "TEAM", fg, tg, ftg, or, dr, tr, a, s, b, to, pf, pts)
	return total
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
	//return formattedString[:len(formattedString)-length]
	return formattedString[:length]
}
