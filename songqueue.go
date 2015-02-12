/*
 * MumbleDJ
 * By Matthieu Grieger
 * songqueue.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"
)

// SongQueue type declaration.
type SongQueue struct {
	queue []*Song
}

// Initializes a new queue and returns the new SongQueue.
func NewSongQueue() *SongQueue {
	return &SongQueue{
		queue: make([]*Song, 0),
	}
}

// Adds a Song to the SongQueue.
func (q *SongQueue) AddSong(s *Song) error {
	beforeLen := q.Len()
	q.queue = append(q.queue, s)
	if len(q.queue) == beforeLen+1 {
		return nil
	} else {
		return errors.New("Could not add Song to the SongQueue.")
	}
}

// Returns the current Song.
func (q *SongQueue) CurrentSong() *Song {
	return q.queue[0]
}

// Moves to the next Song in SongQueue. NextSong() removes the first Song in the queue.
func (q *SongQueue) NextSong() {
	if q.CurrentSong().playlist != nil {
		if s, err := q.PeekNext(); err == nil {
			if q.CurrentSong().playlist.id != s.playlist.id {
				q.CurrentSong().playlist.DeleteSkippers()
			}
		} else {
			q.CurrentSong().playlist.DeleteSkippers()
		}
	}
	q.queue = q.queue[1:]
}

// Peeks at the next Song and returns it.
func (q *SongQueue) PeekNext() (*Song, error) {
	if q.Len() > 1 {
		return q.queue[1], nil
	} else {
		return nil, errors.New("There isn't a Song coming up next.")
	}
}

// Returns the length of the SongQueue.
func (q *SongQueue) Len() int {
	return len(q.queue)
}

// A traversal function for SongQueue. Allows a visit function to be passed in which performs
// the specified action on each queue item.
func (q *SongQueue) Traverse(visit func(i int, s *Song)) {
	for sQueue, queueSong := range q.queue {
		visit(sQueue, queueSong)
	}
}

// OnSongFinished event. Deletes Song that just finished playing, then queues the next Song (if exists).
func (q *SongQueue) OnSongFinished() {
	if q.Len() != 0 {
		if dj.queue.CurrentSong().dontSkip == true {
			dj.queue.CurrentSong().dontSkip = false
			q.PrepareAndPlayNextSong()
		} else {
			q.NextSong()
			if q.Len() != 0 {
				q.PrepareAndPlayNextSong()
			}
		}
	}
}

// Prepares next song and plays it if the download succeeds. Otherwise the function will print an error message
// to the channel and skip to the next song.
func (q *SongQueue) PrepareAndPlayNextSong() {
	if err := q.CurrentSong().Download(); err == nil {
		q.CurrentSong().Play()
	} else {
		dj.client.Self.Channel.Send(AUDIO_FAIL_MSG, false)
		q.OnSongFinished()
	}
}
