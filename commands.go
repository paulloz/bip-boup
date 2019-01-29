package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Function            func([]string, *CommandEnvironment) *discordgo.MessageEmbed
	HelpText            string
	Arguments           []CommandArgument
	RequiredArguments   []string
	RequiredPermissions int

	IsAliasTo string
}

type CommandArgument struct {
	Name        string
	ArgType     string
	Description string
}

type CommandEnvironment struct {
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	User    *discordgo.User
	Member  *discordgo.Member

	Message *discordgo.Message
}

func initCommands() {
	BotData.Commands = make(map[string]*Command)

	BotData.Commands["?"] = &Command{IsAliasTo: "help"}
	BotData.Commands["help"] = &Command{Function: commandHelp, HelpText: "Montre une liste de commande que vous pouvez utiliser."}
	BotData.Commands["ping"] = &Command{Function: commandPing, HelpText: "Retourne le ping moyen vers Discord."}
}

func callCommand(commandName string, args []string, env *CommandEnvironment) *discordgo.MessageEmbed {
	if command, exists := BotData.Commands[commandName]; exists {
		if len(command.IsAliasTo) > 0 {
			return callCommand(command.IsAliasTo, args, env)
		}
		if len(args) >= len(command.RequiredArguments) {
			return command.Function(args, env)
		}
	}

	return nil
}

func commandHelp(args []string, env *CommandEnvironment) *discordgo.MessageEmbed {
	// Get all commands
	var commands []string
	for command := range BotData.Commands {
		commands = append(commands, command)
	}

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
	}
}

func commandPing(args []string, env *CommandEnvironment) *discordgo.MessageEmbed {
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

	return &discordgo.MessageEmbed{
		Title:       "Pong!",
		Description: fmt.Sprintf("Le ping moyen est de ``%dms``. Un total de ``%d/%d`` paquets ont été perdus.\n", pingAverage, failCount, len(pingResults)),
	}
}
