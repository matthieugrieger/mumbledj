/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
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
	sanitizedCommand := sanitize.HTML(command)
	if strings.Contains(sanitizedCommand, " ") {
		index := strings.Index(sanitizedCommand, " ")
		com, argument = sanitizedCommand[0:index], sanitizedCommand[(index+1):]
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
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Skip command
	case dj.conf.Aliases.SkipAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminSkip) {
			skip(user, username, false, false)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Skip playlist command
	case dj.conf.Aliases.SkipPlaylistAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminAddPlaylists) {
			skip(user, username, false, true)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Forceskip command
	case dj.conf.Aliases.AdminSkipAlias:
		if dj.HasPermission(username, true) {
			skip(user, username, true, false)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Playlist forceskip command
	case dj.conf.Aliases.AdminSkipPlaylistAlias:
		if dj.HasPermission(username, true) {
			skip(user, username, true, true)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Help command
	case dj.conf.Aliases.HelpAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminHelp) {
			help(user)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Volume command
	case dj.conf.Aliases.VolumeAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminVolume) {
			volume(user, username, argument)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Move command
	case dj.conf.Aliases.MoveAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminMove) {
			move(user, argument)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Reload command
	case dj.conf.Aliases.ReloadAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminReload) {
			reload(user)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Reset command
	case dj.conf.Aliases.ResetAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminReset) {
			reset(username)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Numsongs command
	case dj.conf.Aliases.NumSongsAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminNumSongs) {
			numSongs()
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Nextsong command
	case dj.conf.Aliases.NextSongAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminNextSong) {
			nextSong(user)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Currentsong command
	case dj.conf.Aliases.CurrentSongAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminCurrentSong) {
			currentSong(user)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Kill command
	case dj.conf.Aliases.KillAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminKill) {
			kill()
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	default:
		dj.SendPrivateMessage(user, COMMAND_DOESNT_EXIST_MSG)
	}
}

// Performs add functionality. Checks input URL for YouTube format, and adds
// the URL to the queue if the format matches.
func add(user *gumble.User, username, url string) {
	if url == "" {
		dj.SendPrivateMessage(user, NO_ARGUMENT_MSG)
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
			if err := dj.queue.AddItem(newSong); err == nil {
				dj.client.Self().Channel().Send(fmt.Sprintf(SONG_ADDED_HTML, username, newSong.title), false)
				if dj.queue.Len() == 1 && !dj.audioStream.IsPlaying() {
					if err := dj.queue.CurrentItem().(*Song).Download(); err == nil {
						dj.queue.CurrentItem().(*Song).Play()
					} else {
						dj.SendPrivateMessage(user, AUDIO_FAIL_MSG)
						dj.queue.CurrentItem().(*Song).Delete()
					}
				}
			}
		} else {
			// Check to see if we have a playlist URL instead.
			youtubePlaylistPattern := `https?:\/\/www\.youtube\.com\/playlist\?list=([\w-]+)`
			if re, err := regexp.Compile(youtubePlaylistPattern); err == nil {
				if re.MatchString(url) {
					if dj.HasPermission(username, dj.conf.Permissions.AdminAddPlaylists) {
						shortUrl = re.FindStringSubmatch(url)[1]
						newPlaylist := NewPlaylist(username, shortUrl)
						if dj.queue.AddItem(newPlaylist); err == nil {
							dj.client.Self().Channel().Send(fmt.Sprintf(PLAYLIST_ADDED_HTML, username, newPlaylist.title), false)
							if dj.queue.Len() == 1 && !dj.audioStream.IsPlaying() {
								if err := dj.queue.CurrentItem().(*Playlist).songs.CurrentItem().(*Song).Download(); err == nil {
									dj.queue.CurrentItem().(*Playlist).songs.CurrentItem().(*Song).Play()
								} else {
									dj.SendPrivateMessage(user, AUDIO_FAIL_MSG)
									dj.queue.CurrentItem().(*Playlist).songs.CurrentItem().(*Song).Delete()
								}
							}
						}
					} else {
						dj.SendPrivateMessage(user, NO_PLAYLIST_PERMISSION_MSG)
					}
				} else {
					dj.SendPrivateMessage(user, INVALID_URL_MSG)
				}
			}
		}
	}
}

// Performs skip functionality. Adds a skip to the skippers slice for the current song, and then
// evaluates if a skip should be performed. Both skip and forceskip are implemented here.
func skip(user *gumble.User, username string, admin, playlistSkip bool) {
	if dj.audioStream.IsPlaying() {
		if playlistSkip {
			if dj.queue.CurrentItem().ItemType() == "playlist" {
				if err := dj.queue.CurrentItem().AddSkip(username); err == nil {
					if admin {
						dj.client.Self().Channel().Send(ADMIN_PLAYLIST_SKIP_MSG, false)
					} else {
						dj.client.Self().Channel().Send(fmt.Sprintf(PLAYLIST_SKIP_ADDED_HTML, username), false)
					}
					if dj.queue.CurrentItem().SkipReached(len(dj.client.Self().Channel().Users())) || admin {
						dj.queue.CurrentItem().(*Playlist).skipped = true
						dj.client.Self().Channel().Send(PLAYLIST_SKIPPED_HTML, false)
						if err := dj.audioStream.Stop(); err != nil {
							panic(errors.New("An error occurred while stopping the current song."))
						}
					}
				}
			} else {
				dj.SendPrivateMessage(user, NO_PLAYLIST_PLAYING_MSG)
			}
		} else {
			var currentItem QueueItem
			if dj.queue.CurrentItem().ItemType() == "playlist" {
				currentItem = dj.queue.CurrentItem().(*Playlist).songs.CurrentItem()
			} else {
				currentItem = dj.queue.CurrentItem()
			}
			if err := currentItem.AddSkip(username); err == nil {
				if admin {
					dj.client.Self().Channel().Send(ADMIN_SONG_SKIP_MSG, false)
				} else {
					dj.client.Self().Channel().Send(fmt.Sprintf(SKIP_ADDED_HTML, username), false)
				}
				if currentItem.SkipReached(len(dj.client.Self().Channel().Users())) || admin {
					dj.client.Self().Channel().Send(SONG_SKIPPED_HTML, false)
					if err := dj.audioStream.Stop(); err != nil {
						panic(errors.New("An error occurred while stopping the current song."))
					}
				}
			}
		}
	} else {
		dj.SendPrivateMessage(user, NO_MUSIC_PLAYING_MSG)
	}
}

