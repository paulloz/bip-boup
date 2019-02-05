package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/moovweb/gokogiri/xpath"
)

func commandHLTB(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	url := "https://howlongtobeat.com/search_results.php?page=1"
	params := fmt.Sprintf("queryString=%s&t=games&sorthead=popular&sortd=Normal Order&plat=&length_type=main&length_min=&length_max=&detail=", strings.Join(args, " "))

	doc, err := httpPostAsHtml(url, params)
	if err != nil {
		return nil, ""
	}
	defer doc.Free()

	titlesResult := xpath.Compile("//div[@class='search_list_details']/h3")
	titles, _ := doc.Root().Search(titlesResult)
	gameResult := xpath.Compile("//div[@class='search_list_details_block']/div")
	games, _ := doc.Root().Search(gameResult)

	var title string
	for _, titleFound := range titles {
		title = titleFound.Content()
		break
	}

	var content []string
	for _, game := range games {
		tmp := strings.Split(game.Content(), "\n")
		content = []string{tmp[2], tmp[6]}
		break
	}

	return nil, formatResult(title, content, args)
}

func formatResult(title string, times []string, args []string) string {
	if title == "" {
		return fmt.Sprintf("Aucun résultat pour %s.\n", strings.Join(args, " "))
	}
	result := strings.Trim(title, "\n")
	if times[0] == "--" {
		result += fmt.Sprintf(" n'a pas encore de temps renseigné")
	} else {
		result += fmt.Sprintf(" se finit en %s si tu traces", formatTime(times[0]))
		if times[1] != "--" {
			result += fmt.Sprintf(", %s si tu veux tout faire", formatTime(times[1]))
		}
	}
	return fmt.Sprintf("%s.", result)
}

func formatTime(time string) string {
	tmp := strings.Split(time, " ")
	if tmp[1] == "Hours" {
		return fmt.Sprintf("%s %s", tmp[0], "h")
	}
	if tmp[1] == "Mins" {
		return fmt.Sprintf("%s %s", tmp[0], "m")
	}
	return ""
}
