package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/moovweb/gokogiri/xpath"

	"github.com/paulloz/bip-boup/bot"
	"github.com/paulloz/bip-boup/httpreq"
)

func commandUrban(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	request := strings.Join(args, " ")

	resp := &discordgo.MessageEmbed{
		Title:       request,
		Description: "Aucun résultat.",
		Footer:      &discordgo.MessageEmbedFooter{Text: "Source: https://www.urbandictionary.com/"},
	}

	doc, err := httpreq.HTTPGetAsHTML(fmt.Sprintf("https://www.urbandictionary.com/define.php?term=%s", strings.Replace(request, " ", "%20", -1)))
	if err != nil {
		return resp, ""
	}
	defer doc.Free()

	defs, _ := doc.Root().Search(xpath.Compile("//div[contains(@class, 'def-panel')]"))
	if len(defs) > 0 {
		meaning, _ := defs[0].Search(xpath.Compile("div[contains(@class, 'meaning')]"))

		resp.Description = meaning[0].Content()
	}

	return resp, ""
}

func init() {
	commands["urban"] = &bot.Command{
		Function: commandUrban, HelpText: "Cherche une définition sur Urban Dictionary.",
		Arguments: []bot.CommandArgument{
			{Name: "requête", Description: "La définition à chercher", ArgType: "string"},
		},
		RequiredArguments: []string{"requête"},
	}
}
