# Display Today's NBA Results using Golang
This will display today's NBA results using [golang](https://golang.org) http request. It uses the endpoint `http://data.nba.net`
There are widgets that display information about the game.
- Today's Game widget - display today's schedule that can be highlighted using up and down arrow.
- Highlights Widget - display the summary of the game that is currently chosen in the top left widget.
- Standings Widget - display the standing of the teams that is playing.
- Quarter Scores Widget - display the team scoring for each quarter.
- Game Time Widget - display time remaining if the game is in progress.
- Team Statistics Widget - display the team's game statistics.
- Boxscore Widget - display the player's game statistics

## Run
Execute `go run main.go`

## Build
Execute `go build main.go`

## Library
This is using [termui](github.com/gizak/termui/v3) to render the widget and its details.