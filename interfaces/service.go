/*
 * MumbleDJ
 * By Matthieu Grieger
 * interfaces/service.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package interfaces

import "layeh.com/gumble/gumble"

// Service is an interface of methods to be implemented
// by various service types, such as YouTube or SoundCloud.
type Service interface {
	GetReadableName() string
	GetFormat() string
	CheckAPIKey() error
	CheckURL(string) bool
	GetTracks(string, *gumble.User) ([]Track, error)
}
