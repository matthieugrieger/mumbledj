/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/queue.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"errors"
	"fmt"

	"math/rand"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.reik.pl/mumbledj/interfaces"

	// needed for loading opus codes needed by gumble
	_ "layeh.com/gumble/opus"
)

// Queue holds the audio tracks queue itself along with useful methods for
// performing actions on the queue. It checks if track conform current
// config values.
type Queue struct {
	queue []interfaces.Track
	mutex sync.RWMutex
	// used for blocking if queue is empty
	notEmptyQueue chan struct{}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// NewQueue initializes a new queue and returns it.
func NewQueue() *Queue {
	q := &Queue{
		queue:         make([]interfaces.Track, 0),
		notEmptyQueue: make(chan struct{}, 0),
	}

	return q
}

// Length returns the length of the queue.
func (q *Queue) Length() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	length := len(q.queue)
	return length
}

// Reset removes all tracks from the queue and reset state of queue.
func (q *Queue) Reset() {
	q.mutex.Lock()
	select {
	case q.notEmptyQueue <- struct{}{}:
		// somebody's waiting for new item in queue, message sent already
	default:
		//do nothing, nobody's waiting for new item in queue
	}
	close(q.notEmptyQueue)

	q.queue = make([]interfaces.Track, 0)
	q.notEmptyQueue = make(chan struct{}, 0)
	q.mutex.Unlock()
	q.mutex = sync.RWMutex{}
}

// AppendTrack adds a track to the back of the queue.
func (q *Queue) AppendTrack(t interfaces.Track) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.appendTrack(t)
}

// AppendTracks adds a tracks to the back of the queue.
func (q *Queue) AppendTracks(ts []interfaces.Track) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	var nErr int
	for _, t := range ts {
		err := q.appendTrack(t)
		if err != nil {
			nErr++
		}
	}

	if nErr == len(ts) {
		return errors.New("Could not add tracks to queue")
	}
	if nErr != 0 && nErr < len(ts) {
		return errors.New("Some tracks could not be added to queue")
	}

	return nil
}

func (q *Queue) appendTrack(t interfaces.Track) error {
	beforeLen := len(q.queue)
	// An error should never occur here since maxTrackDuration is restricted to
	// ints. Any error in the configuration will be caught during yaml load.
	maxTrackDuration, _ := time.ParseDuration(fmt.Sprintf("%ds",
		viper.GetInt("queue.max_track_duration")))

	if viper.GetInt("queue.max_track_duration") == 0 || t.GetDuration() <= maxTrackDuration {
		q.queue = append(q.queue, t)
	} else {
		return errors.New("The track is too long to add to the queue")
	}
	if len(q.queue) != beforeLen+1 {
		return errors.New("Could not add track to queue")
	}
	if beforeLen == 0 {
		select {
		case q.notEmptyQueue <- struct{}{}:
			// somebody's waiting for new item in queue, message sent already
		default:
			//do nothing, nobody's waiting for new item in queue
		}
	}
	return nil
}

// AppendTrack adds a track to the back of the queue.
func (q *Queue) PrependTrack(t interfaces.Track) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.prependTrack(t)
}

func (q *Queue) prependTrack(t interfaces.Track) error {
	beforeLen := len(q.queue)
	// An error should never occur here since maxTrackDuration is restricted to
	// ints. Any error in the configuration will be caught during yaml load.
	maxTrackDuration, _ := time.ParseDuration(fmt.Sprintf("%ds",
		viper.GetInt("queue.max_track_duration")))

	if viper.GetInt("queue.max_track_duration") == 0 || t.GetDuration() <= maxTrackDuration {
		q.queue = append([]interfaces.Track{t}, q.queue...)
	} else {
		return errors.New("The track is too long to add to the queue")
	}
	if len(q.queue) != beforeLen+1 {
		return errors.New("Could not add track to queue")
	}
	if beforeLen == 0 {
		select {
		case q.notEmptyQueue <- struct{}{}:
			// somebody's waiting for new item in queue, message sent already
		default:
			//do nothing, nobody's waiting for new item in queue
		}
	}
	return nil
}

