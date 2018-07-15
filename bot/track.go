/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/track.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"time"

	"reik.pl/mumbledj/interfaces"
)

// Track stores all metadata related to an audio track.
type Track struct {
	ID             string
	URL            string
	Title          string
	Author         string
	AuthorURL      string
	Submitter      string
	Service        string
	Filename       string
	ThumbnailURL   string
	Duration       time.Duration
	PlaybackOffset time.Duration
	Playlist       interfaces.Playlist
}

// GetID returns the ID of the track.
func (t Track) GetID() string {
	return t.ID
}

// GetURL returns the URL of the track.
func (t Track) GetURL() string {
	return t.URL
}

// GetTitle returns the title of the track.
func (t Track) GetTitle() string {
	return t.Title
}

// GetAuthor returns the author of the track.
func (t Track) GetAuthor() string {
	return t.Author
}

// GetAuthorURL returns the URL that links to the author of the track.
func (t Track) GetAuthorURL() string {
	return t.AuthorURL
}

// GetSubmitter returns the submitter of the track.
func (t Track) GetSubmitter() string {
	return t.Submitter
}

// GetService returns the name of the service from which the track was retrieved from.
func (t Track) GetService() string {
	return t.Service
}

// GetFilename returns the name of the file stored on disk, if it exists. If no
// file on disk exists an empty string and error are returned.
func (t Track) GetFilename() string {
	return t.Filename
}

// GetThumbnailURL returns the URL to the thumbnail for the track. If no thumbnail
// exists an empty string and error are returned.
func (t Track) GetThumbnailURL() string {
	return t.ThumbnailURL
}

// GetDuration returns the duration of the track.
func (t Track) GetDuration() time.Duration {
	return t.Duration
}

// GetPlaybackOffset returns the playback offset for the track. A duration
// of 0 is given to tracks that do not specify an offset.
func (t Track) GetPlaybackOffset() time.Duration {
	return t.PlaybackOffset
}

// GetPlaylist returns the playlist the track is associated with, if it exists. If
// the track is not associated with a playlist a nil playlist and error are returned.
func (t Track) GetPlaylist() interfaces.Playlist {
	return t.Playlist
}
