/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/skiptracker.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"fmt"
	"sync"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// SkipTracker keeps track of the list of users who have skipped the current
// track or playlist.
type SkipTracker struct {
	TrackSkips    []*gumble.User
	PlaylistSkips []*gumble.User
	trackMutex    sync.RWMutex
	playlistMutex sync.RWMutex
}

// NewSkipTracker returns an empty SkipTracker.
func NewSkipTracker() *SkipTracker {
	return &SkipTracker{
		TrackSkips:    make([]*gumble.User, 0),
		PlaylistSkips: make([]*gumble.User, 0),
	}
}

// AddTrackSkip adds a skip to the SkipTracker for the current track.
func (s *SkipTracker) AddTrackSkip(skipper *gumble.User) error {
	s.trackMutex.Lock()
	for _, user := range s.TrackSkips {
		if user.Name == skipper.Name {
			s.trackMutex.Unlock()
			return fmt.Errorf("%s has already voted to skip the track", skipper.Name)
		}
	}
	s.TrackSkips = append(s.TrackSkips, skipper)
	s.trackMutex.Unlock()
	s.evaluateTrackSkips()
	return nil
}

// AddPlaylistSkip adds a skip to the SkipTracker for the current playlist.
func (s *SkipTracker) AddPlaylistSkip(skipper *gumble.User) error {
	s.playlistMutex.Lock()
	for _, user := range s.PlaylistSkips {
		if user.Name == skipper.Name {
			s.playlistMutex.Unlock()
			return fmt.Errorf("%s has already voted to skip the playlist", skipper.Name)
		}
	}
	s.PlaylistSkips = append(s.PlaylistSkips, skipper)
	s.playlistMutex.Unlock()
	s.evaluatePlaylistSkips()
	return nil
}

// RemoveTrackSkip removes a skip from the SkipTracker for the current track.
func (s *SkipTracker) RemoveTrackSkip(skipper *gumble.User) error {
	s.trackMutex.Lock()
	for i, user := range s.TrackSkips {
		if user.Name == skipper.Name {
			s.TrackSkips = append(s.TrackSkips[:i], s.TrackSkips[i+1:]...)
			s.trackMutex.Unlock()
			return nil
		}
	}
	s.trackMutex.Unlock()
	return fmt.Errorf("%s did not previously vote to skip the track", skipper.Name)
}

// RemovePlaylistSkip removes a skip from the SkipTracker for the current playlist.
func (s *SkipTracker) RemovePlaylistSkip(skipper *gumble.User) error {
	s.playlistMutex.Lock()
	for i, user := range s.PlaylistSkips {
		if user.Name == skipper.Name {
			s.PlaylistSkips = append(s.PlaylistSkips[:i], s.PlaylistSkips[i+1:]...)
			s.playlistMutex.Unlock()
			return nil
		}
	}
	s.playlistMutex.Unlock()
	return fmt.Errorf("%s did not previously vote to skip the playlist", skipper.Name)
}

// NumTrackSkips returns the number of users who have skipped the current track.
func (s *SkipTracker) NumTrackSkips() int {
	s.trackMutex.RLock()
	length := len(s.TrackSkips)
	s.trackMutex.RUnlock()
	return length
}

// NumPlaylistSkips returns the number of users who have skipped the current playlist.
func (s *SkipTracker) NumPlaylistSkips() int {
	s.playlistMutex.RLock()
	length := len(s.PlaylistSkips)
	s.playlistMutex.RUnlock()
	return length
}

// ResetTrackSkips resets the skip slice for the current track.
func (s *SkipTracker) ResetTrackSkips() {
	s.trackMutex.Lock()
	s.TrackSkips = s.TrackSkips[:0]
	s.trackMutex.Unlock()
}

// ResetPlaylistSkips resets the skip slice for the current playlist.
func (s *SkipTracker) ResetPlaylistSkips() {
	s.playlistMutex.Lock()
	s.PlaylistSkips = s.PlaylistSkips[:0]
	s.playlistMutex.Unlock()
}

func (s *SkipTracker) evaluateTrackSkips() {
	s.trackMutex.RLock()
	skipRatio := viper.GetFloat64("queue.track_skip_ratio")
	DJ.Client.Do(func() {
		if float64(len(s.TrackSkips))/float64(len(DJ.Client.Self.Channel.Users)) >= skipRatio {
			// Stopping an audio stream triggers a skip.
			DJ.Queue.StopCurrent()
		}
	})
	s.trackMutex.RUnlock()
}

func (s *SkipTracker) evaluatePlaylistSkips() {
	s.playlistMutex.RLock()
	skipRatio := viper.GetFloat64("queue.playlist_skip_ratio")
	DJ.Client.Do(func() {
		if float64(len(s.PlaylistSkips))/float64(len(DJ.Client.Self.Channel.Users)) >= skipRatio {
			DJ.Queue.SkipPlaylist()
		}
	})
	s.playlistMutex.RUnlock()
}
