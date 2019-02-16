package queue

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type qItem struct {
	Message string `json:"m"`
	Channel string `json:"c"`
	Time    string `json:"t"`
}

type Q struct {
	fileName string
}

func (q *Q) createQueueFile() {
	os.Create(q.fileName)
	os.Chmod(q.fileName, 0600)
}

func (q *Q) loadQueueFile() []qItem {
	fileHandler, err := os.Open(q.fileName)
	defer fileHandler.Close()
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		q.createQueueFile()
		return q.loadQueueFile()
	}

	var items []qItem

	decoder := json.NewDecoder(fileHandler)
	err = decoder.Decode(&items)
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	return items
}

func (q *Q) writeQueueFile(items *[]qItem) {
	fileHandler, err := os.OpenFile(q.fileName, os.O_WRONLY|os.O_TRUNC, 0600)
	defer fileHandler.Close()
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		q.createQueueFile()
		q.writeQueueFile(items)
		return
	}

	encoder := json.NewEncoder(fileHandler)
	err = encoder.Encode(items)
	if err != nil {
		panic(err)
	}
}

func (q *Q) addItemInQueue(item qItem) {
	items := q.loadQueueFile()
	items = append(items, item)
	q.writeQueueFile(&items)
}

func (q *Q) Queue(channel string, message string, t time.Time) {
	item := qItem{
		Message: message,
		Channel: channel,
		Time:    t.UTC().Format(time.ANSIC),
	}
	q.addItemInQueue(item)
}

func (q *Q) GoThrough(send func(string, *discordgo.MessageEmbed) (*discordgo.Message, error)) {
	items := q.loadQueueFile()
	var newItems []qItem

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

	q.writeQueueFile(&newItems)
}

func (q *Q) GetLength() int {
	return len(q.loadQueueFile())
}

func NewQueue(fileName string) *Q {
	return &Q{fileName: fileName}
}
