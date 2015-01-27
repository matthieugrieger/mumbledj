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

// Peeks at the next Song and returns it.
func (q *SongQueue) PeekNext() (*Song, error) {
	if q.Len() != 0 {
		if q.CurrentItem().ItemType() == "playlist" {
			return q.CurrentItem().(*Playlist).songs.queue[1].(*Song), nil
		} else if q.Len() > 1 {
			if q.queue[1].ItemType() == "playlist" {
				return q.queue[1].(*Playlist).songs.queue[0].(*Song), nil
			} else {
				return q.queue[1].(*Song), nil
			}
		} else {
			return nil, errors.New("There is no song coming up next.")
		}
	} else {
		return nil, errors.New("There are no items in the queue.")
	}
}

// Returns the length of the SongQueue.
func (q *SongQueue) Len() int {
	return len(q.queue)
}

// A traversal function for SongQueue. Allows a visit function to be passed in which performs
// the specified action on each queue item. Traverses all individual songs, and all songs
// within playlists.
func (q *SongQueue) Traverse(visit func(i int, item QueueItem)) {
	for iQueue, queueItem := range q.queue {
		if queueItem.ItemType() == "playlist" {
			for iPlaylist, playlistItem := range q.queue[iQueue].(*Playlist).songs.queue {
				visit(iPlaylist, playlistItem)
			}
		} else {
			visit(iQueue, queueItem)
		}
	}
}

// OnItemFinished event. Deletes item that just finished playing, then queues the next item.
func (q *SongQueue) OnItemFinished() {
	if q.Len() != 0 {
		if q.CurrentItem().ItemType() == "playlist" {
			if err := q.CurrentItem().(*Playlist).songs.CurrentItem().(*Song).Delete(); err == nil {
				if q.CurrentItem().(*Playlist).skipped == true {
					if q.Len() > 1 {
						q.NextItem()
						q.PrepareAndPlayNextItem()
					} else {
						q.queue = q.queue[:0]
					}
				} else if q.CurrentItem().(*Playlist).songs.Len() > 1 {
					q.CurrentItem().(*Playlist).songs.NextItem()
					q.PrepareAndPlayNextItem()
				} else {
					if q.Len() > 1 {
						q.NextItem()
						q.PrepareAndPlayNextItem()
					} else {
						q.queue = q.queue[:0]
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
				} else {
					q.queue = q.queue[:0]
				}
			} else {
				panic(err)
			}
		}
	}
}

func (q *SongQueue) PrepareAndPlayNextItem() {
	if q.Len() != 0 {
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
}
