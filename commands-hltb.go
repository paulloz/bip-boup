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

	doc, err := httpPostAsHTML(url, params)
	if err != nil {
		return nil, ""
	}
	defer doc.Free()

	imageResult := xpath.Compile("//div[@class='search_list_image']/a/img")
	images, _ := doc.Root().Search(imageResult)
	titlesResult := xpath.Compile("//div[@class='search_list_details']/h3")
	titles, _ := doc.Root().Search(titlesResult)
	gameResult := xpath.Compile("//div[@class='search_list_details_block']/div")
	games, _ := doc.Root().Search(gameResult)

	var image string
	for _, img := range images {
		image = img.Attr("src")
		break
	}

	var title string
	for _, titleFound := range titles {
		title = titleFound.Content()
		break
	}

	if title == "" {
		return nil, fmt.Sprintf("Aucune correspondance trouvée pour `%s`.", strings.Join(args, " "))
	}

	var content []string
	for _, game := range games {
		tmp := strings.Split(game.Content(), "\n")
		content = tmp[1:7]
		break
	}

	return formatResult(title, image, content, args)
}

func formatResult(title string, imageURL string, content []string, args []string) (*discordgo.MessageEmbed, string) {
	fields := []*discordgo.MessageEmbedField{}
	for i := 0; i < len(content); i += 2 {
		fields = append(fields, embedField(content[i], formatTime(content[i+1]), true))
	}

	return &discordgo.MessageEmbed{
		Title:       title,
		Description: "",
		Image:       &discordgo.MessageEmbedImage{URL: imageURL},
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Source : HowLongToBeat.com",
		},
	}, ""

}

func formatTime(time string) string {
	if time == "--" {
		return "N/A"
	}
	return time
}
