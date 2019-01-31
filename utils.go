package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
)

func choose(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}

	return
}

func every(ss []string, f func(string) string) (ret []string) {
	for _, s := range ss {
		ret = append(ret, f(s))
	}

	return
}

func contains(ss []string, s string) bool {
	for _, _s := range ss {
		if _s == s {
			return true
		}
	}

	return false
}

func capitalize(s string) string {
	if len(s) <= 0 {
		return s
	}

	return strings.ToUpper(string(s[0])) + string(s[1:])
}

func embedField(name string, value string, inline_opt ...bool) *discordgo.MessageEmbedField {
	inline := false
	if len(inline_opt) > 0 {
		inline = inline_opt[0]
	}

	return &discordgo.MessageEmbedField{Name: name, Value: value, Inline: inline}
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func httpGetAsXml(url string) (*xml.XmlDocument, error) {
	body, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	return gokogiri.ParseXml(body)
}

func httpGetAsHtml(url string) (*html.HtmlDocument, error) {
	body, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	return gokogiri.ParseHtml(body)
}
