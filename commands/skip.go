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

	"github.com/layeh/gumble/gumble"
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
		return "", true, errors.New("The queue is currently empty. There is no track to skip")
	}
	if err := DJ.Skips.AddTrackSkip(user); err != nil {
		return "", true, errors.New("You have already voted to skip this track")
	}

	return fmt.Sprintf("<b>%s</b> has voted to skip the current track.", user.Name), false, nil
}
