/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/player.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 * Copyright (c) 2019 Reikion (MIT License)
 */

package bot

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.reik.pl/mumbledj/interfaces"
	"layeh.com/gumble/gumbleffmpeg"

	// needed for loading opus codes needed by gumble
	_ "layeh.com/gumble/opus"
)

// Player should be special goroutine, which prefetch sound from videos and plays audio stream
type Player struct {
	mu     sync.RWMutex
	resume chan struct{}
	//normalFlag       bool
	skippingPlaylistFlag bool
	stopPlayingFlag      bool
	holdOnTrackFlag      bool
	repeatModeFlag       bool
}

func NewPlayer() *Player {
	return &Player{
		resume: make(chan struct{}),
	}
}

// CurrentTrack returns the current Track.
func (p *Player) CurrentTrack() (interfaces.Track, error) {
	track := DJ.Queue.GetTrackNoWait(0)
	if track == nil {
		return nil, errors.New("There are no tracks currently in the queue")
	}

	return track, nil
}

// Skip performs the necessary actions that take place when a track is skipped
// via a command.
func (p *Player) Skip() {
	// Stop if something is playing
	if DJ.AudioStream != nil {
		DJ.AudioStream.Stop()
	}
}
func (p *Player) skip() {
	// Remove all track skips.
	DJ.Skips.ResetTrackSkips()

	track := DJ.Queue.GetTrackNoWait(0)
	if track == nil {
		return
	}

	// If caching is disabled, delete the track from disk.
	if !viper.GetBool("cache.enabled") {
		DJ.YouTubeDL.Delete(track)
	}

	// If automatic track shuffling is enabled, assign a random track in the queue to be the next track.
	if viper.GetBool("queue.automatic_shuffle_on") {
		DJ.Queue.RandomNextTrack(false)
	}

	// Remove all playlist skips if this is the last track of the playlist still in the queue.
	playlist := track.GetPlaylist()
	// make sure that it's playlist
	if playlist != nil {
		id := playlist.GetID()
		playlistIsFinished := true

		DJ.Queue.Traverse(func(i int, t interfaces.Track) {
			if i != 0 && t.GetPlaylist() != nil {
				if t.GetPlaylist().GetID() == id {
					playlistIsFinished = false
				}
			}
		})
		if playlistIsFinished {
			DJ.Skips.ResetPlaylistSkips()
		}
	}

	// Skip the track.
	track = DJ.Queue.RemoveTrack(0)
	if track == nil {
		// Queue is empty if first element is nil
		// Something like this was originally here. I think resets queue to reduce memory usage.
		DJ.Queue.Reset()
	}

}

// SkipPlaylist performs the necessary actions that take place when a playlist
// is skipped via a command.
func (p *Player) SkipPlaylist() {
	track := DJ.Queue.GetTrackNoWait(0)
	if track == nil {
		return
	}

	playlist := track.GetPlaylist()
	if playlist != nil {
		currentPlaylistID := playlist.GetID()

		DJ.Queue.RemoveTrackIf(func(i int, t interfaces.Track) bool {
			if i == 0 {
				// Current playing track will be skipped with following call
				return false
			}
			if otherTrackPlaylist := t.GetPlaylist(); otherTrackPlaylist != nil {
				return otherTrackPlaylist.GetID() == currentPlaylistID
			}
			return false
		})
		DJ.Player.Skip()
	}
}

// PlayCurrentForeverLoop plays tracks from queue and waits for new tracks if queue is empty
func (p *Player) PlayCurrentForeverLoop(ctx context.Context) {
	for {
		p.playCurrent()
		select {
		case <-ctx.Done():
			// loop cancelled
			return
		default:
			// continue
		}
	}
}

