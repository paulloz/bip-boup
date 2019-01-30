package main

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func commandRestart(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	Info.Println("Received a restart command, exiting...")

	os.Exit(0) // Exit, the master process will start a new bot

	return nil, ""
}
