/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/generic_service.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package services

import (
	"errors"
	"regexp"
)

// GenericService is a generic struct that should be embedded
// in other service structs, as it provides useful helper
// methods and properties.
type GenericService struct {
	ReadableName  string
	Format        string
	TrackRegex    []*regexp.Regexp
	PlaylistRegex []*regexp.Regexp
}

// GetReadableName returns the readable name for the service.
func (gs *GenericService) GetReadableName() string {
	return gs.ReadableName
}

// GetFormat returns the youtube-dl format for the service.
func (gs *GenericService) GetFormat() string {
	return gs.Format
}

// CheckURL matches the passed URL with a list of regex patterns
// for valid URLs associated with this service. Returns true if a
// match is found, false otherwise.
func (gs *GenericService) CheckURL(url string) bool {
	if gs.isTrack(url) || gs.isPlaylist(url) {
		return true
	}
	return false
}

func (gs *GenericService) isTrack(url string) bool {
	for _, regex := range gs.TrackRegex {
		if regex.MatchString(url) {
			return true
		}
	}
	return false
}

func (gs *GenericService) isPlaylist(url string) bool {
	for _, regex := range gs.PlaylistRegex {
		if regex.MatchString(url) {
			return true
		}
	}
	return false
}

func (gs *GenericService) getID(url string) (string, error) {
	var allRegex []*regexp.Regexp

	if gs.PlaylistRegex != nil {
		allRegex = append(gs.TrackRegex, gs.PlaylistRegex...)
	} else {
		allRegex = gs.TrackRegex
	}

	for _, regex := range allRegex {
		match := regex.FindStringSubmatch(url)
		if match == nil {
			// Move on to next regex, this one didn't match.
			continue
		}
		result := make(map[string]string)
		for i, name := range regex.SubexpNames() {
			if i < len(match) {
				result[name] = match[i]
			}
		}

		return result["id"], nil
	}
	return "", errors.New("No match found for URL")
}
