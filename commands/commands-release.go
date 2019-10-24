package commands

// Rap-francais-discographie-2019-lyrics
// 3837973

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moovweb/gokogiri/xpath"

	"github.com/paulloz/bip-boup/bot"
	"github.com/paulloz/bip-boup/httpreq"
)

func commandRelease(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	urlFr := "https://genius.com/Rap-francais-discographie-2019-lyrics"
	title := "Sorties récentes et prochaines"

	resp := &discordgo.MessageEmbed{
		Title:       title,
		Description: "Aucun résultat.",
		Footer:      &discordgo.MessageEmbedFooter{Text: "Source: https://genius.com/"},
	}

	doc, err := httpreq.HTTPGetAsHTML(urlFr)
	if err != nil {
		return resp, ""
	}
	defer doc.Free()

	items := xpath.Compile("//div[@class='lyrics']")
	results, _ := doc.Root().Search(items)

	var releases []string
	for _, r := range results {
		releases = strings.Split(r.Content(), "\n")
	}

	res := manageTitle(releases)

	resp.Description = strings.Join(res, "\n")

	return resp, ""
}

func manageTitle(releases []string) []string {
	var titles []string
	for _, release := range releases {
		if release != "" {
			release = strings.TrimSpace(release)
			re, _ := regexp.MatchString(`[*|-] \d\d/\d\d .*`, release)
			if re {
				titles = append(titles, release)
			}
		}
	}

	var res []string
	for _, title := range titles {
		min := time.Now().Add(-24 * 7 * time.Hour)
		max := time.Now().Add(24 * 6 * time.Hour)
		day, _ := strconv.Atoi(strings.Split(strings.Fields(title)[1], "/")[0])
		month, _ := strconv.Atoi(strings.Split(strings.Fields(title)[1], "/")[1])
		check := time.Date(2019, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		if check.After(min) && check.Before(max) {
			res = append(res, title)
		}
	}
	return res
}

func init() {
	commands["release"] = &bot.Command{
		Function: commandRelease, HelpText: "Quels sont les prochains album à sortir.",
	}
}
