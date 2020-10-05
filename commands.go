package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type User struct {
	lang string
	read []string
}

func handleCommands(msg *discordgo.MessageCreate) {
	if strings.HasPrefix(msg.Content, "!lang=") {

	} else if strings.HasPrefix(msg.Content, "!read=") {

	}
}

func checkLang(s string) bool {
	switch s {
	case "en", "ru":
		return true
	default:
		return false
	}
}
