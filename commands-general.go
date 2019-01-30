package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
)

func commandHelp(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	// Get all commands
	var commands []string
	for command := range BotData.Commands {
		commands = append(commands, command)
	}
	sort.Strings(commands)

	fields := []*discordgo.MessageEmbedField{}
	for _, commandName := range commands {
		command := BotData.Commands[commandName]

		if len(command.IsAliasTo) > 0 {
			continue
		}

		if command.RequiredPermissions != 0 {
			continue
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   BotData.CommandPrefix + commandName,
			Value:  command.HelpText,
			Inline: true,
		})
	}

	return &discordgo.MessageEmbed{
		Title:  "Liste des commandes utilisables",
		Fields: fields,
	}, ""
}

func commandPing(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	pingResults := make([]int, 4)

	// Perform the pings
	pingEmbed := &discordgo.MessageEmbed{Title: "Ping!"}
	for i := 0; i < len(pingResults); i++ {
		currentTime := int(time.Now().UnixNano() / 1000000)

		ping, err := BotData.DiscordSession.ChannelMessageSendEmbed(env.Channel.ID, pingEmbed)
		if err != nil {
			pingResults[i] = -1
			continue
		}

		newTime := int(time.Now().UnixNano() / 1000000)

		BotData.DiscordSession.ChannelMessageDelete(env.Channel.ID, ping.ID)

		pingResults[i] = newTime - currentTime
	}

	// Average the results
	pingSum := 0
	failCount := 0
	for i := 0; i < len(pingResults); i++ {
		if pingResults[i] == -1 {
			failCount++
		} else {
			pingSum += pingResults[i]
		}
	}
	pingAverage := int(pingSum / len(pingResults))

	color := 0x8b0000
	if pingAverage < 10 {
		color = 0x90ee90
	} else if pingAverage < 50 {
		color = 0xeead00
	} else if pingAverage < 100 {
		color = 0xda8600
	} else if pingAverage < 150 {
		color = 0xc26001
	} else if pingAverage < 200 {
		color = 0xa73902
	}

	return &discordgo.MessageEmbed{
		Title: "Pong !", Color: color,
		Description: fmt.Sprintf("Le ping moyen est de ``%dms``. Un total de ``%d/%d`` paquets ont été perdus.\n", pingAverage, failCount, len(pingResults)),
	}, ""
}
