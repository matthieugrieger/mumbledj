/*
 * MumbleDJ
 * By Matthieu Grieger
 * service.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

// Song interface. Each service will implement these
// functions in their Song types.
type Song interface {
	Download() error
	Play()
	Delete() error
	AddSkip() error
	RemoveSkip() error
	SkipReached() bool
	Submitter() string
	Title() string
	ID() string
	Filename() string
	Duration() string
	Thumbnail() string
	Playlist() *Playlist
	DontSkip() bool
}

// Playlist interface. Each service will implement these
// functions in their Playlist types.
type Playlist interface {
	AddSkip() error
	RemoveSkip() error
	DeleteSkippers()
	SkipReached() bool
	ID() string
	Title() string
}
