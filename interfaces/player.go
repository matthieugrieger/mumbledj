/*
 * MumbleDJ
 * By Matthieu Grieger
 * interfaces/player.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 * Copyright (c) 2019 Reikion (MIT License)
 */

package interfaces

import (
	"context"
)

// Player is the interface which should be used to interact with audio player goroutine.
type Player interface {
	Skip()
	SkipPlaylist()
	PlayCurrentForeverLoop(context.Context)
	CurrentTrack() (Track, error)
	HoldOnTrack() error
	PauseCurrent() error
	ResumeCurrent()
	StopCurrent() error
	RepeatMode() bool
}
