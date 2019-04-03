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
	"net/url"

	"math/rand"
	"os"
	"sync"
	"time"

	"layeh.com/gumble/gumbleffmpeg"
	_ "layeh.com/gumble/opus"
	"github.com/spf13/viper"
	"reik.pl/mumbledj/interfaces"
)

// Queue holds the audio queue itself along with useful methods for
// performing actions on the queue.
type Queue struct {
	Queue []interfaces.Track
	mutex sync.RWMutex
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// NewQueue initializes a new queue and returns it.
func NewQueue() *Queue {
	return &Queue{
		Queue: make([]interfaces.Track, 0),
	}
}

// Length returns the length of the queue.
func (q *Queue) Length() int {
	q.mutex.RLock()
	length := len(q.Queue)
	q.mutex.RUnlock()
	return length
}

// Reset removes all tracks from the queue.
func (q *Queue) Reset() {
	q.mutex.Lock()
	q.Queue = q.Queue[:0]
	q.mutex.Unlock()
}

// AppendTrack adds a track to the back of the queue.
func (q *Queue) AppendTrack(t interfaces.Track) error {
	q.mutex.Lock()
	beforeLen := len(q.Queue)

	// An error should never occur here since maxTrackDuration is restricted to
	// ints. Any error in the configuration will be caught during yaml load.
	maxTrackDuration, _ := time.ParseDuration(fmt.Sprintf("%ds",
		viper.GetInt("queue.max_track_duration")))

	if viper.GetInt("queue.max_track_duration") == 0 ||
		t.GetDuration() <= maxTrackDuration {
		q.Queue = append(q.Queue, t)
	} else {
		q.mutex.Unlock()
		return errors.New("The track is too long to add to the queue")
	}
	if len(q.Queue) == beforeLen+1 {
		q.mutex.Unlock()
		q.playIfNeeded()
		return nil
	}
	q.mutex.Unlock()
	return errors.New("Could not add track to queue")
}

// InsertTrack inserts track `t` at position `i` in the queue.
func (q *Queue) InsertTrack(i int, t interfaces.Track) error {
	q.mutex.Lock()
	beforeLen := len(q.Queue)

	// An error should never occur here since maxTrackDuration is restricted to
	// ints. Any error in the configuration will be caught during yaml load.
	maxTrackDuration, _ := time.ParseDuration(fmt.Sprintf("%ds",
		viper.GetInt("queue.max_track_duration")))

	if viper.GetInt("queue.max_track_duration") == 0 ||
		t.GetDuration() <= maxTrackDuration {
		q.Queue = append(q.Queue, Track{})
		copy(q.Queue[i+1:], q.Queue[i:])
		q.Queue[i] = t
	} else {
		q.mutex.Unlock()
		return errors.New("The track is too long to add to the queue")
	}
	if len(q.Queue) == beforeLen+1 {
		q.mutex.Unlock()
		q.playIfNeeded()
		return nil
	}
	q.mutex.Unlock()
	return errors.New("Could not add track to queue")
}

// CurrentTrack returns the current Track.
func (q *Queue) CurrentTrack() (interfaces.Track, error) {
	q.mutex.RLock()
	if len(q.Queue) != 0 {
		current := q.Queue[0]
		q.mutex.RUnlock()
		return current, nil
	}
	q.mutex.RUnlock()
	return nil, errors.New("There are no tracks currently in the queue")
}

// GetTrack takes an `index` argument to determine which track to return.
// If the track in position `index` exists, it is returned. Otherwise,
// nil is returned.
func (q *Queue) GetTrack(index int) interfaces.Track {
	q.mutex.RLock()
	if index >= len(q.Queue) {
		q.mutex.RUnlock()
		return nil
	}
	track := q.Queue[index]
	q.mutex.RUnlock()
	return track
}

// PeekNextTrack peeks at the next track and returns it.
func (q *Queue) PeekNextTrack() (interfaces.Track, error) {
	q.mutex.RLock()
	if len(q.Queue) > 1 {
		if viper.GetBool("queue.automatic_shuffle_on") {
			q.RandomNextTrack(false)
		}
		next := q.Queue[1]
		q.mutex.RUnlock()
		return next, nil
	}
	q.mutex.RUnlock()
	return nil, errors.New("There is no track coming up next")
}

// Traverse is a traversal function for Queue. Allows a visit function to
// be passed in which performs the specified action on each queue item.
func (q *Queue) Traverse(visit func(i int, t interfaces.Track)) {
	q.mutex.RLock()
	if len(q.Queue) > 0 {
		for queueIndex, queueTrack := range q.Queue {
			visit(queueIndex, queueTrack)
		}
	}
	q.mutex.RUnlock()
}

// ShuffleTracks shuffles the queue using an inside-out algorithm.
func (q *Queue) ShuffleTracks() {
	q.mutex.Lock()
	// Skip the first track, as it is likely playing.
	for i := range q.Queue[1:] {
		j := rand.Intn(i + 1)
		q.Queue[i+1], q.Queue[j+1] = q.Queue[j+1], q.Queue[i+1]
	}
	q.mutex.Unlock()
}

// RandomNextTrack sets a random track as the next track to be played.
func (q *Queue) RandomNextTrack(queueWasEmpty bool) {
	q.mutex.Lock()
	if len(q.Queue) > 1 {
		nextTrackIndex := 1
		if queueWasEmpty {
			nextTrackIndex = 0
		}
		swapIndex := nextTrackIndex + rand.Intn(len(q.Queue)-1)
		q.Queue[nextTrackIndex], q.Queue[swapIndex] = q.Queue[swapIndex], q.Queue[nextTrackIndex]
	}
	q.mutex.Unlock()
}

// Skip performs the necessary actions that take place when a track is skipped
// via a command.
func (q *Queue) Skip() {
	// Set AudioStream to nil if it isn't already.
	if DJ.AudioStream != nil {
		DJ.AudioStream = nil
	}

	// Remove all track skips.
	DJ.Skips.ResetTrackSkips()

	q.mutex.Lock()
	// If caching is disabled, delete the track from disk.
	if len(q.Queue) != 0 && !viper.GetBool("cache.enabled") {
		DJ.YouTubeDL.Delete(q.Queue[0])
	}

	// If automatic track shuffling is enabled, assign a random track in the queue to be the next track.
	if viper.GetBool("queue.automatic_shuffle_on") {
		q.mutex.Unlock()
		q.RandomNextTrack(false)
		q.mutex.Lock()
	}

	// Remove all playlist skips if this is the last track of the playlist still in the queue.
	if len(q.Queue) > 0 {
		playlist := q.Queue[0].GetPlaylist()
		// woops, it is empty, better return than panic
		if playlist == nil {
			return
		}
		id := playlist.GetID()
		playlistIsFinished := true

		q.mutex.Unlock()
		q.Traverse(func(i int, t interfaces.Track) {
			if i != 0 && t.GetPlaylist() != nil {
				if t.GetPlaylist().GetID() == id {
					playlistIsFinished = false
				}
			}
		})
		q.mutex.Lock()

		if playlistIsFinished {
			DJ.Skips.ResetPlaylistSkips()
		}
	}

	// Skip the track.
	length := len(q.Queue)
	if length > 1 {
		q.Queue = q.Queue[1:]
	} else {
		q.Queue = make([]interfaces.Track, 0)
	}
	q.mutex.Unlock()

	if err := q.playIfNeeded(); err != nil {
		q.Skip()
	}
}

// SkipPlaylist performs the necessary actions that take place when a playlist
// is skipped via a command.
func (q *Queue) SkipPlaylist() {
	q.mutex.Lock()
	if playlist := q.Queue[0].GetPlaylist(); playlist != nil {
		currentPlaylistID := playlist.GetID()

		// We must loop backwards to prevent missing any elements after deletion.
		// NOTE: We do not remove the first track of the playlist quite yet as that
		// is removed properly with the following Skip() call.
		for i := len(q.Queue) - 1; i >= 1; i-- {
			if otherTrackPlaylist := q.Queue[i].GetPlaylist(); otherTrackPlaylist != nil {
				if otherTrackPlaylist.GetID() == currentPlaylistID {
					q.Queue = append(q.Queue[:i], q.Queue[i+1:]...)
				}
			}
		}
	}
	q.mutex.Unlock()
	q.StopCurrent()
}

// PlayCurrent creates a new audio stream and begins playing the current track.
func (q *Queue) PlayCurrent() error {
	currentTrack := q.GetTrack(0)
	filepath := os.ExpandEnv(viper.GetString("cache.directory") + "/" + currentTrack.GetFilename())
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := DJ.YouTubeDL.Download(q.GetTrack(0)); err != nil {
			return err
		}
	}
	source := gumbleffmpeg.SourceFile(filepath)
	DJ.AudioStream = gumbleffmpeg.New(DJ.Client, source)
	DJ.AudioStream.Offset = currentTrack.GetPlaybackOffset()
	DJ.AudioStream.Volume = DJ.Volume

	if viper.GetString("defaults.player_command") == "avconv" {
		DJ.AudioStream.Command = "avconv"
	}

	if viper.GetBool("queue.announce_new_tracks") {
		message :=
			`<table>
			 	<tr>
					<td align="center"><img src="data:image/JPEG;base64,%s" width=150 /></td>
				</tr>
				<tr>
					<td align="center"><b><a href="%s">%s</a> (%s)</b></td>
				</tr>
				<tr>
					<td align="center">Added by %s</td>
				</tr>
			`
		message = fmt.Sprintf(message, url.QueryEscape(currentTrack.GetThumbnailBase64()), currentTrack.GetURL(),
			currentTrack.GetTitle(), currentTrack.GetDuration().String(), currentTrack.GetSubmitter())
		if currentTrack.GetPlaylist() != nil {
			message = message + fmt.Sprintf(`<tr><td align="center">From playlist "%s"</td></tr>`, currentTrack.GetPlaylist().GetTitle())
		}
		message += `</table>`
		DJ.Client.Self.Channel.Send(message, false)
	}

	DJ.AudioStream.Play()
	go func() {
		DJ.AudioStream.Wait()
		q.Skip()
	}()

	return nil
}

