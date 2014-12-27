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

type SongQueue struct {
	queue *Queue
}

func NewSongQueue() *SongQueue {
	return &SongQueue{
		queue: NewQueue(),
	}
}

func (q *SongQueue) AddSong(s *Song) error {
	beforeLen := q.queue.Len()
	q.queue.Push(s)
	if q.queue.Len() == beforeLen+1 {
		return nil
	} else {
		return errors.New("Could not add Song to the SongQueue.")
	}
}

func (q *SongQueue) NextSong() *Song {
	return q.queue.Poll().(*Song)
}

func (q *SongQueue) Len() int {
	return q.queue.Len()
}
