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
	queue *Queue
}

// Initializes a new queue and returns the new SongQueue.
func NewSongQueue() *SongQueue {
	return &SongQueue{
		queue: NewQueue(),
	}
}

// Adds a song to the SongQueue.
func (q *SongQueue) AddSong(s *Song) error {
	beforeLen := q.queue.Len()
	q.queue.Push(s)
	if q.queue.Len() == beforeLen+1 {
		return nil
	} else {
		return errors.New("Could not add Song to the SongQueue.")
	}
}

// Moves to the next song in SongQueue. NextSong() pops the first value of the queue, and is stored
// in dj.currentSong.
func (q *SongQueue) NextSong() *Song {
	return q.queue.Poll().(*Song)
}

// Returns the length of the SongQueue.
func (q *SongQueue) Len() int {
	return q.queue.Len()
}
