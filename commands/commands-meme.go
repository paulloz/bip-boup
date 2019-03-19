package commands

import (
	"bytes"
	"fmt"
	"image"
	"math/rand"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"

	"github.com/paulloz/bip-boup/bot"
)

var (
	fontSize  float64 = 36
	spongeBob image.Image
)

func commandSpongeBob(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	text := ""

	for _, c := range strings.Join(args, " ") {
		if rand.Float32() < 0.5 {
			text += strings.ToLower(string(c))
		} else {
			text += strings.ToUpper(string(c))
		}
	}

	dir := strings.TrimRight(b.MemeDir, "/")

	if spongeBob == nil {
		file, err := os.Open(fmt.Sprintf("%s/spongebob.jpeg", dir))
		if err != nil {
			return nil, ""
		}
		spongeBob, _, err = image.Decode(file)
		if err != nil {
			return nil, ""
		}
	}

	ctx := gg.NewContext(spongeBob.Bounds().Dx(), spongeBob.Bounds().Dy())
	ctx.DrawImage(spongeBob, 0, 0)

	ctx.LoadFontFace(fmt.Sprintf("%s/Impact.ttf", dir), fontSize)
	ctx.SetHexColor("#fff")
	lines := ctx.WordWrap(text, float64(spongeBob.Bounds().Dx()))
	for i, line := range lines {
		ctx.DrawStringAnchored(line, float64((spongeBob.Bounds().Dx() / 2)), float64(((fontSize * 2) + (float64(i) * fontSize))), 0.5, 0.0)
	}

	var buffer bytes.Buffer
	ctx.EncodePNG(&buffer)

	env.Session.ChannelFileSend(env.Channel.ID, "spongebob.png", &buffer)

	return nil, ""
}

func init() {
	commands["spongebob"] = &bot.Command{
		Function: commandSpongeBob, HelpText: "Crée une image de mocking SpongeBob",
		Arguments: []bot.CommandArgument{
			{Name: "texte", Description: "Le texte à afficher sur l'image", ArgType: "string"},
		},
		RequiredArguments: []string{"text"},
	}
}
