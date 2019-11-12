package commands

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/paulloz/bip-boup/bot"
	"github.com/paulloz/bip-boup/embed"
)

func commandTeam(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {

	names := args[1:]
	numberOfTeams, _ := strconv.ParseFloat(args[0], 64)

	teams := make([]string, int(numberOfTeams))

	cpt := 0.0
	for len(names) > 0 {
		rand.Seed(time.Now().UnixNano())
		player := rand.Intn(len(names))
		teams[int(math.Mod(cpt, numberOfTeams))] += " " + names[player]
		cpt++
		names = remove(names, player)
	}

	fields := []*discordgo.MessageEmbedField{}

	for key, value := range teams {
		fields = append(fields, embed.EmbedField(fmt.Sprint("Équipe ", key+1), fmt.Sprint(value), false))
	}

	return &discordgo.MessageEmbed{
		Title:       "Voici vos équipes",
		Description: fmt.Sprint("C'est ", env.Message.Author.Mention(), " qui me l'a demandé."),
		Fields:      fields,
	}, ""
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func init() {
	commands["team"] = &bot.Command{
		Function: commandTeam,
		HelpText: "Crée vos équipe",
		Arguments: []bot.CommandArgument{
			{Name: "requête", Description: "Une requête sous la forme `X Agathe Benoît Cedric`, X étant le nombre d'équipe et la liste des prénoms à répartir", ArgType: "string"},
		},
		RequiredArguments: []string{"requête"},
	}
}
