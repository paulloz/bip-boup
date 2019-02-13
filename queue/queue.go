package queue

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	queueFile = "/tmp/bip-boup.queue"
)

type queueItem struct {
	Message string `json:"m"`
	Channel string `json:"c"`
	Time    string `json:"t"`
}

func createQueueFile() {
	os.Create(queueFile)
	os.Chmod(queueFile, 0600)
}

func loadQueueFile() []queueItem {
	fileHandler, err := os.Open(queueFile)
	defer fileHandler.Close()
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		createQueueFile()
		return loadQueueFile()
	}

	var items []queueItem

	decoder := json.NewDecoder(fileHandler)
	err = decoder.Decode(&items)
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	return items
}

func writeQueueFile(items *[]queueItem) {
	fileHandler, err := os.OpenFile(queueFile, os.O_WRONLY|os.O_TRUNC, 0600)
	defer fileHandler.Close()
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		createQueueFile()
		writeQueueFile(items)
		return
	}

	encoder := json.NewEncoder(fileHandler)
	err = encoder.Encode(items)
	if err != nil {
		panic(err)
	}
}

func addItemInQueue(item queueItem) {
	items := loadQueueFile()
	items = append(items, item)
	writeQueueFile(&items)
}

func Queue(channel string, message string, t time.Time) {
	item := queueItem{
		Message: message,
		Channel: channel,
		Time:    t.UTC().Format(time.ANSIC),
	}
	addItemInQueue(item)
}

func GoThroughQueue(send func(string, *discordgo.MessageEmbed) (*discordgo.Message, error)) {
	items := loadQueueFile()
	var newItems []queueItem

	now := time.Now().UTC().Truncate(time.Minute)
	for _, item := range items {
		parsed, err := time.Parse(time.ANSIC, item.Time)
		if err != nil {
			continue
		}

		if now.Unix() == parsed.Truncate(time.Minute).Unix() {
			send(item.Channel, &discordgo.MessageEmbed{Title: "Rappel", Description: item.Message})
		} else if now.Unix() < parsed.Unix() {
			newItems = append(newItems, item)
		}
	}

	writeQueueFile(&newItems)
}
