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
				if songTitle, err := add(username, argument); err == nil {
					dj.client.Self().Channel().Send(fmt.Sprintf(SONG_ADDED_HTML, username, songTitle), false)
					if dj.queue.Len() == 1 && !dj.audioStream.IsPlaying() {
						dj.currentSong = dj.queue.NextSong()
						if err := dj.currentSong.Download(); err == nil {
							dj.currentSong.Play()
						} else {
							panic(err)
						}
					}
				} else {
					user.Send(INVALID_URL_MSG)
				}
			}
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	case dj.conf.Aliases.SkipAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminSkip) {
			if err := skip(username, false); err == nil {
				dj.client.Self().Channel().Send(fmt.Sprintf(SKIP_ADDED_HTML, username), false)
			}
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	case dj.conf.Aliases.AdminSkipAlias:
		if dj.HasPermission(username, true) {
			if err := skip(username, true); err == nil {
				dj.client.Self().Channel().Send(ADMIN_SONG_SKIP_MSG, false)
			}
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	case dj.conf.Aliases.VolumeAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminVolume) {
			if argument == "" {
				dj.client.Self().Channel().Send(fmt.Sprintf(CUR_VOLUME_HTML, dj.conf.Volume.DefaultVolume), false)
			} else {
				if err := volume(username, argument); err == nil {
					dj.client.Self().Channel().Send(fmt.Sprintf(VOLUME_SUCCESS_HTML, username, argument), false)
				} else {
					user.Send(NOT_IN_VOLUME_RANGE_MSG)
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
				if err := move(argument); err == nil {
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
			if err := kill(); err == nil {
				fmt.Println("Kill successful. Goodbye!")
				os.Exit(0)
			} else {
				user.Send(KILL_ERROR_MSG)
			}
		} else {
			user.Send(NO_PERMISSION_MSG)
		}
	default:
		user.Send(COMMAND_DOESNT_EXIST_MSG)
	}
}

func add(user, url string) (string, error) {
	youtubePatterns := []string{
		`https?:\/\/www\.youtube\.com\/watch\?v=([\w-]+)`,
		`https?:\/\/youtube\.com\/watch\?v=([\w-]+)`,
		`https?:\/\/youtu.be\/([\w-]+)`,
		`https?:\/\/youtube.com\/v\/([\w-]+)`,
		`https?:\/\/www.youtube.com\/v\/([\w-]+)`,
	}
	matchFound := false

	for _, pattern := range youtubePatterns {
		if re, err := regexp.Compile(pattern); err == nil {
			if re.MatchString(url) {
				matchFound = true
				break
			}
		}
	}

	if matchFound {
		urlMatch := strings.Split(url, "=")
		shortUrl := urlMatch[1]
		newSong := NewSong(user, shortUrl)
		if err := dj.queue.AddSong(newSong); err == nil {
			return newSong.title, nil
		} else {
			return "", errors.New("Could not add the Song to the queue.")
		}
	} else {
		return "", errors.New("The URL provided did not match a YouTube URL.")
	}
}

func skip(user string, admin bool) error {
	if err := dj.currentSong.AddSkip(user); err == nil {
		if dj.currentSong.SkipReached(len(dj.client.Self().Channel().Users())) || admin {
			if err := dj.audioStream.Stop(); err == nil {
				dj.OnSongFinished()
				return nil
			} else {
				return errors.New("An error occurred while stopping the current song.")
			}
		} else {
			return errors.New("Not enough skips have been reached to skip the song.")
		}
	} else {
		return errors.New("An error occurred while adding a skip to the current song.")
	}
}

func volume(user, value string) error {
	if parsedVolume, err := strconv.ParseFloat(value, 32); err == nil {
		newVolume := float32(parsedVolume)
		if newVolume >= dj.conf.Volume.LowestVolume && newVolume <= dj.conf.Volume.HighestVolume {
			dj.conf.Volume.DefaultVolume = newVolume
			return nil
		} else {
			return errors.New("The volume supplied was not in the allowed range.")
		}
	} else {
		return errors.New("An error occurred while parsing the volume string.")
	}
}

func move(channel string) error {
	if dj.client.Channels().Find(channel) != nil {
		dj.client.Self().Move(dj.client.Channels().Find(channel))
		return nil
	} else {
		return errors.New("The channel provided does not exist.")
	}
}

func kill() error {
	songsDir := fmt.Sprintf("%s/.mumbledj/songs", dj.homeDir)
	if err := os.RemoveAll(songsDir); err != nil {
		return errors.New("An error occurred while deleting the audio files.")
	} else {
		if err := os.Mkdir(songsDir, 0777); err != nil {
			return errors.New("An error occurred while recreating the songs directory.")
		}
	}
	if err := dj.client.Disconnect(); err == nil {
		return nil
	} else {
		return errors.New("An error occurred while disconnecting from the server.")
	}
}
