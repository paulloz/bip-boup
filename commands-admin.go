package main

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func commandRestart(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	Info.Println("Received a restart command, exiting...")

	Bot.DiscordSession.ChannelMessageSend(env.Channel.ID, "Je reviens dans un instant...")
	fileHandler, err := os.Create("/tmp/bip-boup.restart")
	if err == nil {
		fileHandler.Write([]byte(env.Channel.ID))
		fileHandler.Close()
	}

	os.Exit(0) // Exit, the master process will start a new bot
	return nil, ""
}
