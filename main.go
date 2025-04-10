package main

import (
	"github.com/charmbracelet/log"
	"github.com/itspacchu/anilist-chart/api"
	"github.com/itspacchu/anilist-chart/processing"
)

func main() {
	// log.SetLevel(log.DebugLevel)
	processing.InitFont()
	if err := api.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
