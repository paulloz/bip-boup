package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gojp/kana"
	"github.com/ikawaha/kagome/tokenizer"
	"github.com/moovweb/gokogiri"
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
	BotData.Commands["help"] = &Command{Function: commandHelp, HelpText: "Montre une liste de commande que vous pouvez utiliser."}
	BotData.Commands["ping"] = &Command{Function: commandPing, HelpText: "Retourne le ping moyen vers Discord."}
	BotData.Commands["nightcore"] = &Command{
		Function: commandNightcore,
		HelpText: "Cherche du nightcore sur YouTube.",
		Arguments: []CommandArgument{
			{Name: "recherche", Description: "La recherche à faire sur YouTube", ArgType: "string"},
		},
	}
	BotData.Commands["furigana"] = &Command{
		Function: commandFurigana,
		HelpText: "Ajoute des furiganas à un texte en Japonais.",
		Arguments: []CommandArgument{
			{Name: "texte", Description: "Le texte dans lequel on veut les furiganas", ArgType: "string"},
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
		}
	}

	return nil, ""
}

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

func commandNightcore(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	if len(args) < 1 {
		return nil, ""
	}

	resp, err := http.Get(fmt.Sprintf("%s&search_query=nightcore+%s", "https://www.youtube.com/results?sp=EgIQAQ%253D%253D", strings.Join(args, "+")))
	if err != nil {
		return nil, ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ""
	}

	doc, err := gokogiri.ParseHtml(body)
	if err != nil {
		return nil, ""
	}
	defer doc.Free()

	rootNode := doc.Root()
	resultsNode, _ := rootNode.Search("//h3/a")
	for _, node := range resultsNode {
		if strings.Contains(strings.ToLower(node.Content()), "nightcore") {
			return nil, fmt.Sprintf("%s - https://www.youtube.com%s", node.Content(), node.Attr("href"))
		}
	}

	return nil, ""
}

func commandFurigana(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	text := strings.Join(args, " ")
	if len(text) <= 0 {
		return nil, ""
	}

	t := tokenizer.New()
	tokens := t.Tokenize(text)
	tokenized := make([]string, len(tokens))
	for _, token := range tokens {
		switch token.Class {
		case tokenizer.UNKNOWN:
			tokenized = append(tokenized, token.Surface)
		case tokenizer.KNOWN:
			features := token.Features()
			katakanas := features[len(features)-2]
			if kana.IsHiragana(token.Surface) {
				tokenized = append(tokenized, token.Surface)
			} else {
				value := fmt.Sprintf("%s(%s)", token.Surface, kana.RomajiToHiragana(kana.KanaToRomaji(katakanas)))
				tokenized = append(tokenized, value)
			}
		}
	}

	return &discordgo.MessageEmbed{Title: strings.Join(tokenized, " "), Color: 0xffffff}, ""
}
