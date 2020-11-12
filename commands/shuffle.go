/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/shuffle.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// ShuffleCommand is a command that shuffles the audio queue.
type ShuffleCommand struct{}

// Aliases returns the current aliases for the command.
func (c *ShuffleCommand) Aliases() []string {
	return viper.GetStringSlice("commands.shuffle.aliases")
}

// Description returns the description for the command.
func (c *ShuffleCommand) Description() string {
	return viper.GetString("commands.shuffle.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *ShuffleCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.shuffle.is_admin")
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
func (c *ShuffleCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	length := DJ.Queue.Length()
	if length == 0 {
		return "", true, errors.New(viper.GetString("commands.common_messages.no_tracks_error"))
	}
	if length <= 2 {
		return "", true, errors.New(viper.GetString("commands.shuffle.messages.not_enough_tracks_error"))
	}

	DJ.Queue.ShuffleTracks()

	return viper.GetString("commands.shuffle.messages.shuffled"), false, nil
}
