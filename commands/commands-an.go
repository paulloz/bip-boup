package commands

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/moovweb/gokogiri"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/paulloz/bip-boup/bot"
	"github.com/paulloz/bip-boup/embed"
	"github.com/paulloz/bip-boup/httpreq"
	ss "github.com/paulloz/bip-boup/strings"
)

// Depute ...
type Depute struct {
	D struct {
		BirthDate        string `json:"date_naissance"`
		Name             string `json:"nom"`
		Circo            string `json:"nom_circo"`
		Profession       string `json:"profession"`
		Sex              string `json:"sexe"`
		Slug             string `json:"slug"`
		Twitter          string `json:"twitter"`
		URLAN            string `json:"url_an"`
		URLNosDeputes    string `json:"url_nosdeputes"`
		Responsabilities []struct {
			R struct {
				Organism string `json:"organisme"`
				Function string `json:"function"`
			} `json:"responsabilite"`
		} `json:"responsabilites"`
		Group struct {
			Organism string `json:"organisme"`
			Function string `json:"fonction"`
		} `json:"groupe"`
	} `json:"depute"`
}

func commandDirectANNoSession(args []string, env *bot.CommandEnvironment) (*discordgo.MessageEmbed, string) {
	body, err := httpreq.HTTPGet("http://data.assemblee-nationale.fr/static/openData/repository/15/vp/seances/seances_publique_libre_office.csv")
	if err != nil {
		return nil, ""
	}

	now := time.Now().Add(-(time.Hour / 2))
	d := now.Format("2006-01-02")
	t := now.Format("15:04")

	sessions := []*struct {
		Key   string
		Value []string
	}{}

	reader := csv.NewReader(bytes.NewReader(body))
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if record[0] == d {
			if record[1] >= t {
				value := regexp.MustCompile("\\s{2,}").Split(record[2], -1) // They use 4, 5 or 6 spaces to split data ¯\_(ツ)_/¯
				value = ss.Every(value, func(s string) string { return fmt.Sprintf("  - %s.", s) })
				sessions = append(sessions, &struct {
					Key   string
					Value []string
				}{Key: record[1], Value: value})
			}
		}
	}

	sort.SliceStable(sessions, func(i, j int) bool {
		return sessions[i].Key < sessions[j].Key
	})

	if len(sessions) > 0 {
		fields := []*discordgo.MessageEmbedField{}
		for _, session := range sessions {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "à " + session.Key,
				Value:  strings.Join(session.Value[:3], "\n"),
				Inline: false,
			})
		}

		return &discordgo.MessageEmbed{
			Title:  "Pas de séance en cours, prochaines séances",
			Fields: fields,
		}, ""
	}

	return nil, ""
}

func commandDirectAN(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	url := "http://videos.assemblee-nationale.fr/direct.1"

	doc, err := httpreq.HTTPGetAsHTML(url)
	if err != nil {
		return nil, ""
	}
	defer doc.Free()

	rootNode := doc.Root()
	playerTitleNode, _ := rootNode.Search("//div[contains(@class, 'playerTitle')]")
	if len(playerTitleNode) > 0 {
		descriptionNode, _ := rootNode.Search("//div[contains(@class, 'txtEditorial')]")
		subjects := ss.Choose(strings.Split(descriptionNode[0].Content(), "- "), func(s string) bool { return len(s) > 0 })
		subjects = ss.Every(subjects, func(s string) string { return fmt.Sprintf(" - %s.", s) })

		return &discordgo.MessageEmbed{
			Title: "Séance en cours",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{Name: "Ordre du jour", Value: strings.Join(subjects, "\n"), Inline: false},
				&discordgo.MessageEmbedField{Name: "Direct", Value: url, Inline: false},
			},
		}, ""
	}

	return commandDirectANNoSession(args, env)
}

