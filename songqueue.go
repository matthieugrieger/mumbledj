/*
 * MumbleDJ
 * By Matthieu Grieger
 * songqueue.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

type SongQueue struct {
	queue Queue
}

func NewSongQueue() *SongQueue {
	return &SongQueue{}
}

func (q *SongQueue) AddSong(s *Song) bool {
	return false
}

func (q *SongQueue) NextSong() bool {
	return false
}

//func (q *SongQueue) CurrentSong() Song {
//	return false
//}
