package bot

import (
	"encoding/json"
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/paulloz/bip-boup/queue"
)

type Command struct {
	Function          func([]string, *CommandEnvironment, *Bot) (*discordgo.MessageEmbed, string)
	HelpText          string
	Arguments         []CommandArgument
	RequiredArguments []string

	RequiredPermissions int
	IsAdmin             bool

	IsAliasTo string
}

type CommandArgument struct {
	Name        string
	ArgType     string
	Description string
}

type CommandEnvironment struct {
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	User    *discordgo.User
	Member  *discordgo.Member

	Message *discordgo.Message

	Session *discordgo.Session
}

type BotConfig struct {
	BotName   string `json:"-"`
	AuthToken string `json:"AuthToken"`

	CommandPrefix string `json:"CommandPrefix"`

	Admins []string `json:"Admins"`

	CacheDir string `json:"-"`
	Database string `json:"database"`
	MemeDir  string `json:"MemeDir"`

	Modified bool `json:"-"`

	RepoURL   string `json:"-"`
	GitCommit string `json:"-"`
}

type Bot struct {
	BotConfig

	DiscordSession *discordgo.Session `json:"-"`

	Commands map[string]*Command `json:"-"`

	Queue *queue.Q `json:"-"`
}

func (b *BotConfig) IsUserAdmin(user *discordgo.User) bool {
	for _, adminID := range b.Admins {
		if user.ID == adminID {
			return true
		}
	}
	return false
}

func NewBot(file string, instanceID string) *Bot {
	bot := &Bot{}

	fileHandler, err := os.Open(file)
	defer fileHandler.Close()
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(fileHandler)
	err = decoder.Decode(&bot)
	if err != nil {
		panic(err)
	}

	bot.checkConfig(instanceID)

	bot.initCache()

	bot.Queue = queue.NewQueue(bot.Database)

	return bot
}
