package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Cache ...
type Cache struct {
	LastModified string
	Values       *map[string]string
}

func getCacheFileName(name string) string {
	return fmt.Sprintf("%s/%s", Bot.CacheDir, name)
}

func initCache() {
	tempDir, err := ioutil.TempDir("", "bipboupcache")
	if err != nil {
		panic(err)
	}

	Bot.CacheDir = tempDir
}

func clearCache(leaveDirOpt ...bool) {
	os.RemoveAll(Bot.CacheDir)
	if len(leaveDirOpt) > 0 && leaveDirOpt[0] {
		initCache()
	}
}

func getCache(name string) (cache *Cache) {
	fileHandler, err := os.Open(getCacheFileName(name))
	defer fileHandler.Close()
	if err != nil {
		return
	}

	decoder := json.NewDecoder(fileHandler)
	err = decoder.Decode(&cache)
	if err != nil {
		cache = nil
	}

	return
}

func setCache(name string, lastModified string, values *map[string]string) {
	// First we check the cacheDir still exists. If not, we create a new one.
	fileInfo, err := os.Stat(Bot.CacheDir)
	if err != nil || !fileInfo.IsDir() {
		clearCache(true)
	}

	fileHandler, err := os.OpenFile(getCacheFileName(name), os.O_CREATE|os.O_WRONLY, 0644)
	defer fileHandler.Close()
	if err != nil {
		return
	}

	encoder := json.NewEncoder(fileHandler)
	err = encoder.Encode(&Cache{LastModified: lastModified, Values: values})
}
