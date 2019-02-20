package commands

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/moovweb/gokogiri/xpath"

	"github.com/paulloz/bip-boup/bot"
	"github.com/paulloz/bip-boup/embed"
	"github.com/paulloz/bip-boup/httpreq"
)

func commandTomatoes(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	resp := &discordgo.MessageEmbed{Title: "Rotten Tomatoes", Color: 0xa73902}

	request := strings.Join(args, " ")

	doc, err := httpreq.HTTPGetAsHTML(fmt.Sprintf("https://www.rottentomatoes.com/search/?search=%s", strings.Replace(request, " ", "%20", -1)))
	if err != nil {
		return nil, ""
	}
	defer doc.Free()

	results, err := doc.Root().Search(xpath.Compile("//div[@id='search-results-root']/following-sibling::script"))
	if err != nil {
		return nil, ""
	}

	if len(results) <= 0 {
		resp.Description = fmt.Sprintf("Je n'ai trouvé aucun résultat pour **%s** :(", request)
		return resp, ""
	}

	scriptContent := results[0].Content()
	startIndex := strings.Index(scriptContent, request) + len(request) + 3
	endIndex := strings.LastIndex(scriptContent[:strings.LastIndex(scriptContent, "}")], "}") + 1

	var parsed struct {
		Movies []struct {
			Name  string `json:"name"`
			Year  int    `json:"year"`
			URL   string `json:"url"`
			Image string `json:"image"`
			Class string `json:"meterClass"`
			Score int    `json:"meterScore"`
		} `json:"movies"`
	}

	if err := json.Unmarshal([]byte(scriptContent[startIndex:endIndex]), &parsed); err != nil {
		return nil, ""
	}

	if len(parsed.Movies) <= 0 {
		resp.Description = fmt.Sprintf("Je n'ai trouvé aucun résultat pour **%s** :(", request)
		return resp, ""
	}

	url := "https://www.rottentomatoes.com" + parsed.Movies[0].URL

	resp.Title = fmt.Sprintf("%s (%d)", parsed.Movies[0].Name, parsed.Movies[0].Year)
	resp.Image = &discordgo.MessageEmbedImage{URL: parsed.Movies[0].Image, Width: 120}

	resp.Fields = append(resp.Fields, embed.EmbedField(
		"Tomatometer",
		fmt.Sprintf("%d%% (%s)", parsed.Movies[0].Score, strings.Replace(parsed.Movies[0].Class, "_", " ", -1)),
		true,
	))
	resp.Fields = append(resp.Fields, embed.EmbedField("Plus d'informations", url, true))

	switch parsed.Movies[0].Class {
	case "certified_fresh":
		resp.Color = 0x90ee90
	case "fresh":
		resp.Color = 0xeead00
	case "rotten":
		resp.Color = 0xc26001
	}

	resp.Footer = &discordgo.MessageEmbedFooter{
		Text: "Source : http://rottentomatoes.com/",
	}

	doc, err = httpreq.HTTPGetAsHTML(url)
	if err != nil {
		return resp, ""
	}

	consensus, _ := doc.Root().Search(xpath.Compile("//p[contains(@class, 'critic_consensus')]"))
	resp.Description = regexp.MustCompile("Critics? Consensus:").ReplaceAllString(consensus[0].Content(), "")

	return resp, ""
}

func init() {
	commands["tomatoes"] = &bot.Command{
		Function: commandTomatoes, HelpText: "Cherche les notes d'un film sur Rotten Tomatoes",
		Arguments: []bot.CommandArgument{
			{Name: "film", Description: "Le titre du film à rechercher", ArgType: "string"},
		},
		RequiredArguments: []string{"film"},
	}
}
