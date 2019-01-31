package main

import (
	"encoding/json"
	"fmt"
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

	CacheDir string `json:"-"`
	Modified bool   `json:"-"`
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

	Bot.CacheDir = "/tmp/bip-boup"
	os.Mkdir(Bot.CacheDir, os.ModeDir|0700)
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

type Cache struct {
	LastModified string
	Values       *map[string]string
}

func getCacheFileName(name string) string {
	return fmt.Sprintf("%s/%s.cached", Bot.CacheDir, name)
}

func getCache(name string) (cache *Cache) {
	fileHandler, err := os.Open(getCacheFileName(name))
	if err != nil {
		return
	}
	defer fileHandler.Close()

	decoder := json.NewDecoder(fileHandler)
	err = decoder.Decode(&cache)
	if err != nil {
		cache = nil
	}

	return
}

func setCache(name string, lastModified string, values *map[string]string) {
	fileHandler, err := os.OpenFile(getCacheFileName(name), os.O_CREATE|os.O_WRONLY, 0644)
	defer fileHandler.Close()
	if err != nil {
		Error.Println(err.Error())
		return
	}

	encoder := json.NewEncoder(fileHandler)
	err = encoder.Encode(&Cache{LastModified: lastModified, Values: values})
	if err != nil {
		Error.Println(err.Error())
	}
}

func clearCache() {
	os.RemoveAll(Bot.CacheDir)
}
