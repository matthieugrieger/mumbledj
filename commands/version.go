/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/version.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"fmt"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// VersionCommand is a command that outputs the local MumbleDJ version.
type VersionCommand struct{}

// Aliases returns the current aliases for the command.
func (c *VersionCommand) Aliases() []string {
	return viper.GetStringSlice("commands.version.aliases")
}

// Description returns the description for the command.
func (c *VersionCommand) Description() string {
	return viper.GetString("commands.version.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *VersionCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.version.is_admin")
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
func (c *VersionCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	return fmt.Sprintf(viper.GetString("commands.version.messages.version"), DJ.Version), true, nil
}
