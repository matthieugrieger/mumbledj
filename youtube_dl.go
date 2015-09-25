/*
 * MumbleDJ
 * By Matthieu Grieger
 * youtube_dl.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
)

// YouTubeSong implements the Song interface
type YouTubeSong struct {
	id        string
	title     string
	thumbnail string
	submitter *gumble.User
	duration  string
	url       string
	offset    int
	format    string
	playlist  Playlist
	skippers  []string
	dontSkip  bool
}

// YouTubePlaylist implements the Playlist interface
type YouTubePlaylist struct {
	id    string
	title string
}

// ------------
// YOUTUBE SONG
// ------------

// Download downloads the song via youtube-dl if it does not already exist on disk.
// All downloaded songs are stored in ~/.mumbledj/songs and should be automatically cleaned.
func (dl *YouTubeSong) Download() error {

	// Checks to see if song is already downloaded
	if _, err := os.Stat(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, dl.Filename())); os.IsNotExist(err) {
		cmd := exec.Command("youtube-dl", "--no-mtime", "--output", fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, dl.Filename()), "--format", dl.format, "--prefer-ffmpeg", dl.url)
		output, err := cmd.CombinedOutput()
		if err == nil {
			if dj.conf.Cache.Enabled {
				dj.cache.CheckMaximumDirectorySize()
			}
			return nil
		} else {
			args := ""
			for s := range cmd.Args {
				args += cmd.Args[s] + " "
			}
			fmt.Printf(args + "\n" + string(output) + "\n" + "youtube-dl: " + err.Error() + "\n")
			return errors.New("Song download failed.")
		}
	}
	return nil
}

// Play plays the song. Once the song is playing, a notification is displayed in a text message that features the song
// thumbnail, URL, title, duration, and submitter.
func (dl *YouTubeSong) Play() {
	if dl.offset != 0 {
		offsetDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", dl.offset))
		dj.audioStream.Offset = offsetDuration
	}
	dj.audioStream.Source = gumble_ffmpeg.SourceFile(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, dl.Filename()))
	if err := dj.audioStream.Play(); err != nil {
		panic(err)
	} else {
		message := `<table><tr><td align="center"><img src="%s" width=150 /></td></tr><tr><td align="center"><b><a href="%s">%s</a> (%s)</b></td></tr><tr><td align="center">Added by %s</td></tr>`
		message = fmt.Sprintf(message, dl.thumbnail, dl.url, dl.title, dl.duration, dl.submitter.Name)
		if !isNil(dl.playlist) {
			message = fmt.Sprintf(message+`<tr><td align="center">From playlist "%s"</td></tr>`, dl.playlist.Title())
		}
		dj.client.Self.Channel.Send(message+`</table>`, false)

		go func() {
			dj.audioStream.Wait()
			dj.queue.OnSongFinished()
		}()
	}
}

// Delete deletes the song from ~/.mumbledj/songs if the cache is disabled.
func (dl *YouTubeSong) Delete() error {
	if dj.conf.Cache.Enabled == false {
		filePath := fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, dl.Filename())
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err == nil {
				return nil
			}
			return errors.New("Error occurred while deleting audio file.")
		}
		return nil
	}
	return nil
}

// AddSkip adds a skip to the skippers slice. If the user is already in the slice, AddSkip
// returns an error and does not add a duplicate skip.
func (dl *YouTubeSong) AddSkip(username string) error {
	for _, user := range dl.skippers {
		if username == user {
			return errors.New("This user has already skipped the current song.")
		}
	}
	dl.skippers = append(dl.skippers, username)
	return nil
}

// RemoveSkip removes a skip from the skippers slice. If username is not in slice, an error is
// returned.
func (dl *YouTubeSong) RemoveSkip(username string) error {
	for i, user := range dl.skippers {
		if username == user {
			dl.skippers = append(dl.skippers[:i], dl.skippers[i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

// SkipReached calculates the current skip ratio based on the number of users within MumbleDJ's
// channel and the number of usernames in the skippers slice. If the value is greater than or equal
// to the skip ratio defined in the config, the function returns true, and returns false otherwise.
func (dl *YouTubeSong) SkipReached(channelUsers int) bool {
	if float32(len(dl.skippers))/float32(channelUsers) >= dj.conf.General.SkipRatio {
		return true
	}
	return false
}

// Submitter returns the name of the submitter of the Song.
func (dl *YouTubeSong) Submitter() string {
	return dl.submitter.Name
}

// Title returns the title of the Song.
func (dl *YouTubeSong) Title() string {
	return dl.title
}

// ID returns the id of the Song.
func (dl *YouTubeSong) ID() string {
	return dl.id
}

// Filename returns the filename of the Song.
func (dl *YouTubeSong) Filename() string {
	return dl.id + dl.format
}

// Duration returns the duration of the Song.
func (dl *YouTubeSong) Duration() string {
	return dl.duration
}

// Thumbnail returns the thumbnail URL for the Song.
func (dl *YouTubeSong) Thumbnail() string {
	return dl.thumbnail
}

// Playlist returns the playlist type for the Song (may be nil).
func (dl *YouTubeSong) Playlist() Playlist {
	return dl.playlist
}

// DontSkip returns the DontSkip boolean value for the Song.
func (dl *YouTubeSong) DontSkip() bool {
	return dl.dontSkip
}

// SetDontSkip sets the DontSkip boolean value for the Song.
func (dl *YouTubeSong) SetDontSkip(value bool) {
	dl.dontSkip = value
}

// ----------------
// YOUTUBE PLAYLIST
// ----------------

// AddSkip adds a skip to the playlist's skippers slice.
func (p *YouTubePlaylist) AddSkip(username string) error {
	for _, user := range dj.playlistSkips[p.ID()] {
		if username == user {
			return errors.New("This user has already skipped the current song.")
		}
	}
	dj.playlistSkips[p.ID()] = append(dj.playlistSkips[p.ID()], username)
	return nil
}

// RemoveSkip removes a skip from the playlist's skippers slice. If username is not in the slice
// an error is returned.
func (p *YouTubePlaylist) RemoveSkip(username string) error {
	for i, user := range dj.playlistSkips[p.ID()] {
		if username == user {
			dj.playlistSkips[p.ID()] = append(dj.playlistSkips[p.ID()][:i], dj.playlistSkips[p.ID()][i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

// DeleteSkippers removes the skippers entry in dj.playlistSkips.
func (p *YouTubePlaylist) DeleteSkippers() {
	delete(dj.playlistSkips, p.ID())
}

// SkipReached calculates the current skip ratio based on the number of users within MumbleDJ's
// channel and the number of usernames in the skippers slice. If the value is greater than or equal
// to the skip ratio defined in the config, the function returns true, and returns false otherwise.
func (p *YouTubePlaylist) SkipReached(channelUsers int) bool {
	if float32(len(dj.playlistSkips[p.ID()]))/float32(channelUsers) >= dj.conf.General.PlaylistSkipRatio {
		return true
	}
	return false
}

// ID returns the id of the YouTubePlaylist.
func (p *YouTubePlaylist) ID() string {
	return p.id
}

// Title returns the title of the YouTubePlaylist.
func (p *YouTubePlaylist) Title() string {
	return p.title
}
