/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/register.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// RegisterCommand is a command that registers the bot on the server.
type RegisterCommand struct{}

// Aliases returns the current aliases for the command.
func (c *RegisterCommand) Aliases() []string {
	return viper.GetStringSlice("commands.register.aliases")
}

// Description returns the description for the command.
func (c *RegisterCommand) Description() string {
	return viper.GetString("commands.register.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *RegisterCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.register.is_admin")
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
func (c *RegisterCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if DJ.Client.Self.IsRegistered() {
		return "", true, errors.New(viper.GetString("commands.register.messages.already_registered_error"))
	}

	DJ.Client.Self.Register()

	return viper.GetString("commands.register.messages.registered"), true, nil
}
