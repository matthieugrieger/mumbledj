/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

import (
	"fmt"
	"github.com/kennygrant/sanitize"
	"github.com/layeh/gumble/gumble"
	"regexp"
	"strconv"
	"strings"
)

func parseCommand(user *gumble.User, username, command string) {
	var com, argument string
	if strings.Contains(command, " ") {
		sanitizedCommand := sanitize.HTML(command)
		parsedCommand := strings.Split(sanitizedCommand, " ")
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
					// TODO: Replace this message with a more informative one.
					dj.client.Self().Channel().Send(fmt.Sprintf("%s has added a song to the queue.", username), false)
				} else {
					user.Send(INVALID_URL_MSG)
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
				dj.client.Self().Channel().Send(fmt.Sprintf(SKIP_ADDED_HTML, username), false)
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
					fmt.Println("Volume change successful!")
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
				} else {
					user.Send(CHANNEL_DOES_NOT_EXIST_MSG)
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
		urlMatch := strings.Split(url, "=")
		shortUrl := urlMatch[1]
		if dj.queue.AddSong(NewSong(user, shortUrl)) {
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
	parsedVolume, err := strconv.ParseFloat(value, 32)
	if err == nil {
		newVolume := float32(parsedVolume)
		if newVolume >= dj.conf.Volume.LowestVolume && newVolume <= dj.conf.Volume.HighestVolume {
			dj.conf.Volume.DefaultVolume = newVolume
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func move(user, channel string) bool {
	return true
}

func kill(user string) bool {
	return true
}
