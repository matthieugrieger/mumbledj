/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/addnext.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"

	"github.com/layeh/gumble/gumble"
	"github.com/matthieugrieger/mumbledj/interfaces"
	"github.com/spf13/viper"
)

// AddNextCommand is a command that adds an audio track associated with a supported
// URL to the queue as the next item.
type AddNextCommand struct{}

// Aliases returns the current aliases for the command.
func (c *AddNextCommand) Aliases() []string {
	return viper.GetStringSlice("commands.addnext.aliases")
}

// Description returns the description for the command.
func (c *AddNextCommand) Description() string {
	return viper.GetString("commands.addnext.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *AddNextCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.addnext.is_admin")
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
func (c *AddNextCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
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
	// We must loop backwards here to preserve the track order when inserting tracks.
	for i := len(allTracks) - 1; i >= 0; i-- {
		insertIndex := 1
		if DJ.Queue.Length() == 0 {
			insertIndex = 0
		}
		if err = DJ.Queue.InsertTrack(insertIndex, allTracks[i]); err != nil {
			numTooLong++
		} else {
			numAdded++
			lastTrackAdded = allTracks[i]
		}
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
