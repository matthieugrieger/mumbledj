/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/skipplaylist.go
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

// SkipPlaylistCommand is a command that places a vote to skip the current
// playlist.
type SkipPlaylistCommand struct{}

// Aliases returns the current aliases for the command.
func (c *SkipPlaylistCommand) Aliases() []string {
	return viper.GetStringSlice("commands.skipplaylist.aliases")
}

// Description returns the description for the command.
func (c *SkipPlaylistCommand) Description() string {
	return viper.GetString("commands.skipplaylist.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *SkipPlaylistCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.skipplaylist.is_admin")
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
func (c *SkipPlaylistCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	var (
		currentTrack interfaces.Track
		err          error
	)

	if currentTrack, err = DJ.Queue.CurrentTrack(); err != nil {
		return "", true, errors.New(viper.GetString("commands.common_messages.no_tracks_error"))
	}

	if playlist := currentTrack.GetPlaylist(); playlist == nil {
		return "", true, errors.New(viper.GetString("commands.skipplaylist.messages.no_playlist_error"))
	}
	if currentTrack.GetPlaylist().GetSubmitter() == user.Name {
		DJ.Queue.SkipPlaylist()
		return fmt.Sprintf(viper.GetString("commands.skipplaylist.messages.submitter_voted"), user.Name), false, nil
	}
	if err := DJ.Skips.AddPlaylistSkip(user); err != nil {
		return "", true, errors.New(viper.GetString("commands.skipplaylist.messages.already_voted_error"))
	}

	return fmt.Sprintf(viper.GetString("commands.skipplaylist.messages.voted"), user.Name), false, nil
}
