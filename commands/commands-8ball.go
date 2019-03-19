package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/paulloz/bip-boup/bot"
)

func command8Ball(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {

	gifs := []string{"http://gph.is/2iZqHD9", "http://gph.is/2tdjYNH", "http://gph.is/2sI6Igk", "http://gph.is/2mUXGgg",
		"http://gph.is/2mRFnW0", "http://gph.is/2pGyYTO", "https://gph.is/2uX3TvQ", "http://gph.is/2jthsLC", "http://gph.is/29dbWuO"}

	answers := []string{"Essaie plus tard.", "Essaie encore.", "J'ai pas d'avis.", "C'est ton destin.", "Le sort en est jeté.",
		"Une chance sur deux.", "Repose ta question.", "D'après moi oui.", "C'est certain.", "Oui, absolument.",
		"Tu peux compter dessus.", "Sans aucun doute.", "C'est très probable.", "Oui !", "C'est bien parti pour :wink:", "C'est non, désolé…",
		"Peu probable…", "Faut pas rêver…", "N'y compte même pas.", "Impossible !", "Je m'en fout complet, demande à quelqu'un d'autre et laisse moi pioncer bordel…"}

	msg, _ := env.Session.ChannelMessageSend(env.Channel.ID, gifs[rand.Intn(len(gifs))])

	time.Sleep(7 * time.Second)

	env.Session.ChannelMessageDelete(env.Channel.ID, msg.ID)

	return nil, fmt.Sprintf("%s : %s", env.Message.Author.Mention(), answers[rand.Intn(len(answers))])
}

func init() {
	commands["8ball"] = &bot.Command{
		Function: command8Ball, HelpText: "Aide à choisir dans les moments difficiles.",
		Arguments: []bot.CommandArgument{
			{Name: "une question", Description: "Une interrogation sur laquelle vous souhaitez une réponse", ArgType: "string"},
		},
		RequiredArguments: []string{"une question"},
	}
}
