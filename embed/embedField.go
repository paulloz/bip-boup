package embed

import (
	"github.com/bwmarrin/discordgo"
)

func EmbedField(name string, value string, inlineOpt ...bool) *discordgo.MessageEmbedField {
	inline := false
	if len(inlineOpt) > 0 {
		inline = inlineOpt[0]
	}

	return &discordgo.MessageEmbedField{Name: name, Value: value, Inline: inline}
}
