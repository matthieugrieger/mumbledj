/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/currenttrack.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"os"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// KillCommand is a command that safely kills the bot.
type KillCommand struct{}

// Aliases returns the current aliases for the command.
func (c *KillCommand) Aliases() []string {
	return viper.GetStringSlice("commands.kill.aliases")
}

// Description returns the description for the command.
func (c *KillCommand) Description() string {
	return viper.GetString("commands.kill.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *KillCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.kill.is_admin")
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
func (c *KillCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if err := DJ.Cache.DeleteAll(); err != nil {
		return "", true, err
	}
	if err := DJ.Client.Disconnect(); err != nil {
		return "", true, err
	}

	os.Exit(0)
	return "", true, nil
}
