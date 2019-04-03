/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/help.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"fmt"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// HelpCommand is a command that outputs a help message that shows the
// available commands and their aliases.
type HelpCommand struct{}

// Aliases returns the current aliases for the command.
func (c *HelpCommand) Aliases() []string {
	return viper.GetStringSlice("commands.help.aliases")
}

// Description returns the description for the command.
func (c *HelpCommand) Description() string {
	return viper.GetString("commands.help.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *HelpCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.help.is_admin")
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
func (c *HelpCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	commandString := "<b>%s</b> -- %s<br>"
	regularCommands := ""
	adminCommands := ""
	totalString := ""

	for _, command := range Commands {
		currentString := fmt.Sprintf(commandString, command.Aliases(), command.Description())
		if command.IsAdminCommand() {
			adminCommands += currentString
		} else {
			regularCommands += currentString
		}
	}

	totalString = viper.GetString("commands.help.messages.commands_header") + regularCommands

	isAdmin := false
	if viper.GetBool("admins.enabled") {
		isAdmin = DJ.IsAdmin(user)
	} else {
		isAdmin = true
	}

	if isAdmin {
		totalString += viper.GetString("commands.help.messages.admin_commands_header") + adminCommands
	}

	return totalString, true, nil
}
