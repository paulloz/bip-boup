package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func command8Ball(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {

	gifs := []string{"http://gph.is/2iZqHD9", "http://gph.is/2tdjYNH", "http://gph.is/2sI6Igk", "http://gph.is/2mUXGgg",
		"http://gph.is/2mRFnW0", "http://gph.is/2pGyYTO", "https://gph.is/2uX3TvQ", "http://gph.is/2jthsLC", "http://gph.is/29dbWuO"}

	answers := []string{"Essaie plus tard.", "Essaie encore.", "J'ai pas d'avis.", "C'est ton destin.", "Le sort en est jeté.",
		"Une chance sur deux.", "Repose ta question.", "D'après moi oui.", "C'est certain.", "Oui, absolument.",
		"Tu peux compter dessus.", "Sans aucun doute.", "C'est très probable.", "Oui !", "C'est bien parti pour :wink:", "C'est non, désole…",
		"Peu probable…", "Faut pas rêver…", "N'y compte même pas.", "Impossible !", "Je m'en fout complet, demande à quelqu'un d'autre et laisse moi pioncer bordel…"}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	msg, _ := Bot.DiscordSession.ChannelMessageSend(env.Channel.ID, gifs[rng.Intn(len(gifs))])

	time.Sleep(7 * time.Second)

	Bot.DiscordSession.ChannelMessageDelete(env.Channel.ID, msg.ID)

	return nil, fmt.Sprintf("%s : %s", env.Message.Author.Mention(), answers[rng.Intn(len(answers))])
}
