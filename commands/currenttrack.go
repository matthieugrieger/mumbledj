/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/currenttrack.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"

	"layeh.com/gumble/gumble"
	"github.com/RichardNysater/mumbledj/interfaces"
	"github.com/spf13/viper"
)

// CurrentTrackCommand is a command that outputs information related to
// the track that is currently playing (if one exists).
type CurrentTrackCommand struct{}

// Aliases returns the current aliases for the command.
func (c *CurrentTrackCommand) Aliases() []string {
	return viper.GetStringSlice("commands.currenttrack.aliases")
}

// Description returns the description for the command.
func (c *CurrentTrackCommand) Description() string {
	return viper.GetString("commands.currenttrack.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *CurrentTrackCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.currenttrack.is_admin")
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
func (c *CurrentTrackCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	var (
		currentTrack interfaces.Track
		err          error
	)

	if currentTrack, err = DJ.Queue.CurrentTrack(); err != nil {
		return "", true, errors.New(viper.GetString("commands.common_messages.no_tracks_error"))
	}

	return fmt.Sprintf(viper.GetString("commands.currenttrack.messages.current_track"),
		currentTrack.GetTitle(), currentTrack.GetSubmitter()), true, nil
}
