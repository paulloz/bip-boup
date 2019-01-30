package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Bot *BotConfig

	IsThisABot bool
	MasterPID  int
)

func init() {
	var configFile string

	flag.BoolVar(&IsThisABot, "bot", false, "launch bot without a master process")
	flag.IntVar(&MasterPID, "masterpid", -1, "this master process' PID")
	flag.StringVar(&configFile, "config", "config.json", "path to the .json configuration file")

	cliCommand := flag.String("command", "", "perform a command instead of running the bot")

	Bot = &BotConfig{CommandPrefix: "!", Debug: false}

	flag.Parse()

	if IsThisABot {
		initLog("BOT")

		initConfig(configFile)

		if len(*cliCommand) > 0 {
			initCommands()
			callCommand(*cliCommand, flag.Args(), &CommandEnvironment{})
			os.Exit(0)
		}
	} else {
		initLog("MASTER")
	}
}

func main() {
	Info.Println("Bip-boup © pauljoannon: 2018-2019")
	Info.Println("Current PID is " + strconv.Itoa(os.Getpid()))

	if IsThisABot {
		if len(Bot.AuthToken) > 0 {
			Info.Println("Creating a Discord session...")

			discord, err := discordgo.New("Bot " + Bot.AuthToken)
			if err != nil {
				panic(err)
			}

			if Bot.Debug {
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
			Error.Println("Couldn't find an authentification token, exiting...")
			os.Exit(1)
		}
	} else {
		botPID := spawnBot()

		Info.Println("Waiting for SIGINT syscall signal to terminate...")

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT)
		watchdog := time.Tick(10 * time.Second)

		for {
			select {
			case _, ok := <-sc:
				if ok {
					Info.Println("Waiting for the bot process to exit...")

					botProcess, _ := os.FindProcess(botPID)
					botProcess.Signal(syscall.SIGINT)
					botProcess.Wait()
					os.Exit(0)
				}
			case <-watchdog:
				if !isBotAlive(botPID) {
					botPID = spawnBot()
				}
			}
		}
	}
}

func discordReady(session *discordgo.Session, event *discordgo.Ready) {
	Bot.DiscordSession = session
	Bot.BotName = session.State.User.Username

	Info.Println("Registering commands...")
	initCommands()

	Info.Println("Everything is ready!")
}
