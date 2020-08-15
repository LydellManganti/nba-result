package main

import (
	"nba-result/controllers"
	"nba-result/models"
)

func main() {
	ls := models.NewLayoutService()
	ns := models.NewNbaService()
	l := ls.Initialise()
	controllers.Render(l, ls, ns)
}
