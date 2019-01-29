package main

import (
	"github.com/bwmarrin/discordgo"
)

func discordMessageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
	message, err := session.ChannelMessage(event.ChannelID, event.ID)
	if err != nil {
		return
	}
	if message.Author.ID == session.State.User.ID {
		return
	}

	go handleMessage(session, message)
}
