/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

import (
	"fmt"
	"github.com/layeh/gumble/gumble"
	"regexp"
	"strings"
)

func parseCommand(user *gumble.User, username, command string) {
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
		if dj.HasPermission(username, dj.conf.Permissions.AdminAdd) {
			if argument == "" {
				user.Send(NO_ARGUMENT_MSG)
			} else {
				success := add(username, argument)
				if success {
					fmt.Println("Add successful!")
					dj.client.Me().Channel().Send(fmt.Sprintf("%s has added a song to the queue.", username))
				}
			}
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	case dj.conf.Aliases.SkipAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminSkip) {
			success := skip(username, false)
			if success {
				fmt.Println("Skip successful!")
			}
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	case dj.conf.Aliases.AdminSkipAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminSkip) {
			success := skip(username, true)
			if success {
				fmt.Println("Forceskip successful!")
			}
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	case dj.conf.Aliases.VolumeAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminVolume) {
			if argument == "" {
				dj.client.Self().Channel().Send(fmt.Sprintf(CUR_VOLUME_HTML, dj.conf.Volume.DefaultVolume), false)
			} else {
				success := volume(username, argument)
				if success {
					fmt.Println("Skip successful!")
				}
			}
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	case dj.conf.Aliases.MoveAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminMove) {
			if argument == "" {
				user.Send(NO_ARGUMENT_MSG)
			} else {
				success := move(username, argument)
				if success {
					fmt.Printf("%s has been moved to %s.", dj.client.Self().Name(), argument)
				}
			}
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	case dj.conf.Aliases.ReloadAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminReload) {
			err := loadConfiguration()
			if err == nil {
				user.Send(CONFIG_RELOAD_SUCCESS_MSG)
			} else {
				panic(err)
			}
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	case dj.conf.Aliases.KillAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminKill) {
			success := kill(username)
			if success {
				fmt.Println("Kill successful!")
			}
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	default:
		user.Send(COMMAND_DOESNT_EXIST_MSG)
	}
}

func add(user, url string) bool {
	youtubePatterns := []string{
		`https?:\/\/www\.youtube\.com\/watch\?v=([\w-]+)`,
		`https?:\/\/youtube\.com\/watch\?v=([\w-]+)`,
		`https?:\/\/youtu.be\/([\w-]+)`,
		`https?:\/\/youtube.com\/v\/([\w-]+)`,
		`https?:\/\/www.youtube.com\/v\/([\w-]+)`,
	}
	matchFound := false

	for _, pattern := range youtubePatterns {
		re, err := regexp.Compile(pattern)
		if err == nil {
			if re.MatchString(url) {
				matchFound = true
				break
			}
		}
	}

	if matchFound {
		if dj.queue.AddSong(NewSong(user, url)) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
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
