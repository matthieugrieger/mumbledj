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

// QueueItem type declaration. QueueItem is an interface that groups together Song and Playlist
// types in a queue.
type QueueItem interface {
	AddSkip(string) error
	RemoveSkip(string) error
	SkipReached(int) bool
	ItemType() string
}

// SongQueue type declaration. Serves as a wrapper around the queue structure defined in queue.go.
type SongQueue struct {
	queue []QueueItem
}

// Initializes a new queue and returns the new SongQueue.
func NewSongQueue() *SongQueue {
	return &SongQueue{
		queue: make([]QueueItem, 0),
	}
}

// Adds an item to the SongQueue.
func (q *SongQueue) AddItem(i QueueItem) error {
	beforeLen := q.Len()
	q.queue = append(q.queue, i)
	if len(q.queue) == beforeLen+1 {
		return nil
	} else {
		return errors.New("Could not add QueueItem to the SongQueue.")
	}
}

// Returns the current QueueItem.
func (q *SongQueue) CurrentItem() QueueItem {
	return q.queue[0]
}

// Moves to the next item in SongQueue. NextItem() removes the first value in the queue.
func (q *SongQueue) NextItem() {
	q.queue = q.queue[1:]
}

// Returns the length of the SongQueue.
func (q *SongQueue) Len() int {
	return len(q.queue)
}

// OnItemFinished event. Deletes item that just finished playing, then queues the next item.
func (q *SongQueue) OnItemFinished() {
	if q.CurrentItem().ItemType() == "playlist" {
		if err := q.CurrentItem().(*Playlist).songs.CurrentItem().(*Song).Delete(); err == nil {
			if q.CurrentItem().(*Playlist).skipped == true {
				if q.Len() > 1 {
					q.NextItem()
					q.PrepareAndPlayNextItem()
				} else {
					q.queue = q.queue[1:]
				}
			} else if q.CurrentItem().(*Playlist).songs.Len() > 1 {
				q.CurrentItem().(*Playlist).songs.NextItem()
				q.PrepareAndPlayNextItem()
			} else {
				if q.Len() > 1 {
					q.NextItem()
					q.PrepareAndPlayNextItem()
				}
			}
		} else {
			panic(err)
		}
	} else {
		if err := q.CurrentItem().(*Song).Delete(); err == nil {
			if q.Len() > 1 {
				q.NextItem()
				q.PrepareAndPlayNextItem()
			}
		} else {
			panic(err)
		}
	}
}

func (q *SongQueue) PrepareAndPlayNextItem() {
	if q.CurrentItem().ItemType() == "playlist" {
		if err := q.CurrentItem().(*Playlist).songs.CurrentItem().(*Song).Download(); err == nil {
			q.CurrentItem().(*Playlist).songs.CurrentItem().(*Song).Play()
		} else {
			username := q.CurrentItem().(*Playlist).submitter
			user := dj.client.Self().Channel().Users().Find(username)
			user.Send(AUDIO_FAIL_MSG)
			q.OnItemFinished()
		}
	} else {
		if err := q.CurrentItem().(*Song).Download(); err == nil {
			q.CurrentItem().(*Song).Play()
		} else {
			username := q.CurrentItem().(*Song).submitter
			user := dj.client.Self().Channel().Users().Find(username)
			user.Send(AUDIO_FAIL_MSG)
			q.OnItemFinished()
		}
	}
}
