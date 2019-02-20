package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/src-d/go-git.v4"

	"github.com/paulloz/bip-boup/bot"
	"github.com/paulloz/bip-boup/log"
)

func commandRestart(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	log.Info.Println("Received a restart command, exiting...")

	env.Session.ChannelMessageSend(env.Channel.ID, "Je reviens dans un instant...")
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

func commandUpdate(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	errorEmbed := &discordgo.MessageEmbed{Title: "Erreur de mise à jour", Color: 0xa73902}

	updateEmbed, err := env.Session.ChannelMessageSendEmbed(env.Channel.ID, &discordgo.MessageEmbed{
		Title: "Mise à jour", Color: 0x90ee90, Description: "Mise à jour en cours...",
	})

	output, err := exec.Command("go", "version").CombinedOutput()
	if len(output) <= 0 || err != nil {
		env.Session.ChannelMessageDelete(env.Channel.ID, updateEmbed.ID)
		errorEmbed.Description = "Impossible d'exécuter ``go version``."
		return errorEmbed, ""
	}

	tmpDir, err := ioutil.TempDir(b.CacheDir, "update")
	if err != nil {
		env.Session.ChannelMessageDelete(env.Channel.ID, updateEmbed.ID)
		errorEmbed.Description = "Impossible de créer un dossier temporaire (``" + err.Error() + "``)."
		return errorEmbed, ""
	}
	defer os.RemoveAll(tmpDir)

	repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{URL: b.RepoURL, Depth: 1})
	if err != nil {
		env.Session.ChannelMessageDelete(env.Channel.ID, updateEmbed.ID)
		errorEmbed.Description = "Impossible de cloner le dépôt (``" + err.Error() + "``)."
		return errorEmbed, ""
	}

	ref, err := repo.Head()
	hash := ref.Hash().String()[:7]
	if hash == b.GitCommit {
		env.Session.ChannelMessageDelete(env.Channel.ID, updateEmbed.ID)
		return &discordgo.MessageEmbed{Title: "Mise à jour", Color: 0x90ee90, Description: "Aucune mise à jour nécessaire."}, ""
	}

	for i := 0; i < 2; i++ {
		dirs := ""
		dots := "\\."
		for j := 0; j < i; j++ {
			dirs += "/*"
			dots += "\\."
		}
		goFiles, _ := filepath.Glob(tmpDir + dirs + "/*.go")
		for _, goFile := range goFiles {
			sedCommand := exec.Command("sed", "-i", "s/github\\.com\\/paulloz\\/bip-boup\\//"+dots+"\\//g", goFile)
			sedCommand.Dir = tmpDir
			output, err = sedCommand.CombinedOutput()
			if err != nil {
				env.Session.ChannelMessageDelete(env.Channel.ID, updateEmbed.ID)
				errorEmbed.Description = fmt.Sprintf("Échec de la commande sed.\n%s\n%s", err.Error(), output)
				return errorEmbed, ""
			}
		}
	}

	outputFile := fmt.Sprintf("%s/%s", tmpDir, os.Args[0])
	buildCommand := exec.Command("go", "build", "-ldflags", ("-X main.GitCommit=" + hash), "-o", outputFile)
	buildCommand.Dir = tmpDir
	output, err = buildCommand.CombinedOutput()
	if err != nil {
		env.Session.ChannelMessageDelete(env.Channel.ID, updateEmbed.ID)
		errorEmbed.Description = fmt.Sprintf("Impossible de compiler ``%s``.\n%s\n%s", hash, err.Error(), output)
		return errorEmbed, ""
	}
	os.Rename(os.Args[0], os.Args[0]+".old")
	os.Rename(outputFile, os.Args[0])

	fileHandler, err := os.Create("/tmp/bip-boup.update")
	if err == nil {
		fileHandler.Write([]byte(env.Channel.ID + "\n" + updateEmbed.ID + "\n" + b.CacheDir))
		fileHandler.Close()
	}

	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		panic(err)
	}
	process.Signal(syscall.SIGINT) // Exit, the master process will start a new bot

	return nil, ""
}

func commandQueue(args []string, env *bot.CommandEnvironment, b *bot.Bot) (*discordgo.MessageEmbed, string) {
	n := b.Queue.GetLength()
	s := func() string {
		if n == 1 {
			return ""
		}
		return "s"
	}()

	return &discordgo.MessageEmbed{
		Title:       "Queue",
		Description: fmt.Sprintf("Il y a actuellement %d message%s dans la queue.", n, s),
	}, ""
}

func init() {
	commands["restart"] = &bot.Command{
		Function: commandRestart, IsAdmin: true,
		HelpText: "Arrête et redémarre le bot.",
	}
	commands["update"] = &bot.Command{
		Function: commandUpdate, IsAdmin: true,
		HelpText: "Met à jour et redémarre le bot.",
	}
	commands["queue"] = &bot.Command{
		Function: commandQueue, IsAdmin: true,
		HelpText: "Renvoie le contenu de la queue.",
	}
}
