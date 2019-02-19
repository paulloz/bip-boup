package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/paulloz/bip-boup/bot"
)

func InitCommands() map[string]*bot.Command {
	cmds := make(map[string]*bot.Command)

	cmds["restart"] = &bot.Command{
		Function: commandRestart, IsAdmin: true,
		HelpText: "Arrête et redémarre le bot.",
	}
	cmds["update"] = &bot.Command{
		Function: commandUpdate, IsAdmin: true,
		HelpText: "Met à jour et redémarre le bot.",
	}
	cmds["queue"] = &bot.Command{
		Function: commandQueue, IsAdmin: true,
		HelpText: "Renvoie le contenu de la queue.",
	}

	cmds["help"] = &bot.Command{
		Function: commandHelp,
		HelpText: "Montre une liste de commande que vous pouvez utiliser ou bien l'aide d'une commande spécifique.",
		Arguments: []bot.CommandArgument{
			{Name: "commande", Description: "Commande dont on veut afficher l'aide", ArgType: "commande"},
		},
	}
	cmds["?"] = &bot.Command{IsAliasTo: "help"}

	cmds["ping"] = &bot.Command{Function: commandPing, HelpText: "Retourne le ping moyen vers Discord."}

	cmds["nightcore"] = &bot.Command{
		Function: commandNightcore, HelpText: "Cherche du nightcore sur YouTube.",
		Arguments: []bot.CommandArgument{
			{Name: "requête", Description: "La recherche à faire sur YouTube", ArgType: "string"},
		},
		RequiredArguments: []string{"requête"},
	}

	cmds["furigana"] = &bot.Command{
		Function: commandFurigana, HelpText: "Ajoute des furiganas à un texte en Japonais.",
		Arguments: []bot.CommandArgument{
			{Name: "texte", Description: "Le texte dans lequel on veut insérer des furiganas", ArgType: "string"},
		},
		RequiredArguments: []string{"texte"},
	}

	cmds["directan"] = &bot.Command{Function: commandDirectAN, HelpText: "Envoie un lien vers la séance publique en cours à l'Assemblée Nationale."}

	cmds["député"] = &bot.Command{
		Function: commandDepute, HelpText: "Montre les informations à propos d'un député disponibles sur NosDéputés.fr.",
		Arguments: []bot.CommandArgument{
			{Name: "prénom", Description: "Le prénom du député", ArgType: "string"},
			{Name: "nom", Description: "Le nom du député", ArgType: "string"},
		},
		RequiredArguments: []string{"nom"},
	}
	cmds["depute"] = &bot.Command{IsAliasTo: "député"}

	cmds["hltb"] = &bot.Command{
		Function: commandHLTB, HelpText: "Donne le temps moyen pour finir un jeu.",
		Arguments: []bot.CommandArgument{
			{Name: "nom du jeu", Description: "La recherche à faire sur howlongtobeat", ArgType: "string"},
		},
		RequiredArguments: []string{"nom du jeu"},
	}

	cmds["rappel"] = &bot.Command{
		Function: commandReminder,
		HelpText: "Enverra un message de rappel à l'heure indiquée",
		Arguments: []bot.CommandArgument{
			{Name: "requête", Description: "Une requête en langage naturel qui sera plus ou moins bien interprétée", ArgType: "string"},
		},
		RequiredArguments: []string{"requête"},
	}

	cmds["8ball"] = &bot.Command{
		Function: command8Ball, HelpText: "Aide à choisir dans les moments difficiles.",
		Arguments: []bot.CommandArgument{
			{Name: "une question", Description: "Une interrogation sur laquelle vous souhaitez une réponse", ArgType: "string"},
		},
		RequiredArguments: []string{"une question"},
	}

	return cmds
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
