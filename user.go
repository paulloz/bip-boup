package main

import (
	"github.com/bwmarrin/discordgo"
)

func isUserAdmin(user *discordgo.User) bool {
	for _, adminID := range Bot.Admins {
		if user.ID == adminID {
			return true
		}
	}
	return false
}
