/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"
	"fmt"
	"github.com/kennygrant/sanitize"
	"github.com/layeh/gumble/gumble"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Called on text message event. Checks the message for a command string, and processes it accordingly if
// it contains a command.
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
	// Add command
	case dj.conf.Aliases.AddAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminAdd) {
			add(user, username, argument)
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	// Skip command
	case dj.conf.Aliases.SkipAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminSkip) {
			skip(username, false)
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	// Forceskip command
	case dj.conf.Aliases.AdminSkipAlias:
		if dj.HasPermission(username, true) {
			skip(username, true)
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	// Volume command
	case dj.conf.Aliases.VolumeAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminVolume) {
			volume(user, username, argument)
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	// Move command
	case dj.conf.Aliases.MoveAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminMove) {
			move(user, argument)
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	// Reload command
	case dj.conf.Aliases.ReloadAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminReload) {
			reload(user)
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	// Kill command
	case dj.conf.Aliases.KillAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminKill) {
			kill()
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	default:
		user.Send(COMMAND_DOESNT_EXIST_MSG)
	}
}

// Performs add functionality. Checks input URL for YouTube format, and adds
// the URL to the queue if the format matches.
func add(user *gumble.User, username, url string) {
	if url == "" {
		user.Send(NO_ARGUMENT_MSG)
	} else {
		youtubePatterns := []string{
			`https?:\/\/www\.youtube\.com\/watch\?v=([\w-]+)`,
			`https?:\/\/youtube\.com\/watch\?v=([\w-]+)`,
			`https?:\/\/youtu.be\/([\w-]+)`,
			`https?:\/\/youtube.com\/v\/([\w-]+)`,
			`https?:\/\/www.youtube.com\/v\/([\w-]+)`,
		}
		matchFound := false
		shortUrl := ""

		for _, pattern := range youtubePatterns {
			if re, err := regexp.Compile(pattern); err == nil {
				if re.MatchString(url) {
					matchFound = true
					shortUrl = re.FindStringSubmatch(url)[1]
					break
				}
			}
		}

		if matchFound {
			newSong := NewSong(username, shortUrl)
			if err := dj.queue.AddSong(newSong); err == nil {
				dj.client.Self().Channel().Send(fmt.Sprintf(SONG_ADDED_HTML, username, newSong.title), false)
				if dj.queue.Len() == 1 && !dj.audioStream.IsPlaying() {
					dj.currentSong = dj.queue.NextSong()
					if err := dj.currentSong.Download(); err == nil {
						dj.currentSong.Play()
					} else {
						user.Send(AUDIO_FAIL_MSG)
						dj.currentSong.Delete()
					}
				}
			} else {
				panic(errors.New("Could not add the Song to the queue."))
			}
		} else {
			user.Send(INVALID_URL_MSG)
		}
	}
}

// Performs skip functionality. Adds a skip to the skippers slice for the current song, and then
// evaluates if a skip should be performed. Both skip and forceskip are implemented here.
func skip(user string, admin bool) {
	if err := dj.currentSong.AddSkip(user); err == nil {
		if admin {
			dj.client.Self().Channel().Send(ADMIN_SONG_SKIP_MSG, false)
		} else {
			dj.client.Self().Channel().Send(fmt.Sprintf(SKIP_ADDED_HTML, user), false)
		}
		if dj.currentSong.SkipReached(len(dj.client.Self().Channel().Users())) || admin {
			dj.client.Self().Channel().Send(SONG_SKIPPED_HTML, false)
			if err := dj.audioStream.Stop(); err == nil {
				dj.OnSongFinished()
			} else {
				panic(errors.New("An error occurred while stopping the current song."))
			}
		}
	} else {
		panic(errors.New("An error occurred while adding a skip to the current song."))
	}
}

// Performs volume functionality. Checks input value against LowestVolume and HighestVolume from
// config to determine if the volume should be applied. If in the correct range, the new volume
// is applied and is immediately in effect.
func volume(user *gumble.User, username, value string) {
	if value == "" {
		dj.client.Self().Channel().Send(fmt.Sprintf(CUR_VOLUME_HTML, dj.audioStream.Volume()), false)
	} else {
		if parsedVolume, err := strconv.ParseFloat(value, 32); err == nil {
			newVolume := float32(parsedVolume)
			if newVolume >= dj.conf.Volume.LowestVolume && newVolume <= dj.conf.Volume.HighestVolume {
				dj.audioStream.SetVolume(newVolume)
				dj.client.Self().Channel().Send(fmt.Sprintf(VOLUME_SUCCESS_HTML, username, dj.audioStream.Volume()), false)
			} else {
				user.Send(fmt.Sprintf(NOT_IN_VOLUME_RANGE_MSG, dj.conf.Volume.LowestVolume, dj.conf.Volume.HighestVolume))
			}
		} else {
			user.Send(fmt.Sprintf(NOT_IN_VOLUME_RANGE_MSG, dj.conf.Volume.LowestVolume, dj.conf.Volume.HighestVolume))
		}
	}
}

// Performs move functionality. Determines if the supplied channel is valid and moves the bot
// to the channel if it is.
func move(user *gumble.User, channel string) {
	if channel == "" {
		user.Send(NO_ARGUMENT_MSG)
	} else {
		if dj.client.Channels().Find(channel) != nil {
			dj.client.Self().Move(dj.client.Channels().Find(channel))
		} else {
			user.Send(CHANNEL_DOES_NOT_EXIST_MSG)
		}
	}
}

// Performs reload functionality. Tells command submitter if the reload completed successfully.
func reload(user *gumble.User) {
	if err := loadConfiguration(); err == nil {
		user.Send(CONFIG_RELOAD_SUCCESS_MSG)
	} else {
		panic(err)
	}
}

// Performs kill functionality. First cleans the ~/.mumbledj/songs directory to get rid of any
// excess m4a files. The bot then safely disconnects from the server.
func kill() {
	songsDir := fmt.Sprintf("%s/.mumbledj/songs", dj.homeDir)
	if err := os.RemoveAll(songsDir); err != nil {
		panic(errors.New("An error occurred while deleting the audio files."))
	} else {
		if err := os.Mkdir(songsDir, 0777); err != nil {
			panic(errors.New("An error occurred while recreating the songs directory."))
		}
	}
	if err := dj.client.Disconnect(); err == nil {
		fmt.Println("Kill successful. Goodbye!")
		os.Exit(0)
	} else {
		panic(errors.New("An error occurred while disconnecting from the server."))
	}
}
