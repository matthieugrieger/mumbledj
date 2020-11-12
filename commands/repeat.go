/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/currenttrack.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 * Copyright (c) 2019 Reikion (MIT License)
 */

package commands

import (
	"github.com/spf13/viper"
	"layeh.com/gumble/gumble"
)

// RepeatCommand is a command that safely kills the bot.
type RepeatCommand struct{}

// Aliases returns the current aliases for the command.
func (c *RepeatCommand) Aliases() []string {
	return viper.GetStringSlice("commands.repeat.aliases")
}

// Description returns the description for the command.
func (c *RepeatCommand) Description() string {
	return viper.GetString("commands.repeat.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *RepeatCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.repeat.is_admin")
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
func (c *RepeatCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	flag := DJ.Player.RepeatMode()
	if flag {
		return viper.GetString("commands.repeat.messages.enabled"), false, nil
	}
	return viper.GetString("commands.repeat.messages.disabled"), false, nil
}
