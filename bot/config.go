package bot

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/paulloz/bip-boup/log"
)

func (b *BotConfig) checkConfig(instanceID string) {
	if len(b.CommandPrefix) <= 0 {
		b.CommandPrefix = "!"
	}

	if len(b.Database) <= 0 {
		db := "/tmp/%s-queue.json"
		dbSuffix := "bip-boup"
		if len(instanceID) > 0 {
			dbSuffix = instanceID
		}
		b.Database = fmt.Sprintf(db, dbSuffix)
	}

	b.RepoURL = "https://github.com/paulloz/bip-boup.git"

	b.Modified = false
}

func (b *BotConfig) SaveConfig(file string) {
	if !b.Modified {
		return
	}

	fileHandler, err := os.OpenFile(file, os.O_WRONLY, 0644)
	defer fileHandler.Close()
	if err != nil {
		return
	}

	encoder := json.NewEncoder(fileHandler)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(&b)
	if err != nil {
		log.Error.Println(err.Error())
	}
}
