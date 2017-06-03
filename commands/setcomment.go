/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/setcomment.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"fmt"
	"strings"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// SetCommentCommand is a command that changes the Mumble comment of the bot.
type SetCommentCommand struct{}

// Aliases returns the current aliases for the command.
func (c *SetCommentCommand) Aliases() []string {
	return viper.GetStringSlice("commands.setcomment.aliases")
}

// Description returns the description for the command.
func (c *SetCommentCommand) Description() string {
	return viper.GetString("commands.setcomment.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *SetCommentCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.setcomment.is_admin")
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
func (c *SetCommentCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if len(args) == 0 {
		DJ.Client.Do(func() {
			DJ.Client.Self.SetComment("")
		})
		return viper.GetString("commands.setcomment.messages.comment_removed"), true, nil
	}

	var newComment string
	for _, arg := range args {
		newComment += arg + " "
	}
	strings.TrimSpace(newComment)

	DJ.Client.Do(func() {
		DJ.Client.Self.SetComment(newComment)
	})

	return fmt.Sprintf(viper.GetString("commands.setcomment.messages.comment_changed"),
		newComment), true, nil
}
