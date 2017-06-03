/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/forceskipplaylist.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"

	"github.com/layeh/gumble/gumble"
	"github.com/RichardNysater/mumbledj/interfaces"
	"github.com/spf13/viper"
)

// ForceSkipPlaylistCommand is a command that immediately skips the current
// playlist.
type ForceSkipPlaylistCommand struct{}

// Aliases returns the current aliases for the command.
func (c *ForceSkipPlaylistCommand) Aliases() []string {
	return viper.GetStringSlice("commands.forceskipplaylist.aliases")
}

// Description returns the description for the command.
func (c *ForceSkipPlaylistCommand) Description() string {
	return viper.GetString("commands.forceskipplaylist.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *ForceSkipPlaylistCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.forceskipplaylist.is_admin")
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
func (c *ForceSkipPlaylistCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	var (
		currentTrack interfaces.Track
		err          error
	)

	if currentTrack, err = DJ.Queue.CurrentTrack(); err != nil {
		return "", true, errors.New(viper.GetString("commands.common_messages.no_tracks_error"))
	}

	if playlist := currentTrack.GetPlaylist(); playlist == nil {
		return "", true, errors.New(viper.GetString("commands.forceskipplaylist.messages.no_playlist_error"))
	}

	DJ.Queue.SkipPlaylist()

	return fmt.Sprintf(viper.GetString("commands.forceskipplaylist.messages.playlist_skipped"),
		user.Name), false, nil
}
