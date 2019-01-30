package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gojp/kana"
	"github.com/ikawaha/kagome/tokenizer"
)

func commandFurigana(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	text := strings.Join(args, " ")
	if len(text) <= 0 {
		return nil, ""
	}

	t := tokenizer.New()
	tokens := t.Tokenize(text)
	tokenized := make([]string, len(tokens))
	for _, token := range tokens {
		switch token.Class {
		case tokenizer.UNKNOWN:
			tokenized = append(tokenized, token.Surface)
		case tokenizer.KNOWN:
			features := token.Features()
			katakanas := features[len(features)-2]
			if kana.IsHiragana(token.Surface) {
				tokenized = append(tokenized, token.Surface)
			} else {
				value := fmt.Sprintf("%s(%s)", token.Surface, kana.RomajiToHiragana(kana.KanaToRomaji(katakanas)))
				tokenized = append(tokenized, value)
			}
		}
	}

	return &discordgo.MessageEmbed{Title: strings.Join(tokenized, " "), Color: 0xffffff}, ""
}
