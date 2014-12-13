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
		case "add":
			success := add(username, argument)
			if success {
				fmt.Println("Add successful!")
			}
		case "skip":
			success := skip(username, false)
			if success {
				fmt.Println("Skip successful!")
			}
		case "forceskip":
			success := skip(username, true)
			if success {
				fmt.Println("Forceskip successful!")
			}
		case "volume":
			success := volume(username, argument)
			if success {
				fmt.Println("Volume change successful!")
			}
		case "move":
			success := move(username, argument)
			if success {
				fmt.Println("Move successful!")
			}
		case "reload":
			conf, err := loadConfiguration()
			if err == nil {
				dj.conf = conf
				fmt.Println("Reload successful!")
			}
		case "kill":
			success := kill(username)
			if success {
				fmt.Println("Kill successful!")
			}
		case "test":
			fmt.Printf("Title: %s\n", dj.conf.title)
	}
}

func add(user, url string) bool {
	fmt.Println("Add requested!")
	return true
}

func skip(user string, admin bool) bool {
	if admin {
		fmt.Println("Admin skip requested!")
	} else {
		fmt.Println("Skip requested!")
	}
	return true
}

func volume(user, value string) bool {
	fmt.Println("Volume change requested!")
	return true
}

func move(user, channel string) bool {
	fmt.Println("Move requested!")
	return true
}

func kill(user string) bool {
	fmt.Println("Kill requested!")
	return true
}
