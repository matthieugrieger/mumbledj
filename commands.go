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
	"os"
	"strconv"
	"strings"

	"github.com/layeh/gumble/gumble"
)

// parseCommand views incoming chat messages and determines if there is a valid command within them.
// If a command exists, the arguments (if any) will be parsed and sent to the appropriate helper
// function to perform the command's task.
func parseCommand(user *gumble.User, username, command string) {
	var com, argument string
	split := strings.Split(command, "\n")
	splitString := split[0]
	if strings.Contains(splitString, " ") {
		index := strings.Index(splitString, " ")
		com, argument = splitString[0:index], splitString[(index+1):]
	} else {
		com = command
		argument = ""
	}

	switch com {
	// Add commands
	case dj.conf.Aliases.AddAlias, dj.conf.Aliases.AddAlias2:
		if dj.HasPermission(username, dj.conf.Permissions.AdminAdd) {
			add(user, argument)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Skip command
	case dj.conf.Aliases.SkipAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminSkip) {
			skip(user, false, false)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Skip playlist command
	case dj.conf.Aliases.SkipPlaylistAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminAddPlaylists) {
			skip(user, false, true)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Forceskip command
	case dj.conf.Aliases.AdminSkipAlias:
		if dj.HasPermission(username, true) {
			skip(user, true, false)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Playlist forceskip command
	case dj.conf.Aliases.AdminSkipPlaylistAlias:
		if dj.HasPermission(username, true) {
			skip(user, true, true)
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
	// Volume commands
	case dj.conf.Aliases.VolumeAlias, dj.conf.Aliases.VolumeAlias2:
		if dj.HasPermission(username, dj.conf.Permissions.AdminVolume) {
			volume(user, argument)
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
	// Setcomment command
	case dj.conf.Aliases.SetCommentAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminSetComment) {
			setComment(user, argument)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Numcached command
	case dj.conf.Aliases.NumCachedAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminNumCached) {
			numCached(user)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Cachesize command
	case dj.conf.Aliases.CacheSizeAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminCacheSize) {
			cacheSize(user)
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
	// Shuffle command
	case dj.conf.Aliases.ShuffleAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminShuffle) {
			shuffleSongs(user, username)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	// Shuffleon command
  	case dj.conf.Aliases.ShuffleOnAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminShuffleToggle) {
			toggleAutomaticShuffle(true, user, username)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
  	// Shuffleoff command
	case dj.conf.Aliases.ShuffleOffAlias:
		if dj.HasPermission(username, dj.conf.Permissions.AdminShuffleToggle) {
			toggleAutomaticShuffle(false, user, username)
		} else {
			dj.SendPrivateMessage(user, NO_PERMISSION_MSG)
		}
	default:
		dj.SendPrivateMessage(user, COMMAND_DOESNT_EXIST_MSG)
	}
}


// add performs !add functionality. Checks input URL for service, and adds
// the URL to the queue if the format matches.
func add(user *gumble.User, url string) error {
	if url == "" {
		dj.SendPrivateMessage(user, NO_ARGUMENT_MSG)
		return errors.New("NO_ARGUMENT")
	} else {
		err := FindServiceAndAdd(user, url)
		if err != nil {
			dj.SendPrivateMessage(user, err.Error())
		}
		return err
	}
}

// skip performs !skip functionality. Adds a skip to the skippers slice for the current song, and then
// evaluates if a skip should be performed. Both skip and forceskip are implemented here.
func skip(user *gumble.User, admin, playlistSkip bool) {
	if dj.audioStream.IsPlaying() {
		if playlistSkip {
			if dj.queue.CurrentSong().Playlist() != nil {
				if err := dj.queue.CurrentSong().Playlist().AddSkip(user.Name); err == nil {
					submitterSkipped := false
					if admin {
						dj.client.Self.Channel.Send(ADMIN_PLAYLIST_SKIP_MSG, false)
					} else if dj.queue.CurrentSong().Submitter() == user.Name {
						dj.client.Self.Channel.Send(fmt.Sprintf(PLAYLIST_SUBMITTER_SKIP_HTML, user.Name), false)
						submitterSkipped = true
					} else {
						dj.client.Self.Channel.Send(fmt.Sprintf(PLAYLIST_SKIP_ADDED_HTML, user.Name), false)
					}
					if submitterSkipped || dj.queue.CurrentSong().Playlist().SkipReached(len(dj.client.Self.Channel.Users)) || admin {
						id := dj.queue.CurrentSong().Playlist().ID()
						dj.queue.CurrentSong().Playlist().DeleteSkippers()
						for i := 0; i < len(dj.queue.queue); i++ {
							if dj.queue.queue[i].Playlist() != nil {
								if dj.queue.queue[i].Playlist().ID() == id {
									dj.queue.queue = append(dj.queue.queue[:i], dj.queue.queue[i+1:]...)
									i--
								}
							}
						}
						if dj.queue.Len() != 0 {
							// Set dontSkip to true to avoid audioStream.Stop() callback skipping the new first song.
							dj.queue.CurrentSong().SetDontSkip(true)
						}
						if !(submitterSkipped || admin) {
							dj.client.Self.Channel.Send(PLAYLIST_SKIPPED_HTML, false)
						}
						if err := dj.audioStream.Stop(); err != nil {
							panic(errors.New("An error occurred while stopping the current song."))
						}
					}
				}
			} else {
				dj.SendPrivateMessage(user, NO_PLAYLIST_PLAYING_MSG)
			}
		} else {
			if err := dj.queue.CurrentSong().AddSkip(user.Name); err == nil {
				submitterSkipped := false
				if admin {
					dj.client.Self.Channel.Send(ADMIN_SONG_SKIP_MSG, false)
				} else if dj.queue.CurrentSong().Submitter() == user.Name {
					dj.client.Self.Channel.Send(fmt.Sprintf(SUBMITTER_SKIP_HTML, user.Name), false)
					submitterSkipped = true
				} else {
					dj.client.Self.Channel.Send(fmt.Sprintf(SKIP_ADDED_HTML, user.Name), false)
				}
				if submitterSkipped || dj.queue.CurrentSong().SkipReached(len(dj.client.Self.Channel.Users)) || admin {
					if !(submitterSkipped || admin) {
						dj.client.Self.Channel.Send(SONG_SKIPPED_HTML, false)
					}
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

// help performs !help functionality. Displays a list of valid commands.
func help(user *gumble.User) {
	dj.SendPrivateMessage(user, HELP_HTML)
}

// volume performs !volume functionality. Checks input value against LowestVolume and HighestVolume from
// config to determine if the volume should be applied. If in the correct range, the new volume
// is applied and is immediately in effect.
func volume(user *gumble.User, value string) {
	if value == "" {
		dj.client.Self.Channel.Send(fmt.Sprintf(CUR_VOLUME_HTML, dj.audioStream.Volume), false)
	} else {
		if parsedVolume, err := strconv.ParseFloat(value, 32); err == nil {
			newVolume := float32(parsedVolume)
			if newVolume >= dj.conf.Volume.LowestVolume && newVolume <= dj.conf.Volume.HighestVolume {
				dj.audioStream.Volume = newVolume
				dj.client.Self.Channel.Send(fmt.Sprintf(VOLUME_SUCCESS_HTML, user.Name, dj.audioStream.Volume), false)
			} else {
				dj.SendPrivateMessage(user, fmt.Sprintf(NOT_IN_VOLUME_RANGE_MSG, dj.conf.Volume.LowestVolume, dj.conf.Volume.HighestVolume))
			}
		} else {
			dj.SendPrivateMessage(user, fmt.Sprintf(NOT_IN_VOLUME_RANGE_MSG, dj.conf.Volume.LowestVolume, dj.conf.Volume.HighestVolume))
		}
	}
}

// move performs !move functionality. Determines if the supplied channel is valid and moves the bot
// to the channel if it is.
func move(user *gumble.User, channel string) {
	if channel == "" {
		dj.SendPrivateMessage(user, NO_ARGUMENT_MSG)
	} else {
		if channels := strings.Split(channel, "/"); dj.client.Channels.Find(channels...) != nil {
			dj.client.Self.Move(dj.client.Channels.Find(channels...))
		} else {
			dj.SendPrivateMessage(user, CHANNEL_DOES_NOT_EXIST_MSG+" "+channel)
		}
	}
}

// reload performs !reload functionality. Tells command submitter if the reload completed successfully.
func reload(user *gumble.User) {
	if err := loadConfiguration(); err == nil {
		dj.SendPrivateMessage(user, CONFIG_RELOAD_SUCCESS_MSG)
	}
}

// reset performs !reset functionality. Clears the song queue, stops playing audio, and deletes all
// remaining songs in the ~/.mumbledj/songs directory.
func reset(username string) {
	dj.queue.queue = dj.queue.queue[:0]
	if dj.audioStream.IsPlaying() {
		if err := dj.audioStream.Stop(); err != nil {
			panic(err)
		}
	}
	if err := deleteSongs(); err == nil {
		dj.client.Self.Channel.Send(fmt.Sprintf(QUEUE_RESET_HTML, username), false)
	} else {
		panic(err)
	}
}

// numSongs performs !numsongs functionality. Uses the SongQueue traversal function to traverse the
// queue with a function call that increments a counter. Once finished, the bot outputs
// the number of songs in the queue to chat.
func numSongs() {
	songCount := 0
	dj.queue.Traverse(func(i int, song Song) {
		songCount++
	})
	dj.client.Self.Channel.Send(fmt.Sprintf(NUM_SONGS_HTML, songCount), false)
}

// nextSong performs !nextsong functionality. Uses the SongQueue PeekNext function to peek at the next
// item if it exists. The user will then be sent a message containing the title and submitter
// of the next item if it exists.
func nextSong(user *gumble.User) {
	if song, err := dj.queue.PeekNext(); err != nil {
		dj.SendPrivateMessage(user, NO_SONG_NEXT_MSG)
	} else {
		dj.SendPrivateMessage(user, fmt.Sprintf(NEXT_SONG_HTML, song.Title(), song.Submitter()))
	}
}

// currentSong performs !currentsong functionality. Sends the user who submitted the currentsong command
// information about the song currently playing.
func currentSong(user *gumble.User) {
	if dj.audioStream.IsPlaying() {
		if dj.queue.CurrentSong().Playlist() == nil {
			dj.SendPrivateMessage(user, fmt.Sprintf(CURRENT_SONG_HTML, dj.queue.CurrentSong().Title(), dj.queue.CurrentSong().Submitter()))
		} else {
			dj.SendPrivateMessage(user, fmt.Sprintf(CURRENT_SONG_PLAYLIST_HTML, dj.queue.CurrentSong().Title(),
				dj.queue.CurrentSong().Submitter(), dj.queue.CurrentSong().Playlist().Title()))
		}
	} else {
		dj.SendPrivateMessage(user, NO_MUSIC_PLAYING_MSG)
	}
}

// setComment performs !setcomment functionality. Sets the bot's comment to whatever text is supplied in the argument.
func setComment(user *gumble.User, comment string) {
	dj.client.Self.SetComment(comment)
	dj.SendPrivateMessage(user, COMMENT_UPDATED_MSG)
}

// numCached performs !numcached functionality. Displays the number of songs currently cached on disk at ~/.mumbledj/songs.
func numCached(user *gumble.User) {
	if dj.conf.Cache.Enabled {
		dj.cache.Update()
		dj.SendPrivateMessage(user, fmt.Sprintf(NUM_CACHED_MSG, dj.cache.NumSongs))
	} else {
		dj.SendPrivateMessage(user, CACHE_NOT_ENABLED_MSG)
	}
}

// cacheSize performs !cachesize functionality. Displays the total file size of the cached audio files.
func cacheSize(user *gumble.User) {
	if dj.conf.Cache.Enabled {
		dj.cache.Update()
		dj.SendPrivateMessage(user, fmt.Sprintf(CACHE_SIZE_MSG, float64(dj.cache.TotalFileSize/1048576)))
	} else {
		dj.SendPrivateMessage(user, CACHE_NOT_ENABLED_MSG)
	}
}

// kill performs !kill functionality. First cleans the ~/.mumbledj/songs directory to get rid of any
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

// deleteSongs deletes songs from ~/.mumbledj/songs.
func deleteSongs() error {
	songsDir := fmt.Sprintf("%s/.mumbledj/songs", dj.homeDir)
	if err := os.RemoveAll(songsDir); err != nil {
		return errors.New("An error occurred while deleting the audio files.")
	}
	if err := os.Mkdir(songsDir, 0777); err != nil {
		return errors.New("An error occurred while recreating the songs directory.")
	}
	return nil
}

// shuffles the song list
func shuffleSongs(user *gumble.User, username string) {
	if dj.queue.Len() > 1 {
		dj.queue.ShuffleSongs()
		dj.client.Self.Channel.Send(fmt.Sprintf(SHUFFLE_SUCCESS_MSG, username), false)
	} else {
		dj.SendPrivateMessage(user, CANT_SHUFFLE_MSG)
	}
}

// handles toggling of automatic shuffle playing
func toggleAutomaticShuffle(activate bool, user *gumble.User, username string){
	if (dj.conf.General.AutomaticShuffleOn != activate){
		dj.conf.General.AutomaticShuffleOn = activate
		if (activate){
			dj.client.Self.Channel.Send(fmt.Sprintf(SHUFFLE_ON_MESSAGE, username), false)
		} else{
			dj.client.Self.Channel.Send(fmt.Sprintf(SHUFFLE_OFF_MESSAGE, username), false)
		}
	} else if (activate){
		dj.SendPrivateMessage(user, SHUFFLE_ACTIVATED_ERROR_MESSAGE)
	} else{
		dj.SendPrivateMessage(user, SHUFFLE_DEACTIVATED_ERROR_MESSAGE)
	}
}