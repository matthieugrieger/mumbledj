/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/add.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"layeh.com/gumble/gumble"
	"reik.pl/mumbledj/interfaces"
)

// AddCommand is a command that adds an audio track associated with a supported
// URL to the queue.
type AddCommand struct {
	mutex sync.Mutex
}

// Aliases returns the current aliases for the command.
func (c *AddCommand) Aliases() []string {
	return viper.GetStringSlice("commands.add.aliases")
}

// Description returns the description for the command.
func (c *AddCommand) Description() string {
	return viper.GetString("commands.add.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *AddCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.add.is_admin")
}

// Execute executes the command with the given user and arguments.
// Return value descriptions:
//    string: A message to be returned to the user upon successful execution.
//    bool:   Whether the message should be private or not. true = private,
//            false = public (sent to whole channel).
//    error:  An error message to be returned upon unsuccessful execution.
//            If no error has occurred, pass nil instead.
// Example return statement:
//    return "This is a private message!", true, nil
func (c *AddCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	var (
		allTracks      []interfaces.Track
		tracks         []interfaces.Track
		service        interfaces.Service
		err            error
		lastTrackAdded interfaces.Track
	)

	if len(args) == 0 {
		return "", true, errors.New(viper.GetString("commands.add.messages.no_url_error"))
	}

	for _, arg := range args {
		if service, err = DJ.GetService(arg); err == nil {
			tracks, err = service.GetTracks(arg, user)
			if err == nil {
				allTracks = append(allTracks, tracks...)
			}
		}
	}

	if len(allTracks) == 0 {
		return "", true, errors.New(viper.GetString("commands.add.messages.no_valid_tracks_error"))
	}

	numTooLong := 0
	numAdded := 0
	for _, track := range allTracks {
		c.mutex.Lock()
		if err = DJ.Queue.AppendTrack(track); err != nil {
			numTooLong++
		} else {
			numAdded++
			lastTrackAdded = track
		}
		c.mutex.Unlock()
	}

	if numAdded == 0 {
		return "", true, errors.New(viper.GetString("commands.add.messages.tracks_too_long_error"))
	} else if numAdded == 1 {
		return fmt.Sprintf(viper.GetString("commands.add.messages.one_track_added"),
			user.Name, lastTrackAdded.GetTitle(), lastTrackAdded.GetService()), false, nil
	}

	retString := fmt.Sprintf(viper.GetString("commands.add.messages.many_tracks_added"), user.Name, numAdded)
	if numTooLong != 0 {
		retString += fmt.Sprintf(viper.GetString("commands.add.messages.num_tracks_too_long"), numTooLong)
	}
	return retString, false, nil
}
