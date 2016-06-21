/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/nexttrack.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"

	"github.com/layeh/gumble/gumble"
	"github.com/spf13/viper"
)

// NextTrackCommand is a command that outputs information related to the next
// track in the queue (if one exists).
type NextTrackCommand struct{}

// Aliases returns the current aliases for the command.
func (c *NextTrackCommand) Aliases() []string {
	return viper.GetStringSlice("commands.nexttrack.aliases")
}

// Description returns the description for the command.
func (c *NextTrackCommand) Description() string {
	return viper.GetString("commands.nexttrack.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *NextTrackCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.nexttrack.is_admin")
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
func (c *NextTrackCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	length := DJ.Queue.Length()
	if length == 0 {
		return "", true, errors.New("There are no tracks in the queue")
	}
	if length == 1 {
		return "", true, errors.New("The current track is the only track in the queue")
	}

	nextTrack, _ := DJ.Queue.PeekNextTrack()

	return fmt.Sprintf("The next track is \"%s\", added by <b>%s</b>.",
		nextTrack.GetTitle(), nextTrack.GetSubmitter()), true, nil
}
