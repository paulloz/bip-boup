package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	BotData *Bot
)

type Bot struct {
	Commands       map[string]*Command
	DiscordSession *discordgo.Session

	BotName       string
	CommandPrefix string
	Debug         bool
}

func discordReady(session *discordgo.Session, event *discordgo.Ready) {
	BotData.DiscordSession = session
	BotData.BotName = session.State.User.Username

	Info.Println("Registering commands...")
	initCommands()

	Info.Println("Everything is ready!")
}

func init() {
	BotData = &Bot{CommandPrefix: "!", Debug: false}
	initLog()

	cliCommand := flag.String("cmd", "", "perform a command instead of running the bot")
	flag.Parse()

	if len(*cliCommand) > 0 {
		initCommands()
		callCommand(*cliCommand, flag.Args(), &CommandEnvironment{})
		os.Exit(0)
	}
}

func main() {
	Info.Println("Bip-boup © pauljoannon: 2018-2019")
	Info.Println("Current PID is " + strconv.Itoa(os.Getpid()))

	token := os.Getenv("DISCORD_TOKEN")
	if len(token) > 0 {
		Info.Println("Creating a Discord session...")

		discord, err := discordgo.New("Bot " + token)
		if err != nil {
			panic(err)
		}

		if BotData.Debug {
			discord.LogLevel = discordgo.LogInformational
		}

		Info.Println("Registering Discord event handlers...")
		discord.AddHandler(discordMessageCreate)
		discord.AddHandler(discordReady)

		Info.Println("Connecting to Discord...")
		err = discord.Open()
		if err != nil {
			panic(err)
		}

		Info.Println("Waiting for SIGINT syscall signal to terminate...")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT)
		<-sc

		Info.Println("Disconnecting from Discord...")
		discord.Close()
	} else {

	}
}
