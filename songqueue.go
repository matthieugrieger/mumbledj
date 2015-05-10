/*
 * MumbleDJ
 * By Matthieu Grieger
 * songqueue.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"
	"fmt"
	"time"
)

// SongQueue type declaration.
type SongQueue struct {
	queue []Song
}

// NewSongQueue initializes a new queue and returns it.
func NewSongQueue() *SongQueue {
	return &SongQueue{
		queue: make([]Song, 0),
	}
}

// AddSong adds a Song to the SongQueue.
func (q *SongQueue) AddSong(s Song) error {
	beforeLen := q.Len()
	q.queue = append(q.queue, s)
	if len(q.queue) == beforeLen+1 {
		return nil
	}
	return errors.New("Could not add Song to the SongQueue.")
}

// CurrentSong returns the current Song.
func (q *SongQueue) CurrentSong() Song {
	return q.queue[0]
}

// NextSong moves to the next Song in SongQueue. NextSong() removes the first Song in the queue.
func (q *SongQueue) NextSong() {
	if q.CurrentSong().Playlist() != nil {
		if s, err := q.PeekNext(); err == nil {
			if s.Playlist() != nil && (q.CurrentSong().Playlist().ID() != s.Playlist().ID()) {
				q.CurrentSong().Playlist().DeleteSkippers()
			}
		} else {
			q.CurrentSong().Playlist().DeleteSkippers()
		}
	}
	q.queue = q.queue[1:]
}

// PeekNext peeks at the next Song and returns it.
func (q *SongQueue) PeekNext() (Song, error) {
	if q.Len() > 1 {
		return q.queue[1], nil
	}
	return nil, errors.New("There isn't a Song coming up next.")
}

// Len returns the length of the SongQueue.
func (q *SongQueue) Len() int {
	return len(q.queue)
}

// Traverse is a traversal function for SongQueue. Allows a visit function to be passed in which performs
// the specified action on each queue item.
func (q *SongQueue) Traverse(visit func(i int, s Song)) {
	for sQueue, queueSong := range q.queue {
		visit(sQueue, queueSong)
	}
}

// OnSongFinished event. Deletes Song that just finished playing, then queues the next Song (if exists).
func (q *SongQueue) OnSongFinished() {
	resetOffset, _ := time.ParseDuration(fmt.Sprintf("%ds", 0))
	dj.audioStream.Offset = resetOffset
	if q.Len() != 0 {
		if dj.queue.CurrentSong().DontSkip() == true {
			dj.queue.CurrentSong().SetDontSkip(false)
			q.PrepareAndPlayNextSong()
		} else {
			q.NextSong()
			if q.Len() != 0 {
				q.PrepareAndPlayNextSong()
			}
		}
	}
}

// PrepareAndPlayNextSong prepares next song and plays it if the download succeeds.
// Otherwise the function will print an error message to the channel and skip to the next song.
func (q *SongQueue) PrepareAndPlayNextSong() {
	if err := q.CurrentSong().Download(); err == nil {
		q.CurrentSong().Play()
	} else {
		dj.client.Self.Channel.Send(AUDIO_FAIL_MSG, false)
		q.OnSongFinished()
	}
}
