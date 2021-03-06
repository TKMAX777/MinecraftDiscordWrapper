package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandContent put minecraft command content
type CommandContent struct {
	Command string
	Options string
}

// MinecraftCommand handles minecraft function
type MinecraftCommand struct {
	sendChannel chan CommandContent
	idRegExp    *regexp.Regexp
}

// Handler handle say commands sent to discord
func (c *MinecraftCommand) Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID != Settings.Discord.GuildID || m.ChannelID != Settings.Discord.ChannelID {
		return
	}

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if m.Message.Content == "" {
		return
	}

	userDict, err := ReadNameDict()
	if err != nil {
		return
	}

	user, ok := userDict.findUserFromDiscordID(m.Author.ID)
	if !ok {
		user, ok = userDict.findUserFromDiscordID("Default")
		if !ok {
			return
		}
		user.Name = "Unknown"
	}

	for _, text := range strings.Split(m.Message.Content, "\n") {
		var command CommandContent
		var msg = strings.Split(text, " ")

		var permissions = GetPermissions(user.PermissionCode)
		command.Command, ok = permissions[msg[0]]
		if !ok {
			_, ok = permissions["say"]
			if !ok || !user.SendAllMessages {
				return
			}
			msg = append([]string{"say"}, msg...)
			command.Command = "/say"
		}

		if len(msg) < 2 {
			return
		}

		if msg[1] == ";" {
			msg = msg[:1]
		}

		switch command.Command {
		case "/msg", "/say":
			msg[1] = fmt.Sprintf("[%s]%s", user.Name, msg[1])
			command.Options = strings.Join(msg[1:], " ")

			for _, match := range c.idRegExp.FindAllStringSubmatch(command.Options, -1) {
				fmt.Println(match)
				u, ok := userDict.findUserFromDiscordID(match[1])
				if !ok {
					continue
				}
				command.Options = strings.ReplaceAll(command.Options, "!"+u.DiscordID, u.Name)
			}

		default:
			command.Options = strings.Join(msg[1:], " ")
		}

		fmt.Printf("[Discord]%v\n", text)

		c.sendChannel <- command
	}

	return
}
