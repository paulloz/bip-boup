package main

import (
	"encoding/json"
	"os"

	"github.com/bwmarrin/discordgo"
)

type BotConfig struct {
	Commands       map[string]*Command `json:"-"`
	DiscordSession *discordgo.Session  `json:"-"`

	BotName   string `json:"-"`
	AuthToken string `json:"AuthToken"`

	CommandPrefix string `json:"CommandPrefix"`

	Admins []string `json:"Admins"`
	Debug  bool     `json:"-"`
}

func initConfig(file string) {
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
}
