/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/numtracks.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"fmt"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// NumTracksCommand is a command that outputs the current number of tracks
// in the queue.
type NumTracksCommand struct{}

// Aliases returns the current aliases for the command.
func (c *NumTracksCommand) Aliases() []string {
	return viper.GetStringSlice("commands.numtracks.aliases")
}

// Description returns the description for the command.
func (c *NumTracksCommand) Description() string {
	return viper.GetString("commands.numtracks.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *NumTracksCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.numtracks.is_admin")
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
func (c *NumTracksCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	length := DJ.Queue.Length()
	if length == 1 {
		return viper.GetString("commands.numtracks.messages.one_track"), true, nil
	}

	return fmt.Sprintf(viper.GetString("commands.numtracks.messages.plural_tracks"), length), true, nil
}
