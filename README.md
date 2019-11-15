# Dislay Today's NBA Results
This will display today's NBA results using go http request. It uses the endpoint `http://data.nba.net`
There are 3 main widgets that display information about the game.
- Top Left widget - display today's schedule that can be highlighted using up and down arrow.
- Top Right Widget - display the summary of the game that is currently chosen in the top left widget.
- Bottom Widget - display the boxscore of the game (WIP)

## Run
Execute `go run main.go nba.go`

## Build
Execute `go build main.go nba.go`

## Library
This is using [termui](github.com/gizak/termui/v3) to render the widget and its details.