// PauseCurrent pauses the current audio stream if it exists and is not already paused.
func (q *Queue) PauseCurrent() error {
	if DJ.AudioStream == nil {
		return errors.New("There is no track to pause")
	}
	if DJ.AudioStream.State() == gumbleffmpeg.StatePaused {
		return errors.New("The track is already paused")
	}
	DJ.AudioStream.Pause()
	return nil
}

// ResumeCurrent resumes playback of the current audio stream if it exists and is paused.
func (q *Queue) ResumeCurrent() error {
	if DJ.AudioStream == nil {
		return errors.New("There is no track to resume")
	}
	if DJ.AudioStream.State() == gumbleffmpeg.StatePlaying {
		return errors.New("The track is already playing")
	}
	DJ.AudioStream.Play()
	return nil
}

// StopCurrent stops the playback of the current audio stream if it exists.
func (q *Queue) StopCurrent() error {
	if DJ.AudioStream == nil {
		return errors.New("The audio stream is nil")
	}
	DJ.AudioStream.Stop()
	return nil
}

func (q *Queue) playIfNeeded() error {
	if DJ.AudioStream == nil && q.Length() > 0 {
		if err := DJ.YouTubeDL.Download(q.GetTrack(0)); err != nil {
			return err
		}
		if err := q.PlayCurrent(); err != nil {
			return err
		}
	}
	return nil
}