func buildDeputeCache(b *bot.Bot, name string, opt ...string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://www.nosdeputes.fr/deputes/enmandat/xml", nil)
	if err != nil {
		return
	}

	if len(opt) > 0 {
		req.Header.Add("If-Modified-Since", opt[0])
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(body) == 0 {
		return
	}

	doc, err := gokogiri.ParseXml(body)
	if err != nil {
		return
	}
	defer doc.Free()

	values := map[string]string{}

	nodes, _ := doc.Root().Search("//depute")
	for _, node := range nodes {
		slug, _ := node.Search("slug")
		name, _ := node.Search("nom")

		values[slug[0].Content()] = strings.ToLower(name[0].Content())
	}

	b.SetCache(name, resp.Header["Date"][0], &values)
}

func commandDeputeSearch(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	cacheName := "deputes"
	var cache *bot.Cache

	if cache = b.GetCache(cacheName); cache == nil {
		buildDeputeCache(b, cacheName)
	} else {
		// Last-Modified is not working right now
		// buildDeputeCache(cacheName, cache.LastModified)
	}
	if cache = b.GetCache(cacheName); cache == nil {
		return nil, ""
	}

	search := strings.Join(args, " ")
	results := []string{}

	for slug, name := range *cache.Values {
		if strings.Contains(name, search) {
			results = append(results, slug)
		}
	}

	if len(results) == 1 {
		return commandDepute([]string{results[0]}, env, b)
	} else if len(results) > 1 {
		for _, slug := range results {
			env.Message.Content = fmt.Sprintf("%sdepute %s", b.CommandPrefix, slug)
			go commandDepute([]string{slug}, env, b)
		}
	}
	return nil, ""
}

func commandDepute(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	slug := strings.Join(ss.Every(args, func(s string) string {
		_s, _, _ := transform.String(t, s)
		return strings.ToLower(_s)
	}), "-")
	url := fmt.Sprintf("https://www.nosdeputes.fr/%s/json", slug)

	body, err := httpreq.HTTPGet(url)
	if err != nil {
		return nil, ""
	}
	if len(body) == 0 {
		return commandDeputeSearch(args, env, b)
	}

	var depute Depute
	if err = json.Unmarshal(body, &depute); err != nil {
		return nil, ""
	}

	fields := []*discordgo.MessageEmbedField{}

	genderE := ""
	if depute.D.Sex == "F" {
		genderE = "e"
	}

	if len(depute.D.Responsabilities) > 0 {
		respParl := []string{}
		if depute.D.Group.Function != "membre" {
			respParl = append(respParl, fmt.Sprintf("%s, %s.", depute.D.Group.Organism, depute.D.Group.Function))
		}
		for _, responsability := range depute.D.Responsabilities {
			respParl = append(respParl, fmt.Sprintf("%s, %s.", responsability.R.Organism, responsability.R.Function))
		}
		fields = append(fields, embed.EmbedField("Responsabilités parlementaires", strings.Join(respParl, "\n")))
	}

	parsedBirthDate, _ := time.Parse("2006-01-02", depute.D.BirthDate)
	age := int(time.Since(parsedBirthDate).Hours() / 24 / 365)
	fields = append(fields, embed.EmbedField("Âge", fmt.Sprintf("%d ans", age), true))
	fields = append(fields, embed.EmbedField("Profession", depute.D.Profession, true))

	if len(depute.D.Twitter) > 0 {
		fields = append(fields, embed.EmbedField("Twitter", fmt.Sprintf("[@%s](https://twitter.com/%s)", depute.D.Twitter, depute.D.Twitter), true))
	}

	fields = append(fields, embed.EmbedField("Fiches", fmt.Sprintf("[Assemblée Nationale](%s)\n[NosDéputés.fr](%s)", depute.D.URLAN, depute.D.URLNosDeputes), true))

	imageURL := fmt.Sprintf("https://www.nosdeputes.fr/depute/photo/%s/120", depute.D.Slug)

	return &discordgo.MessageEmbed{
		Title:       depute.D.Name,
		Description: fmt.Sprintf("Député%s %s (%s).", genderE, depute.D.Group.Organism, depute.D.Circo),
		Image:       &discordgo.MessageEmbedImage{URL: imageURL},
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Source : NosDéputés.fr par Regards Citoyens à partir de l'Assemblée nationale et du Journal Officiel",
		},
	}, ""
}

func init() {
	commands["directan"] = &bot.Command{Function: commandDirectAN, HelpText: "Envoie un lien vers la séance publique en cours à l'Assemblée Nationale."}

	commands["député"] = &bot.Command{
		Function: commandDepute, HelpText: "Montre les informations à propos d'un député disponibles sur NosDéputés.fr.",
		Arguments: []bot.CommandArgument{
			{Name: "prénom", Description: "Le prénom du député", ArgType: "string"},
			{Name: "nom", Description: "Le nom du député", ArgType: "string"},
		},
		RequiredArguments: []string{"nom"},
	}
	commands["depute"] = &bot.Command{IsAliasTo: "député"}
}
