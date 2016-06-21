/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/reload.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"github.com/layeh/gumble/gumble"
	"github.com/matthieugrieger/mumbledj/bot"
	"github.com/spf13/viper"
)

// ReloadCommand is a command that reloads the configuration values for the bot
// from a config file.
type ReloadCommand struct{}

// Aliases returns the current aliases for the command.
func (c *ReloadCommand) Aliases() []string {
	return viper.GetStringSlice("commands.reload.aliases")
}

// Description returns the description for the command.
func (c *ReloadCommand) Description() string {
	return viper.GetString("commands.reload.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *ReloadCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.reload.is_admin")
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
func (c *ReloadCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if err := bot.ReadConfigFile(); err != nil {
		return "", true, err
	}

	return viper.GetString("commands.reload.messages.reloaded"),
		true, nil
}
