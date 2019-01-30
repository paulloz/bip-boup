package main

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Function            func([]string, *CommandEnvironment) (*discordgo.MessageEmbed, string)
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
	BotData.Commands["help"] = &Command{
		Function: commandHelp,
		HelpText: "Montre une liste de commande que vous pouvez utiliser ou bien l'aide d'une commande spécifique.",
		Arguments: []CommandArgument{
			{Name: "commande", Description: "Commande dont on veut afficher l'aide", ArgType: "commande"},
		},
	}
	BotData.Commands["ping"] = &Command{Function: commandPing, HelpText: "Retourne le ping moyen vers Discord."}
	BotData.Commands["nightcore"] = &Command{
		Function: commandNightcore,
		HelpText: "Cherche du nightcore sur YouTube.",
		Arguments: []CommandArgument{
			{Name: "requête", Description: "La recherche à faire sur YouTube", ArgType: "string"},
		},
		RequiredArguments: []string{"requête"},
	}
	BotData.Commands["furigana"] = &Command{
		Function: commandFurigana,
		HelpText: "Ajoute des furiganas à un texte en Japonais.",
		Arguments: []CommandArgument{
			{Name: "texte", Description: "Le texte dans lequel on veut insérer des furiganas.", ArgType: "string"},
		},
		RequiredArguments: []string{"texte"},
	}
}

func callCommand(commandName string, args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	if command, exists := BotData.Commands[commandName]; exists {
		if len(command.IsAliasTo) > 0 {
			return callCommand(command.IsAliasTo, args, env)
		}

		if len(args) >= len(command.RequiredArguments) {
			return command.Function(args, env)
		} else {
			return callCommand("help", []string{commandName}, env)
		}
	}

	return nil, ""
}
