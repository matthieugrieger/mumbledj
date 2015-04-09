/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/base.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package services

// Song interface. Each service will implement these
// functions in their Song types.
type Song interface {
	Download()
	Play()
	Delete()
	AddSkip()
	RemoveSkip()
	SkipReached()
}

// Playlist interface. Each service will implement these
// functions in their Playlist types.
type Playlist interface {
	AddSkip()
	RemoveSkip()
	DeleteSkippers()
	SkipReached()
}