func (p *Player) playCurrent() error {
	// wait for signal until we can continue playing
	if p.holdOnTrackFlag || p.stopPlayingFlag {
		<-p.resume
	}

	if DJ.Queue.Length() == 0 {
		logrus.Info("Queue is empty, waiting for track to appear.")
	}

	// Blocking call if DJ.Queue is empty
	currentTrack := DJ.Queue.GetTrack(0)

	// Download track
	filepath := os.ExpandEnv(viper.GetString("cache.directory") + "/" + currentTrack.GetFilename())
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := DJ.YouTubeDL.Download(currentTrack); err != nil {
			// Youtube-dl couldn't download track, proceed to next track in queue
			message := fmt.Sprintf("<b>Error:</b> Download of %s track failed, skipping...", currentTrack.GetURL())
			DJ.Client.Self.Channel.Send(message, false)
			DJ.Queue.RemoveTrack(0)
			return err
		}
	}

	// someone's using OhohohoPlayer, wait for signal
	if DJ.AudioStream != nil && DJ.AudioStream.State() != gumbleffmpeg.StateStopped {
		<-p.resume
	}

	p.mu.Lock()
	source := gumbleffmpeg.SourceFile(filepath)
	DJ.AudioStream = gumbleffmpeg.New(DJ.Client, source)
	DJ.AudioStream.Offset = currentTrack.GetPlaybackOffset()
	DJ.AudioStream.Volume = DJ.Volume
	DJ.AudioStream.Play()
	p.mu.Unlock()

	if viper.GetString("defaults.player_command") == "avconv" {
		DJ.AudioStream.Command = "avconv"
	}

	if viper.GetBool("queue.announce_new_tracks") && !p.holdOnTrackFlag {
		message :=
			`<table width=500>
			 	<tr>
					<td align="center"><img src="data:image/JPEG;base64,%s" width=150/></td>
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

	// clean flags before wait
	p.mu.Lock()
	p.holdOnTrackFlag = false
	p.stopPlayingFlag = false
	p.mu.Unlock()

	DJ.AudioStream.Wait()

	if p.repeatModeFlag {
		DJ.Queue.AppendTrack(currentTrack)
	}

	if !p.holdOnTrackFlag && !p.stopPlayingFlag && !p.skippingPlaylistFlag {
		p.skip()
	}

	if p.skippingPlaylistFlag {
		p.mu.Lock()
		p.skippingPlaylistFlag = false
		p.mu.Unlock()
	}

	return nil
}

func (p *Player) HoldOnTrack() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if DJ.AudioStream == nil ||
		(DJ.AudioStream != nil && DJ.AudioStream.State() == gumbleffmpeg.StateStopped) {
		return errors.New("There is no track to hold on")
	}

	p.holdOnTrackFlag = true
	DJ.AudioStream.Stop()
	track := DJ.Queue.RemoveTrack(0)
	if track != nil {
		t, ok := track.(Track)
		if ok {
			t.PlaybackOffset = DJ.AudioStream.Elapsed()
			DJ.Queue.PrependTrack(t)
		}
	}

	return nil
}

// PauseCurrent pauses the current audio stream if it exists and is not already paused.
func (p *Player) PauseCurrent() error {
	if DJ.AudioStream == nil {
		return errors.New("There is no track to pause")
	}
	DJ.AudioStream.Pause()
	return nil
}

// ResumeCurrent resumes playback of the current audio stream if it exists and is paused.
// It also notifies PlayCurrentForeverLoop() that
func (p *Player) ResumeCurrent() {
	if DJ.AudioStream != nil && DJ.AudioStream.State() == gumbleffmpeg.StatePaused {
		DJ.AudioStream.Play()
		return
	}
	select {
	case p.resume <- struct{}{}:
		// resume if PlayCurrentForeverLoop waits for it
	default:
	}

}

// StopCurrent stops the playback of the current audio stream if it exists.
func (p *Player) StopCurrent() error {
	if DJ.AudioStream == nil {
		return errors.New("The audio stream is nil")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stopPlayingFlag = true
	DJ.AudioStream.Stop()
	return nil
}
