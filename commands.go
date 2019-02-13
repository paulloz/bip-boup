package main

import (
	"github.com/bwmarrin/discordgo"
)

// Command ...
type Command struct {
	Function          func([]string, *CommandEnvironment) (*discordgo.MessageEmbed, string)
	HelpText          string
	Arguments         []CommandArgument
	RequiredArguments []string

	RequiredPermissions int
	IsAdmin             bool

	IsAliasTo string
}

// CommandArgument ...
type CommandArgument struct {
	Name        string
	ArgType     string
	Description string
}

// CommandEnvironment ...
type CommandEnvironment struct {
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	User    *discordgo.User
	Member  *discordgo.Member

	Message *discordgo.Message
}

func initCommands() {
	Bot.Commands = make(map[string]*Command)

	Bot.Commands["restart"] = &Command{
		Function: commandRestart, IsAdmin: true,
		HelpText: "Arrête et redémarre le bot.",
	}
	Bot.Commands["update"] = &Command{
		Function: commandUpdate, IsAdmin: true,
		HelpText: "Met à jour et redémarre le bot.",
	}
	Bot.Commands["queue"] = &Command{
		Function: commandQueue, IsAdmin: true,
		HelpText: "Renvoie le contenu de la queue.",
	}

	Bot.Commands["help"] = &Command{
		Function: commandHelp,
		HelpText: "Montre une liste de commande que vous pouvez utiliser ou bien l'aide d'une commande spécifique.",
		Arguments: []CommandArgument{
			{Name: "commande", Description: "Commande dont on veut afficher l'aide", ArgType: "commande"},
		},
	}
	Bot.Commands["?"] = &Command{IsAliasTo: "help"}

	Bot.Commands["ping"] = &Command{Function: commandPing, HelpText: "Retourne le ping moyen vers Discord."}

	Bot.Commands["nightcore"] = &Command{
		Function: commandNightcore, HelpText: "Cherche du nightcore sur YouTube.",
		Arguments: []CommandArgument{
			{Name: "requête", Description: "La recherche à faire sur YouTube", ArgType: "string"},
		},
		RequiredArguments: []string{"requête"},
	}

	Bot.Commands["furigana"] = &Command{
		Function: commandFurigana, HelpText: "Ajoute des furiganas à un texte en Japonais.",
		Arguments: []CommandArgument{
			{Name: "texte", Description: "Le texte dans lequel on veut insérer des furiganas", ArgType: "string"},
		},
		RequiredArguments: []string{"texte"},
	}

	Bot.Commands["directan"] = &Command{Function: commandDirectAN, HelpText: "Envoie un lien vers la séance publique en cours à l'Assemblée Nationale."}

	Bot.Commands["député"] = &Command{
		Function: commandDepute, HelpText: "Montre les informations à propos d'un député disponibles sur NosDéputés.fr.",
		Arguments: []CommandArgument{
			{Name: "prénom", Description: "Le prénom du député", ArgType: "string"},
			{Name: "nom", Description: "Le nom du député", ArgType: "string"},
		},
		RequiredArguments: []string{"nom"},
	}
	Bot.Commands["depute"] = &Command{IsAliasTo: "député"}

	Bot.Commands["hltb"] = &Command{
		Function: commandHLTB, HelpText: "Donne le temps moyen pour finir un jeu.",
		Arguments: []CommandArgument{
			{Name: "nom du jeu", Description: "La recherche à faire sur howlongtobeat", ArgType: "string"},
		},
		RequiredArguments: []string{"nom du jeu"},
	}

	Bot.Commands["rappel"] = &Command{
		Function: commandReminder,
		HelpText: "Enverra un message de rappel à l'heure indiquée",
		Arguments: []CommandArgument{
			{Name: "requête", Description: "Une requête en langage naturel qui sera plus ou moins bien interprétée", ArgType: "string"},
		},
		RequiredArguments: []string{"requête"},
	}
}

func callCommand(commandName string, args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	if command, exists := Bot.Commands[commandName]; exists {
		if len(command.IsAliasTo) > 0 {
			return callCommand(command.IsAliasTo, args, env)
		}

		if !command.IsAdmin || isUserAdmin(env.User) {
			if len(args) >= len(command.RequiredArguments) {
				return command.Function(args, env)
			}
			return callCommand("help", []string{commandName}, env)
		}
	}

	return nil, ""
}
