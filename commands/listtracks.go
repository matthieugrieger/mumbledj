/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/currenttrack.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/layeh/gumble/gumble"
	"github.com/matthieugrieger/mumbledj/interfaces"
	"github.com/spf13/viper"
)

// ListTracksCommand is a command that lists the tracks that are currently
// in the queue.
type ListTracksCommand struct{}

// Aliases returns the current aliases for the command.
func (c *ListTracksCommand) Aliases() []string {
	return viper.GetStringSlice("commands.listtracks.aliases")
}

// Description returns the description for the command.
func (c *ListTracksCommand) Description() string {
	return viper.GetString("commands.listtracks.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *ListTracksCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.listtracks.is_admin")
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
func (c *ListTracksCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if DJ.Queue.Length() == 0 {
		return "", true, errors.New(viper.GetString("commands.common_messages.no_tracks_error"))
	}

	numTracksToList := DJ.Queue.Length()
	if len(args) != 0 {
		if parsedNum, err := strconv.Atoi(args[0]); err == nil {
			numTracksToList = parsedNum
		} else {
			return "", true, errors.New(viper.GetString("commands.listtracks.messages.invalid_integer_error"))
		}
	}

	var buffer bytes.Buffer
	DJ.Queue.Traverse(func(i int, track interfaces.Track) {
		if i < numTracksToList {
			buffer.WriteString(fmt.Sprintf(viper.GetString("commands.listtracks.messages.track_listing"),
				i+1, track.GetTitle(), track.GetSubmitter()))
		}
	})

	return buffer.String(), true, nil
}
