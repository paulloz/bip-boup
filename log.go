package main

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	logFlags = log.Ldate | log.Ltime | log.Lshortfile

	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func initLog() {
	Debug = log.New(ioutil.Discard, "DEBUG: ", logFlags)
	if BotData.Debug {
		Debug = log.New(os.Stdout, "DEBUG: ", logFlags)
	}
	Info = log.New(os.Stdout, "INFO: ", logFlags)
	Warning = log.New(os.Stdout, "WARNING: ", logFlags)
	Error = log.New(os.Stdout, "ERROR: ", logFlags)
}
