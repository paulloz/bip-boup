package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

// BotConfig ...
type BotConfig struct {
	Commands       map[string]*Command `json:"-"`
	DiscordSession *discordgo.Session  `json:"-"`

	BotName   string `json:"-"`
	AuthToken string `json:"AuthToken"`

	CommandPrefix string `json:"CommandPrefix"`

	Admins []string `json:"Admins"`

	CacheDir string `json:"-"`
	Database string `json:"database"`

	Modified bool `json:"-"`

	RepoURL string `json:"-"`
}

func initConfig(file string) {
	Bot = &BotConfig{}

	fileHandler, err := os.Open(file)
	defer fileHandler.Close()
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(fileHandler)
	err = decoder.Decode(&Bot)
	if err != nil {
		panic(err)
	}

	checkConfig()
}

func checkConfig() {
	if len(Bot.CommandPrefix) <= 0 {
		Bot.CommandPrefix = "!"
	}

	if len(Bot.Database) <= 0 {
		db := "/tmp/%s-queue.json"
		dbSuffix := "bip-boup"
		if len(InstanceId) > 0 {
			dbSuffix = InstanceId
		}
		Bot.Database = fmt.Sprintf(db, dbSuffix)
	}

	Bot.RepoURL = "https://github.com/paulloz/bip-boup.git"

	Bot.Modified = false
}

func saveConfig(file string) {
	if !Bot.Modified {
		return
	}

	fileHandler, err := os.OpenFile(file, os.O_WRONLY, 0644)
	defer fileHandler.Close()
	if err != nil {
		return
	}

	encoder := json.NewEncoder(fileHandler)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(&Bot)
	if err != nil {
		Error.Println(err.Error())
	}
}
