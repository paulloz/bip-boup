package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/paulloz/bip-boup/queue"
)

// Globals
var (
	Bot   *BotConfig
	Queue *queue.Q

	InstanceId string

	IsThisABot bool
	MasterPID  int
	ConfigFile string
	IsDebug    bool

	GitCommit string
)

func init() {
	flag.BoolVar(&IsThisABot, "bot", false, "launch bot without a master process")
	flag.IntVar(&MasterPID, "masterpid", -1, "this master process' PID")
	flag.StringVar(&ConfigFile, "config", "config.json", "path to the .json configuration file")
	flag.BoolVar(&IsDebug, "debug", false, "launch in debug mode")
	flag.StringVar(&InstanceId, "id", "", "an instance identifier (not actually used for anything)")
	flag.Parse()

	if IsThisABot {
		initLog("BOT", IsDebug)
		initConfig(ConfigFile)
		Queue = queue.NewQueue(Bot.Database)
		initCache()
	} else {
		initLog("MASTER", IsDebug)
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

			if IsDebug {
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
			func() {
				sc := make(chan os.Signal, 1)
				signal.Notify(sc, syscall.SIGINT)
				ticker := time.NewTicker(time.Second * 30).C
				for {
					select {
					case <-sc:
						return
					case <-ticker:
						Queue.GoThrough(Bot.DiscordSession.ChannelMessageSendEmbed)
					}
				}
			}()

			Info.Println("Disconnecting from Discord...")
			discord.Close()

			saveConfig(ConfigFile)
			clearCache()
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

	restartFile := "/tmp/bip-boup.restart"
	fileData, err := ioutil.ReadFile(restartFile)
	if err == nil {
		Bot.DiscordSession.ChannelMessageSend(string(fileData), "Je suis là !")
		os.Remove(restartFile)
	} else if !os.IsNotExist(err) {
		panic(err)
	}

	os.Remove(os.Args[0] + ".old")

	updateFile := "/tmp/bip-boup.update"
	fileData, err = ioutil.ReadFile(updateFile)
	if err == nil {
		lines := strings.Split(string(fileData), "\n")
		if len(lines) >= 2 {
			Bot.DiscordSession.ChannelMessageDelete(lines[0], lines[1])
			Bot.DiscordSession.ChannelMessageSendEmbed(lines[0], &discordgo.MessageEmbed{
				Title: "Mise à jour", Color: 0x90ee90,
				Description: fmt.Sprintf("Mise à jour vers ``%s`` terminée.", GitCommit),
			})

			if len(lines) >= 3 {
				os.RemoveAll(lines[2])
			}
		}
		os.Remove(updateFile)
	} else if !os.IsNotExist(err) {
		panic(err)
	}
}
