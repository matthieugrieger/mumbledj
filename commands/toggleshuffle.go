/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/toggleshuffle.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// ToggleShuffleCommand is a command that changes the Mumble comment of the bot.
type ToggleShuffleCommand struct{}

// Aliases returns the current aliases for the command.
func (c *ToggleShuffleCommand) Aliases() []string {
	return viper.GetStringSlice("commands.toggleshuffle.aliases")
}

// Description returns the description for the command.
func (c *ToggleShuffleCommand) Description() string {
	return viper.GetString("commands.toggleshuffle.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *ToggleShuffleCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.toggleshuffle.is_admin")
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
func (c *ToggleShuffleCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if viper.GetBool("queue.automatic_shuffle_on") {
		viper.Set("queue.automatic_shuffle_on", false)
		return viper.GetString("commands.toggleshuffle.messages.toggled_off"), false, nil
	}
	viper.Set("queue.automatic_shuffle_on", true)
	return viper.GetString("commands.toggleshuffle.messages.toggled_on"), false, nil
}
