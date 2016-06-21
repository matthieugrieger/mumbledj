/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/pkg_init.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"github.com/matthieugrieger/mumbledj/bot"
	"github.com/matthieugrieger/mumbledj/interfaces"
)

// DJ is an injected MumbleDJ struct.
var DJ *bot.MumbleDJ

// Commands is a slice of all enabled commands.
var Commands []interfaces.Command

func init() {
	Commands = []interfaces.Command{
		new(AddCommand),
		new(AddNextCommand),
		new(CacheSizeCommand),
		new(CurrentTrackCommand),
		new(ForceSkipCommand),
		new(ForceSkipPlaylistCommand),
		new(HelpCommand),
		new(JoinMeCommand),
		new(KillCommand),
		new(ListTracksCommand),
		new(MoveCommand),
		new(NextTrackCommand),
		new(NumCachedCommand),
		new(NumTracksCommand),
		new(PauseCommand),
		new(ReloadCommand),
		new(ResetCommand),
		new(ResumeCommand),
		new(SetCommentCommand),
		new(ShuffleCommand),
		new(SkipCommand),
		new(SkipPlaylistCommand),
		new(ToggleShuffleCommand),
		new(VersionCommand),
		new(VolumeCommand),
	}
}
