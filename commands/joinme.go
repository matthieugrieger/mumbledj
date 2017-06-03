/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/joinme.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"

	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleffmpeg"
	"github.com/spf13/viper"
)

// JoinMeCommand is a command that moves the bot to the channel of the user
// who issued the command.
type JoinMeCommand struct{}

// Aliases returns the current aliases for the command.
func (c *JoinMeCommand) Aliases() []string {
	return viper.GetStringSlice("commands.joinme.aliases")
}

// Description returns the description for the command.
func (c *JoinMeCommand) Description() string {
	return viper.GetString("commands.joinme.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *JoinMeCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.joinme.is_admin")
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
func (c *JoinMeCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if DJ.AudioStream != nil && DJ.AudioStream.State() == gumbleffmpeg.StatePlaying &&
		len(DJ.Client.Self.Channel.Users) > 1 {
		return "", true, errors.New(viper.GetString("commands.joinme.messages.others_are_listening_error"))
	}

	DJ.Client.Do(func() {
		DJ.Client.Self.Move(user.Channel)
	})

	return viper.GetString("commands.joinme.messages.in_your_channel"), true, nil
}