// Performs help functionality. Displays a list of valid commands.
func help(user *gumble.User) {
	dj.SendPrivateMessage(user, HELP_HTML)
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
				dj.SendPrivateMessage(user, fmt.Sprintf(NOT_IN_VOLUME_RANGE_MSG, dj.conf.Volume.LowestVolume, dj.conf.Volume.HighestVolume))
			}
		} else {
			dj.SendPrivateMessage(user, fmt.Sprintf(NOT_IN_VOLUME_RANGE_MSG, dj.conf.Volume.LowestVolume, dj.conf.Volume.HighestVolume))
		}
	}
}

// Performs move functionality. Determines if the supplied channel is valid and moves the bot
// to the channel if it is.
func move(user *gumble.User, channel string) {
	if channel == "" {
		dj.SendPrivateMessage(user, NO_ARGUMENT_MSG)
	} else {
		if dj.client.Channels().Find(channel) != nil {
			dj.client.Self().Move(dj.client.Channels().Find(channel))
		} else {
			dj.SendPrivateMessage(user, CHANNEL_DOES_NOT_EXIST_MSG)
		}
	}
}

// Performs reload functionality. Tells command submitter if the reload completed successfully.
func reload(user *gumble.User) {
	if err := loadConfiguration(); err == nil {
		dj.SendPrivateMessage(user, CONFIG_RELOAD_SUCCESS_MSG)
	}
}

// Performs reset functionality. Clears the song queue, stops playing audio, and deletes all
// remaining songs in the ~/.mumbledj/songs directory.
func reset(username string) {
	dj.queue.queue = dj.queue.queue[:0]
	if err := dj.audioStream.Stop(); err == nil {
		if err := deleteSongs(); err == nil {
			dj.client.Self().Channel().Send(fmt.Sprintf(QUEUE_RESET_HTML, username), false)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}

// Performs numsongs functionality. Uses the SongQueue traversal function to traverse the
// queue with a function call that increments a counter. Once finished, the bot outputs
// the number of songs in the queue to chat.
func numSongs() {
	songCount := 0
	dj.queue.Traverse(func(i int, item QueueItem) {
		songCount += 1
	})
	dj.client.Self().Channel().Send(fmt.Sprintf(NUM_SONGS_HTML, songCount), false)
}

// Performs nextsong functionality. Uses the SongQueue PeekNext function to peek at the next
// item if it exists. The user will then be sent a message containing the title and submitter
// of the next item if it exists.
func nextSong(user *gumble.User) {
	if song, err := dj.queue.PeekNext(); err != nil {
		dj.SendPrivateMessage(user, NO_SONG_NEXT_MSG)
	} else {
		dj.SendPrivateMessage(user, fmt.Sprintf(NEXT_SONG_HTML, song.title, song.submitter))
	}
}

// Performs currentsong functionality. Sends the user who submitted the currentsong command
// information about the song currently playing.
func currentSong(user *gumble.User) {
	if dj.audioStream.IsPlaying() {
		var currentItem *Song
		if dj.queue.CurrentItem().ItemType() == "playlist" {
			currentItem = dj.queue.CurrentItem().(*Playlist).songs.CurrentItem().(*Song)
		} else {
			currentItem = dj.queue.CurrentItem().(*Song)
		}
		dj.SendPrivateMessage(user, fmt.Sprintf(CURRENT_SONG_HTML, currentItem.title, currentItem.submitter))
	} else {
		dj.SendPrivateMessage(user, NO_MUSIC_PLAYING_MSG)
	}
}

// Performs kill functionality. First cleans the ~/.mumbledj/songs directory to get rid of any
// excess m4a files. The bot then safely disconnects from the server.
func kill() {
	if err := deleteSongs(); err != nil {
		panic(err)
	}
	if err := dj.client.Disconnect(); err == nil {
		fmt.Println("Kill successful. Goodbye!")
		os.Exit(0)
	} else {
		panic(errors.New("An error occurred while disconnecting from the server."))
	}
}

// Deletes songs from ~/.mumbledj/songs.
func deleteSongs() error {
	songsDir := fmt.Sprintf("%s/.mumbledj/songs", dj.homeDir)
	if err := os.RemoveAll(songsDir); err != nil {
		return errors.New("An error occurred while deleting the audio files.")
	} else {
		if err := os.Mkdir(songsDir, 0777); err != nil {
			return errors.New("An error occurred while recreating the songs directory.")
		}
		return nil
	}
	return nil
}
