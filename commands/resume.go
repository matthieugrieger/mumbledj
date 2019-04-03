/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/resume.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// ResumeCommand is a command that resumes audio playback.
type ResumeCommand struct{}

// Aliases returns the current aliases for the command.
func (c *ResumeCommand) Aliases() []string {
	return viper.GetStringSlice("commands.resume.aliases")
}

// Description returns the description for the command.
func (c *ResumeCommand) Description() string {
	return viper.GetString("commands.resume.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *ResumeCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.resume.is_admin")
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
func (c *ResumeCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	err := DJ.Queue.ResumeCurrent()
	if err != nil {
		return "", true, errors.New(viper.GetString("commands.resume.messages.audio_error"))
	}
	return fmt.Sprintf(viper.GetString("commands.resume.messages.resumed"), user.Name), false, nil
}
