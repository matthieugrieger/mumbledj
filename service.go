/*
 * MumbleDJ
 * By Matthieu Grieger
 * service.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"
	"fmt"

	"github.com/layeh/gumble/gumble"
)

// Service interface. Each service will implement these functions
type Service interface {
	URLRegex(string) bool
	NewRequest(*gumble.User, string) (string, error)
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

var services []Service

func FindServiceAndAdd(user *gumble.User, url string) error {
	var urlService Service

	// Checks all services to see if any can take the URL
	for _, service := range services {
		if service.URLRegex(url) {
			urlService = service
		}
	}

	if urlService == nil {
		Verbose("INVALID_URL")
		return errors.New("INVALID_URL")
	} else {
		oldLength := dj.queue.Len()
		var title string
		var err error

		if title, err = urlService.NewRequest(user, url); err == nil {
			dj.client.Self.Channel.Send(fmt.Sprintf(SONG_ADDED_HTML, user.Name, title), false)

			// Starts playing the new song if nothing else is playing
			if oldLength == 0 && dj.queue.Len() != 0 && !dj.audioStream.IsPlaying() {
				if err := dj.queue.CurrentSong().Download(); err == nil {
					dj.queue.CurrentSong().Play()
				} else {
					dj.queue.CurrentSong().Delete()
					dj.queue.OnSongFinished()
					return errors.New("FAILED_TO_DOWNLOAD")
				}
			}
		} else {
			dj.SendPrivateMessage(user, err.Error())
		}
		return err
	}
}

// RegexpFromURL loops through an array of patterns to see if it matches the url
func RegexpFromURL(url string, patterns []string) *regexp.Regexp {
	for _, pattern := range patterns {
		if re, err := regexp.Compile(pattern); err == nil {
			if re.MatchString(url) {
				return re
			}
		}
	}
	return nil
}
