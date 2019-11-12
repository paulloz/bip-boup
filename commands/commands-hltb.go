package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/moovweb/gokogiri/xpath"

	"github.com/paulloz/bip-boup/bot"
	"github.com/paulloz/bip-boup/embed"
	"github.com/paulloz/bip-boup/httpreq"
)

func commandHLTB(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	baseURL := "https://howlongtobeat.com/"
	url := fmt.Sprintf("%ssearch_results.php?page=1", baseURL)
	params := fmt.Sprintf("queryString=%s&t=games&sorthead=popular&sortd=Normal Order&plat=&length_type=main&length_min=&length_max=&detail=", strings.Join(args, " "))
	doc, err := httpreq.HTTPPostAsHTML(url, params)
	if err != nil {
		return nil, ""
	}
	defer doc.Free()

	imageResult := xpath.Compile("//div[@class='search_list_image']/a/img")
	images, _ := doc.Root().Search(imageResult)
	titlesResult := xpath.Compile("//div[@class='search_list_details']/h3")
	titles, _ := doc.Root().Search(titlesResult)
	gameResult := xpath.Compile("//div[@class='search_list_details_block']")
	games, _ := doc.Root().Search(gameResult)

	var image string
	for _, img := range images {
		image = fmt.Sprintf("%s%s", baseURL, img.Attr("src"))
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
		content = parseGameResult(strings.Split(game.Content(), "\n"))
		break
	}

	return formatResult(title, image, content, args)
}

func parseGameResult(gameResult []string) []string {
	var res []string
	for _, info := range gameResult {
		if info == " " {
			continue
		}
		res = append(res, strings.Split(info, "  ")...)
	}
	return res
}

func formatResult(title string, imageURL string, content []string, args []string) (*discordgo.MessageEmbed, string) {
	fields := []*discordgo.MessageEmbedField{}
	for i := 0; i < len(content)-1; i += 2 {
		fields = append(fields, embed.EmbedField(content[i], formatTime(content[i+1]), true))
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

func init() {
	commands["hltb"] = &bot.Command{
		Function: commandHLTB, HelpText: "Donne le temps moyen pour finir un jeu.",
		Arguments: []bot.CommandArgument{
			{Name: "nom du jeu", Description: "La recherche à faire sur howlongtobeat", ArgType: "string"},
		},
		RequiredArguments: []string{"nom du jeu"},
	}
}
