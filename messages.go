package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleMessage(session *discordgo.Session, message *discordgo.Message) {
	channel, err := session.State.Channel(message.ChannelID)
	if err != nil {
		return
	}

	isDM := channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM

	var member *discordgo.Member
	guild, err := session.State.Guild(channel.GuildID)
	if err != nil {
		if !isDM {
			return
		}
	} else {
		_member, err := session.GuildMember(guild.ID, message.Author.ID)
		if err != nil && !isDM {
			return
		}
		member = _member
	}

	content := message.Content
	if len(content) <= 0 {
		return
	}

	var responseEmbed *discordgo.MessageEmbed
	var responseText string

	prefix := ""
	if strings.HasPrefix(content, Bot.CommandPrefix) {
		prefix = Bot.CommandPrefix
	}

	if prefix != "" {
		Debug.Println("[" + channel.Name + "] " + message.Author.Username + ": " + content)

		commandContent := strings.TrimPrefix(content, prefix)
		command := strings.Split(commandContent, " ")

		responseEmbed, responseText = callCommand(command[0], command[1:], &CommandEnvironment{
			Guild: guild, Channel: channel,
			User: message.Author, Member: member,
			Message: message,
		})
	}

	if responseEmbed != nil {
		_, err := session.ChannelMessageSendEmbed(message.ChannelID, responseEmbed)
		if err != nil {
			Error.Println(err.Error())
		}
	}

	if len(responseText) > 0 {
		Bot.DiscordSession.ChannelMessageSend(channel.ID, responseText)
	}
}
