package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/paulloz/bip-boup/bot"
	"github.com/paulloz/bip-boup/httpreq"
)

func commandNightcore(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	doc, err := httpreq.HTTPGetAsHTML(fmt.Sprintf("%s&search_query=nightcore+%s", "https://www.youtube.com/results?sp=EgIQAQ%253D%253D", strings.Join(args, "+")))
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
