# Display Today's NBA Results using Golang
This will display today's NBA results using [golang](https://golang.org) http request. It uses the endpoint `http://data.nba.net`
There are 4 main widgets that display information about the game.
- Top Left widget - display today's schedule that can be highlighted using up and down arrow.
- Top Middle Widget - display the summary of the game that is currently chosen in the top left widget.
- Top Right Widget - display the standing of the teams that is playing.
- Bottom Widget - display the Team stats of the game.

## Run
Execute `go run main.go nba.go`

## Build
Execute `go build main.go nba.go`

## Library
This is using [termui](github.com/gizak/termui/v3) to render the widget and its details.