/*
 * MumbleDJ
 * By Matthieu Grieger
 * songqueue.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"
)

// SongQueue type declaration. Serves as a wrapper around the queue structure defined in queue.go.
type SongQueue struct {
	queue []*Song
}

// Initializes a new queue and returns the new SongQueue.
func NewSongQueue() *SongQueue {
	return &SongQueue{
		queue: make([]*Song, 0),
	}
}

// Adds a song to the SongQueue.
func (q *SongQueue) AddSong(s *Song) error {
	beforeLen := len(q.queue)
	q.queue = append(q.queue, s)
	if len(q.queue) == beforeLen+1 {
		return nil
	} else {
		return errors.New("Could not add Song to the SongQueue.")
	}
}

// Moves to the next song in SongQueue. NextSong() pops the first value of the queue, and is stored
// in dj.currentSong.
func (q *SongQueue) NextSong() *Song {
	s, queue := q.queue[0], q.queue[1:]
	q.queue = queue
	return s
}

// Returns the length of the SongQueue.
func (q *SongQueue) Len() int {
	return len(q.queue)
}
