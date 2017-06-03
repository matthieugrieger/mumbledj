/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/forceskip.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// ForceSkipCommand is a command that immediately skips the current track.
type ForceSkipCommand struct{}

// Aliases returns the current aliases for the command.
func (c *ForceSkipCommand) Aliases() []string {
	return viper.GetStringSlice("commands.forceskip.aliases")
}

// Description returns the description for the command.
func (c *ForceSkipCommand) Description() string {
	return viper.GetString("commands.forceskip.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *ForceSkipCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.forceskip.is_admin")
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
func (c *ForceSkipCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if DJ.Queue.Length() == 0 {
		return "", true, errors.New(viper.GetString("commands.common_messages.no_tracks_error"))
	}

	DJ.Queue.StopCurrent()

	return fmt.Sprintf(viper.GetString("commands.forceskip.messages.track_skipped"),
		user.Name), false, nil
}
