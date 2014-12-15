/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

import (
	"fmt"
	"strings"
)

func parseCommand(username, command string) {
	var com, argument string
	if strings.Contains(command, " ") {
		parsedCommand := strings.Split(command, " ")
		com, argument = parsedCommand[0], parsedCommand[1]
	} else {
		com = command
		argument = ""
	}
	
	switch com {
		case dj.conf.Aliases.AddAlias:
			success := add(username, argument)
			if success {
				fmt.Println("Add successful!")
			}
		case dj.conf.Aliases.SkipAlias:
			success := skip(username, false)
			if success {
				fmt.Println("Skip successful!")
			}
		case dj.conf.Aliases.AdminSkipAlias:
			success := skip(username, true)
			if success {
				fmt.Println("Forceskip successful!")
			}
		case dj.conf.Aliases.VolumeAlias:
			success := volume(username, argument)
			if success {
				fmt.Println("Volume change successful!")
			}
		case dj.conf.Aliases.MoveAlias:
			success := move(username, argument)
			if success {
				fmt.Println("Move successful!")
			}
		case dj.conf.Aliases.ReloadAlias:
			err := loadConfiguration()
			if err == nil {
				fmt.Println("Reload successful!")
			}
		case dj.conf.Aliases.KillAlias:
			success := kill(username)
			if success {
				fmt.Println("Kill successful!")
			}
	}
}

func add(user, url string) bool {
	return true
}

func skip(user string, admin bool) bool {
	return true
}

func volume(user, value string) bool {
	return true
}

func move(user, channel string) bool {
	return true
}

func kill(user string) bool {
	return true
}
