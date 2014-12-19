/*
 * MumbleDJ
 * By Matthieu Grieger
 * songqueue.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

type SongQueue struct {
	queue *Queue
}

func NewSongQueue() *SongQueue {
	return &SongQueue{
		queue: NewQueue(),
	}
}

func (q *SongQueue) AddSong(s *Song) bool {
	beforeLen := q.queue.Len()
	q.queue.Push(s)
	if q.queue.Len() == beforeLen+1 {
		return true
	} else {
		return false
	}
	return true
}

func (q *SongQueue) NextSong() *Song {
	return q.queue.Poll().(*Song)
}

func (q *SongQueue) CurrentSong() *Song {
	return q.queue.Peek().(*Song)
}
