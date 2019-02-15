package main

import (
	"os"
	"os/exec"
	"strconv"
)

func spawnBot() int {
	Info.Println("Spawning a new bot process...")

	bot := exec.Command(os.Args[0], "-bot", "-masterpid", strconv.Itoa(os.Getpid()))

	if len(InstanceId) > 0 {
		bot.Args = append(bot.Args, []string{"-id", InstanceId}...)
	}

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

	processState, err := botProcess.Wait()
	if err != nil {
		return false
	}

	return processState.Exited()
}
