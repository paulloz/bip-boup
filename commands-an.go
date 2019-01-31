package main

import (
	"bytes"
	"encoding/csv"
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
)

func commandDirectANNoSession(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	url := "http://data.assemblee-nationale.fr/static/openData/repository/15/vp/seances/seances_publique_libre_office.csv"
	resp, err := http.Get(url)
	if err != nil {
		return nil, ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
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
				value = every(value, func(s string) string { return fmt.Sprintf("  - %s.", s) })
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

func commandDirectAN(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	url := "http://videos.assemblee-nationale.fr/direct.1"
	resp, err := http.Get(url)
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
	playerTitleNode, _ := rootNode.Search("//div[contains(@class, 'playerTitle')]")
	if len(playerTitleNode) > 0 {
		descriptionNode, _ := rootNode.Search("//div[contains(@class, 'txtEditorial')]")
		subjects := choose(strings.Split(descriptionNode[0].Content(), "- "), func(s string) bool { return len(s) > 0 })
		subjects = every(subjects, func(s string) string { return fmt.Sprintf(" - %s.", s) })

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

func commandDepute(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	slug := strings.Join(every(args, func(s string) string {
		_s, _, _ := transform.String(t, s)
		return strings.ToLower(_s)
	}), "-")
	url := fmt.Sprintf("https://www.nosdeputes.fr/%s/xml", slug)

	resp, err := http.Get(url)
	if err != nil {
		return nil, ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(body) == 0 {
		return nil, ""
	}

	doc, err := gokogiri.ParseXml(body)
	if err != nil {
		return nil, ""
	}
	defer doc.Free()

	rootNode := doc.Root()

	fields := []*discordgo.MessageEmbedField{}

	name, _ := rootNode.Search("//depute/nom")
	group, _ := rootNode.Search("//groupe/organisme")
	groupFunction, _ := rootNode.Search("//groupe/fonction")
	circo, _ := rootNode.Search("//depute/nom_circo")
	profession, _ := rootNode.Search("//profession")
	birthDate, _ := rootNode.Search("//date_naissance")

	sex, _ := rootNode.Search("//depute/sexe")
	genderE := ""
	if sex[0].Content() == "F" {
		genderE = "e"
	}

	respParlNodes, _ := rootNode.Search("//responsabilites/responsabilite")
	if len(respParlNodes) > 0 {
		respParl := []string{}
		if groupFunction[0].Content() != "membre" {
			respParl = append(respParl, fmt.Sprintf("%s, %s.", group[0].Content(), groupFunction[0].Content()))
		}
		for _, respParlNode := range respParlNodes {
			org, _ := respParlNode.Search("organisme")
			function, _ := respParlNode.Search("fonction")
			respParl = append(respParl, fmt.Sprintf("%s, %s.", org[0].Content(), function[0].Content()))
		}
		fields = append(fields, embedField("Responsabilités parlementaires", strings.Join(respParl, "\n")))
	}

	// groupParlNodes, _ := rootNode.Search("//groupes_parlementaires/responsabilite")
	// if len(groupParlNodes) > 0 {
	// 	groupParl := []string{}
	// 	for _, groupParlNode := range groupParlNodes {
	// 		org, _ := groupParlNode.Search("organisme")
	// 		function, _ := groupParlNode.Search("fonction")
	// 		groupParl = append(groupParl, fmt.Sprintf("%s, %s.", org[0].Content(), function[0].Content()))
	// 	}
	// 	fields = append(fields, embedField("Groupes parlementaires", strings.Join(groupParl, "\n")))
	// }

	parsedBirthDate, _ := time.Parse("2006-01-02", birthDate[0].Content())
	age := int(time.Since(parsedBirthDate).Hours() / 24 / 365)
	fields = append(fields, embedField("Âge", fmt.Sprintf("%d ans", age), true))
	fields = append(fields, embedField("Profession", profession[0].Content(), true))

	twitter, _ := rootNode.Search("//twitter")
	if len(twitter) > 0 {
		handle := twitter[0].Content()
		fields = append(fields, embedField("Twitter", fmt.Sprintf("[@%s](https://twitter.com/%s)", handle, handle), true))
	}

	ficheAN, _ := rootNode.Search("//url_an")
	fiche, _ := rootNode.Search("//url_nosdeputes")

	fields = append(fields, embedField("Fiches", fmt.Sprintf("[Assemblée Nationale](%s)\n[NosDéputés.fr](%s)", ficheAN[0].Content(), fiche[0].Content()), true))

	imageURL := fmt.Sprintf("https://www.nosdeputes.fr/depute/photo/%s/120", slug)

	return &discordgo.MessageEmbed{
		Title:       name[0].Content(),
		Description: fmt.Sprintf("Député%s %s (%s).", genderE, group[0].Content(), circo[0].Content()),
		Image:       &discordgo.MessageEmbedImage{URL: imageURL},
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Source : NosDéputés.fr par Regards Citoyens à partir de l'Assemblée nationale et du Journal Officiel",
		},
	}, ""
}
