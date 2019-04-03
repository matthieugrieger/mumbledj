/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/pause.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// PauseCommand is a command that pauses audio playback.
type PauseCommand struct{}

// Aliases returns the current aliases for the command.
func (c *PauseCommand) Aliases() []string {
	return viper.GetStringSlice("commands.pause.aliases")
}

// Description returns the description for the command.
func (c *PauseCommand) Description() string {
	return viper.GetString("commands.pause.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *PauseCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.pause.is_admin")
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
func (c *PauseCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	err := DJ.Queue.PauseCurrent()
	if err != nil {
		return "", true, errors.New(viper.GetString("commands.pause.messages.no_audio_error"))
	}
	return fmt.Sprintf(viper.GetString("commands.pause.messages.paused"), user.Name), false, nil
}
