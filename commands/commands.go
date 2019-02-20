package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/paulloz/bip-boup/bot"
)

var (
	commands map[string]*bot.Command = make(map[string]*bot.Command)
)

func InitCommands() map[string]*bot.Command {
	return commands
}

func CallCommand(commandName string, args []string, env *bot.CommandEnvironment, bot *bot.Bot) (*discordgo.MessageEmbed, string) {
	if command, exists := bot.Commands[commandName]; exists {
		if len(command.IsAliasTo) > 0 {
			return CallCommand(command.IsAliasTo, args, env, bot)
		}

		if !command.IsAdmin || bot.IsUserAdmin(env.User) {
			if len(args) >= len(command.RequiredArguments) {
				return command.Function(args, env, bot)
			}
			return CallCommand("help", []string{commandName}, env, bot)
		}
	}

	return nil, ""
}
