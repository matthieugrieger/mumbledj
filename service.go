/*
 * MumbleDJ
 * By Matthieu Grieger
 * service.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"

	"github.com/layeh/gumble/gumble"
)

// Service interface. Each service should implement these functions
type Service interface {
	ServiceName() string
	URLRegex(string) bool                            // Can service deal with URL
	NewRequest(*gumble.User, string) (string, error) // Create song/playlist and add to the queue
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

var services = []Service{YouTube{}}

func findServiceAndAdd(user *gumble.User, url string) (string, error) {
	var urlService Service

	// Checks all services to see if any can take the URL
	for _, service := range services {
		if service.URLRegex(url) {
			urlService = service
		}
	}

	if urlService == nil {
		return "", errors.New("INVALID_URL")
	} else {
		oldLength := dj.queue.Len()
		var title string
		var err error
		if title, err := urlService.NewRequest(user, url); err == nil {

			// Starts playing the new song if nothing else is playing
			if oldLength == 0 && dj.queue.Len() != 0 && !dj.audioStream.IsPlaying() {
				if err := dj.queue.CurrentSong().Download(); err == nil {
					dj.queue.CurrentSong().Play()
				} else {
					dj.queue.CurrentSong().Delete()
					dj.queue.OnSongFinished()
					return "", errors.New("FAILED_TO_DOWNLOAD")
				}
			}
		}
		return title, err
	}
}
