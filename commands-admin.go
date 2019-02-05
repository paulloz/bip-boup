package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/src-d/go-git.v4"
)

func commandRestart(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	Info.Println("Received a restart command, exiting...")

	Bot.DiscordSession.ChannelMessageSend(env.Channel.ID, "Je reviens dans un instant...")
	fileHandler, err := os.Create("/tmp/bip-boup.restart")
	if err == nil {
		fileHandler.Write([]byte(env.Channel.ID))
		fileHandler.Close()
	}

	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		panic(err)
	}
	process.Signal(syscall.SIGINT) // Exit, the master process will start a new bot
	return nil, ""
}

func commandUpdate(args []string, env *CommandEnvironment) (*discordgo.MessageEmbed, string) {
	errorEmbed := &discordgo.MessageEmbed{Title: "Erreur de mise à jour", Color: 0xa73902}

	updateEmbed, err := Bot.DiscordSession.ChannelMessageSendEmbed(env.Channel.ID, &discordgo.MessageEmbed{
		Title: "Mise à jour", Color: 0x90ee90, Description: "Mise à jour en cours...",
	})

	output, err := exec.Command("go", "version").CombinedOutput()
	if len(output) <= 0 || err != nil {
		Bot.DiscordSession.ChannelMessageDelete(env.Channel.ID, updateEmbed.ID)
		errorEmbed.Description = "Impossible d'exécuter ``go version``."
		return errorEmbed, ""
	}

	tmpDir, err := ioutil.TempDir(Bot.CacheDir, "update")
	if err != nil {
		Bot.DiscordSession.ChannelMessageDelete(env.Channel.ID, updateEmbed.ID)
		errorEmbed.Description = "Impossible de créer un dossier temporaire (``" + err.Error() + "``)."
		return errorEmbed, ""
	}
	defer os.RemoveAll(tmpDir)

	repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{URL: Bot.RepoURL, Depth: 1})
	if err != nil {
		Bot.DiscordSession.ChannelMessageDelete(env.Channel.ID, updateEmbed.ID)
		errorEmbed.Description = "Impossible de cloner le dépôt (``" + err.Error() + "``)."
		return errorEmbed, ""
	}

	ref, err := repo.Head()
	hash := ref.Hash().String()[:7]
	if hash == GitCommit {
		Bot.DiscordSession.ChannelMessageDelete(env.Channel.ID, updateEmbed.ID)
		return &discordgo.MessageEmbed{Title: "Mise à jour", Color: 0x90ee90, Description: "Aucune mise à jour nécessaire."}, ""
	}

	outputFile := fmt.Sprintf("%s/%s", tmpDir, os.Args[0])
	buildCommand := exec.Command("go", "build", "-ldflags", "-X main.GitCommit="+hash, "-o", outputFile)
	buildCommand.Dir = tmpDir
	output, err = buildCommand.CombinedOutput()
	if err != nil {
		Bot.DiscordSession.ChannelMessageDelete(env.Channel.ID, updateEmbed.ID)
		errorEmbed.Description = fmt.Sprintf("Impossible de compiler ``%s``.\n%s", hash, err.Error())
		return errorEmbed, ""
	}
	os.Rename(os.Args[0], os.Args[0]+".old")
	os.Rename(outputFile, os.Args[0])

	processToKillPID := os.Getpid()
	if MasterPID > 0 {
		processToKillPID = MasterPID
	}
	newProcess := exec.Command(os.Args[0], "-killoldpid", fmt.Sprintf("%d", processToKillPID))
	if MasterPID < 0 {
		newProcess.Args = append(newProcess.Args, "-bot")
	}
	newProcess.Stdout = os.Stdout
	newProcess.Stderr = os.Stderr
	err = newProcess.Start()
	if err != nil {
		panic(err)
	}
	newProcess.Process.Release()

	fileHandler, err := os.Create("/tmp/bip-boup.update")
	if err == nil {
		fileHandler.Write([]byte(env.Channel.ID + "\n" + updateEmbed.ID))
		fileHandler.Close()
	}

	return nil, ""
}
