package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
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
