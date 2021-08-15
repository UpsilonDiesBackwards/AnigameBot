package main

import "github.com/bwmarrin/discordgo"

var context appContext

type appContext struct {
	config      *config
	discord     *discordgo.Session
	gameContext *gameContext
}

type gameContext struct {
	currentLocation int
	currentFloor    int
}
