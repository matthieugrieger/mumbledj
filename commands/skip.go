/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/skip.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// SkipCommand is a command that places a vote to skip the current track.
type SkipCommand struct{}

// Aliases returns the current aliases for the command.
func (c *SkipCommand) Aliases() []string {
	return viper.GetStringSlice("commands.skip.aliases")
}

// Description returns the description for the command.
func (c *SkipCommand) Description() string {
	return viper.GetString("commands.skip.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *SkipCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.skip.is_admin")
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
func (c *SkipCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if DJ.Queue.Length() == 0 {
		return "", true, errors.New(viper.GetString("commands.common_messages.no_tracks_error"))
	}
	if DJ.Queue.GetTrack(0).GetSubmitter() == user.Name {
		// The user who submitted the track is skipping, this means we skip this track immediately.
		DJ.Queue.StopCurrent()
		return fmt.Sprintf(viper.GetString("commands.skip.messages.submitter_voted"), user.Name), false, nil
	}
	if err := DJ.Skips.AddTrackSkip(user); err != nil {
		return "", true, errors.New(viper.GetString("commands.skip.messages.already_voted_error"))
	}

	return fmt.Sprintf(viper.GetString("commands.skip.messages.voted"), user.Name), false, nil
}
