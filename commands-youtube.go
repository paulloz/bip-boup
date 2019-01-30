package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/moovweb/gokogiri"
)

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
