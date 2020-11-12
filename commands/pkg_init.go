/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/pkg_init.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 * Copyright (c) 2019 Reikion (MIT License)
 */

package commands

import (
	"go.reik.pl/mumbledj/assets"
	"go.reik.pl/mumbledj/bot"
	"go.reik.pl/mumbledj/interfaces"
)

// DJ is an injected MumbleDJ struct.
var DJ *bot.MumbleDJ

// Assets embedded in binary
var Assets = assets.Assets

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
		new(RegisterCommand),
		new(ReloadCommand),
		new(RepeatCommand),
		new(ResetCommand),
		new(ResumeCommand),
		new(SearchCommand),
		new(SetCommentCommand),
		new(ShuffleCommand),
		new(SkipCommand),
		new(SkipPlaylistCommand),
		new(ToggleShuffleCommand),
		new(VersionCommand),
		new(VolumeCommand),
		// mine
		new(OhohohoCommand),
	}
}
