package commands

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"

	"github.com/paulloz/bip-boup/bot"
	"github.com/paulloz/bip-boup/when/fr"
)

func commandReminder(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	w := when.New(nil)
	w.Add(common.All...)
	w.Add(fr.All...)

	text := strings.Join(args, " ")
	r, err := w.Parse(text, time.Now())
	if err != nil || r == nil {
		return nil, ""
	}

	split := strings.Split(text, r.Text)
	longest := [2]int{-1, -1}
	for i := range split {
		split[i] = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(split[i]), "que"), "à"))
		if len(split[i]) > longest[1] {
			longest[0] = i
			longest[1] = len(split[i])
		}
	}

	if longest[1] <= 0 {
		return nil, ""
	}

	b.Queue.Queue(env.Channel.ID, split[longest[0]], r.Time)

	return nil, "ok :ok_hand:"
}

func init() {
	commands["rappel"] = &bot.Command{
		Function: commandReminder,
		HelpText: "Enverra un message de rappel à l'heure indiquée",
		Arguments: []bot.CommandArgument{
			{Name: "requête", Description: "Une requête en langage naturel qui sera plus ou moins bien interprétée", ArgType: "string"},
		},
		RequiredArguments: []string{"requête"},
	}
}
