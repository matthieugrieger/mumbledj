/*
 * MumbleDJ
 * By Matthieu Grieger
 * interfaces/skiptracker.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package interfaces

import "layeh.com/gumble/gumble"

// SkipTracker is the interface which should be interacted with for skip operations.
// Using the SkipTracker interface ensures thread safety.
type SkipTracker interface {
	AddTrackSkip(*gumble.User) error
	AddPlaylistSkip(*gumble.User) error
	RemoveTrackSkip(*gumble.User) error
	RemovePlaylistSkip(*gumble.User) error
	NumTrackSkips() int
	NumPlaylistSkips() int
	ResetTrackSkips()
	ResetPlaylistSkips()
}
