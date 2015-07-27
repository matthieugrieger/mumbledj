/*
 * MumbleDJ
 * By Matthieu Grieger
 * service.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

// Service interface. Each service should implement these functions
type Service interface {
	ServiceName() string
	URLRegex(string) bool                  // Can service deal with URL
	NewRequest(*gumble.User, string) error // Create song/playlist and add to the queue
}

// Song interface. Each service will implement these
// functions in their Song types.
type Song interface {
	Download() error
	Play()
	Delete() error
	AddSkip(string) error
	RemoveSkip(string) error
	SkipReached(int) bool
	Submitter() string
	Title() string
	ID() string
	Filename() string
	Duration() string
	Thumbnail() string
	Playlist() Playlist
	DontSkip() bool
	SetDontSkip(bool)
}

// Playlist interface. Each service will implement these
// functions in their Playlist types.
type Playlist interface {
	AddSkip(string) error
	RemoveSkip(string) error
	DeleteSkippers()
	SkipReached(int) bool
	ID() string
	Title() string
}

var services = []Service{
	new(YoutubeService),
}