func (q *Queue) InsertTrack(i int, t interfaces.Track) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	beforeLen := len(q.queue)

	// An error should never occur here since maxTrackDuration is restricted to
	// ints. Any error in the configuration will be caught during yaml load.
	maxTrackDuration, _ := time.ParseDuration(fmt.Sprintf("%ds",
		viper.GetInt("queue.max_track_duration")))

	if viper.GetInt("queue.max_track_duration") == 0 ||
		t.GetDuration() <= maxTrackDuration {
		q.queue = append(q.queue, Track{})
		copy(q.queue[i+1:], q.queue[i:])
		q.queue[i] = t
	} else {
		return errors.New("The track is too long to add to the queue")
	}
	if len(q.queue) == beforeLen+1 {
		return nil
	}
	return errors.New("Could not add track to queue")
}

// GetTrack takes an `index` argument to determine which track to return.
// If the track in position `index` exists, it is returned. Otherwise,
// nil is returned. If queue is empty it will block calling goroutine until
// track appear.
func (q *Queue) GetTrack(index int) interfaces.Track {
	q.mutex.RLock()

	if len(q.queue) == 0 {
		q.mutex.RUnlock()
		<-q.notEmptyQueue
		q.mutex.RLock()
	}
	if index >= len(q.queue) {
		q.mutex.RUnlock()
		return nil
	}

	track := q.queue[index]
	q.mutex.RUnlock()
	return track
}

// GetTrackNoWait takes an `index` argument to determine which track to return.
// If the track in position `index` exists, it is returned. Otherwise,
// nil is returned.
func (q *Queue) GetTrackNoWait(index int) interfaces.Track {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	if index >= len(q.queue) {
		return nil
	}
	track := q.queue[index]
	return track
}

// RemoveTrack takes an `index` argument [0:q.Length(queue)) and removes track connected with that
// index from the queue. If the track in position `index` exists, it is returned. Otherwise,
// nil is returned. Note that it may not be the same track as in GetTrack, because other goroutine could
// RemoveTrack earlier.
func (q *Queue) RemoveTrack(index int) interfaces.Track {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if index >= len(q.queue) {
		return nil
	}
	track := q.queue[index]
	q.queue = append(q.queue[:index], q.queue[index+1:]...)
	return track
}

// RemoveTrackIf removes item if given function returns true.
// Function returns number of removed elements
// TODO: Improve implementation to reuse memory instead creating new slice
func (q *Queue) RemoveTrackIf(fun func(int, interfaces.Track) bool) int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	newQueue := []interfaces.Track{}
	var removed int

	for i, track := range q.queue {
		if fun(i, track) {
			removed++
		} else {
			newQueue = append(newQueue, track)
		}
	}
	q.queue = newQueue

	return removed
}

// PeekNextTrack peeks at the next track and returns it.
func (q *Queue) PeekNextTrack() (interfaces.Track, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	if len(q.queue) > 1 {
		if viper.GetBool("queue.automatic_shuffle_on") {
			q.RandomNextTrack(false)
		}
		next := q.queue[1]
		return next, nil
	}
	return nil, errors.New("There is no track coming up next")
}

// Traverse is a traversal function for Queue. Allows a visit function to
// be passed in which performs the specified action on each queue item.
func (q *Queue) Traverse(visit func(i int, t interfaces.Track)) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	if len(q.queue) > 0 {
		for queueIndex, queueTrack := range q.queue {
			visit(queueIndex, queueTrack)
		}
	}
}

// ShuffleTracks shuffles the queue using an inside-out algorithm.
func (q *Queue) ShuffleTracks() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.queue) >= 3 {
		// Skip the first track, as it is likely playing.
		for i := range q.queue[1:] {
			j := rand.Intn(i + 1)
			q.queue[i+1], q.queue[j+1] = q.queue[j+1], q.queue[i+1]
		}
	}
}

// RandomNextTrack sets a random track as the next track to be played.
func (q *Queue) RandomNextTrack(queueWasEmpty bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.queue) > 1 {
		nextTrackIndex := 1
		if queueWasEmpty {
			nextTrackIndex = 0
		}
		swapIndex := nextTrackIndex + rand.Intn(len(q.queue)-1)
		q.queue[nextTrackIndex], q.queue[swapIndex] = q.queue[swapIndex], q.queue[nextTrackIndex]
	}
}

func (q *Queue) notify() {
    // TODO: code deduplication and using sync.Cond
}
