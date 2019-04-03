/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/move.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"
	"strings"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// MoveCommand is a command that moves the bot from one channel to another.
type MoveCommand struct{}

// Aliases returns the current aliases for the command.
func (c *MoveCommand) Aliases() []string {
	return viper.GetStringSlice("commands.move.aliases")
}

// Description returns the description for the command.
func (c *MoveCommand) Description() string {
	return viper.GetString("commands.move.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *MoveCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.move.is_admin")
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
func (c *MoveCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if len(args) == 0 {
		return "", true, errors.New(viper.GetString("commands.move.messages.no_channel_provided_error"))
	}
	channel := ""
	for _, arg := range args {
		channel += arg + " "
	}
	channel = strings.TrimSpace(channel)
	if channels := strings.Split(channel, "/"); DJ.Client.Channels.Find(channels...) != nil {
		DJ.Client.Do(func() {
			DJ.Client.Self.Move(DJ.Client.Channels.Find(channels...))
		})
	} else {
		return "", true, errors.New(viper.GetString("commands.move.messages.channel_doesnt_exist_error"))
	}

	return fmt.Sprintf(viper.GetString("commands.move.messages.move_successful"), channel), true, nil
}
