package main

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func spawnBot() int {
	Info.Println("Spawning a new bot process...")

	bot := exec.Command(os.Args[0], "-bot", "-masterpid", strconv.Itoa(os.Getpid()))

	bot.Stdout = os.Stdout
	bot.Stderr = os.Stderr

	err := bot.Start()
	if err != nil {
		panic(err)
	}

	return bot.Process.Pid
}

func isBotAlive(pid int) bool {
	botProcess, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	return botProcess.Signal(syscall.Signal(0)) == nil
}
