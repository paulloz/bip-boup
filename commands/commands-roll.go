package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/paulloz/bip-boup/bot"
	"github.com/paulloz/bip-boup/embed"
)

func commandRoll(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	var dices []int

	for _, a := range args {
		tmp := strings.Split(a, "d")
		var err error
		var number int
		var max int
		if len(tmp) == 1 {
			number = 1
			max, err = strconv.Atoi(tmp[0])
		} else {
			if tmp[0] == "" {
				number = 1
			} else {
				number, err = strconv.Atoi(tmp[0])
			}
			max, err = strconv.Atoi(tmp[1])
		}
		if err != nil {
			return nil, "La requête n'as pas le bon format."
		}
		dices = append(dices, createDices(number, max)...)
	}

	if len(dices) > 50 {
		return nil, "Wesh, calme-toi " + env.Message.Author.Mention() + " et demande moins de 50 dés."
	}

	path := b.CacheDir + "rng"

	values, err := getDices(path, len(dices), b.BotConfig.RandomDotOrgToken, dices)
	if err != nil {
		return nil, fmt.Sprintf("L'erreur suivante est apparue: `%s`", err.Error())
	}
	return formatDices(dices, values, env.Message.Author.Mention())
}

func formatDices(dices []int, values []int, user string) (*discordgo.MessageEmbed, string) {

	fields := []*discordgo.MessageEmbedField{}
	res := make(map[int][]int)
	for i := 0; i < len(values); i++ {
		res[dices[i]] = append(res[dices[i]], values[i])
	}

	sum := 0
	for key, value := range res {
		fields = append(fields, embed.EmbedField(strconv.Itoa(key), fmt.Sprint(value), true))
		for _, i := range value {
			sum += i
		}
	}
	fields = append(fields, embed.EmbedField("Somme", fmt.Sprint(sum), false))

	return &discordgo.MessageEmbed{
		Title:       "Nouveau tirage",
		Description: fmt.Sprint("C'est ", user, " qui me l'a demandé."),
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Randomisation fourni par random.org",
		},
	}, ""
}

func createDices(number int, max int) []int {
	var res []int
	for i := 0; i < number; i++ {
		res = append(res, max)
	}
	return res
}

func getDices(path string, howMany int, token string, dices []int) ([]int, error) {
	file, _ := os.Open(path)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) == 0 || lines[0] == "" {
		callRandomOrg(token, path)
	}

	var numbers []int
	if len(lines) < howMany {
		d, err := getDices(path, len(lines), token, dices)
		if err != nil {
			return nil, errors.New("Erreur dans la création de dés")
		}

		numbers = append(numbers, d...)
		if len(callRandomOrg(token, path)) > 0 {
			return nil, errors.New("Erreur dans l'appel à Random.org")
		}

		d, err = getDices(path, howMany-len(lines), token, dices[len(lines):])
		if err != nil {
			return nil, err
		}

		numbers = append(numbers, d...)
	} else {
		for i := 0; i < howMany; i++ {
			n, err := strconv.ParseFloat(lines[i], 64)
			if err != nil {
				return nil, errors.New("Erreur dans la transformation des résultats")
			}
			d := n * float64(dices[i])
			numbers = append(numbers, int(math.Ceil(d)))
		}
		file, _ := os.Create(path)

		for _, nb := range lines[howMany:] {
			fmt.Fprintln(file, nb)
		}
		file.Close()
	}
	return numbers, nil
}

func callRandomOrg(token string, db string) string {
	url := "https://api.random.org/json-rpc/1/invoke"
	params := Params{
		APIKey:        token,
		N:             1000,
		DecimalPlaces: 5,
	}

	request := Request{
		Jsonrpc: "2.0",
		Method:  "generateDecimalFractions",
		Params:  params,
		ID:      "osef",
	}

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(request)

	req, err := http.NewRequest("POST", url, buffer)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	randomOrg, _ := ioutil.ReadAll(resp.Body)

	var data RandomOrg
	json.Unmarshal(randomOrg, &data)
	if data.Result.BitsUsed == 0 {
		return "Error: " + data.Failed.Message
	}

	writeNumbers(data.Result.Random.Data, db)

	return ""
}

func writeNumbers(numbers []float64, path string) {
	file, _ := os.Create(path)

	for _, nb := range numbers {
		fmt.Fprintln(file, nb)
	}
	file.Close()
}

type Request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      string `json:"id"`
}

type Params struct {
	APIKey        string `json:"apiKey"`
	N             int    `json:"n"`
	DecimalPlaces int    `json:"decimalPlaces"`
}

type RandomOrg struct {
	Jsonrpc string `json:"jsonrpc"`
	Failed  Failed `json:"error"`
	Result  Result `json:"result"`
}

type Result struct {
	Random        Random `json:"random"`
	BitsUsed      int    `json:"bitsused"`
	BitsLeft      int    `json:"bitsleft"`
	RequestsLeft  int    `json:"requestsleft"`
	AdvisoryDelay int    `json:"advisorydelay"`
}

type Random struct {
	Data           []float64 `json:"data"`
	CompletionTime string    `json:"completiontime"`
}

type Failed struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

func init() {
	commands["roll"] = &bot.Command{
		Function: commandRoll,
		HelpText: "Lance autant de dés que vous le souhaitez",
		Arguments: []bot.CommandArgument{
			{Name: "requête", Description: "Une requête sous la forme xdy, x étant le nombre de dés (vide accepté) et y le type de dés", ArgType: "string"},
		},
		RequiredArguments: []string{"requête"},
	}
}
