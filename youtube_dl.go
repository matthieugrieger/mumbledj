/*
 * MumbleDJ
 * By Matthieu Grieger
 * youtube_dl.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/jsonq"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
)

// AudioTrack implements the Song interface
type AudioTrack struct {
	id        string
	title     string
	thumbnail string
	submitter *gumble.User
	duration  int
	url       string
	offset    int
	format    string
	playlist  Playlist
	skippers  []string
	dontSkip  bool
	service   Service
}

// AudioPlaylist implements the Playlist interface
type AudioPlaylist struct {
	id    string
	title string
}

// ------------
// YOUTUBEDL SONG
// ------------

// Download downloads the song via youtube-dl if it does not already exist on disk.
// All downloaded songs are stored in ~/.mumbledj/songs and should be automatically cleaned.
func (dl *AudioTrack) Download() error {
	player := "--prefer-ffmpeg"
	if dj.conf.General.PlayerCommand == "avconv" {
		player = "--prefer-avconv"
	}

	// Checks to see if song is already downloaded
	if _, err := os.Stat(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, dl.Filename())); os.IsNotExist(err) {
		cmd := exec.Command("youtube-dl", "--verbose", "--no-mtime", "--output", fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, dl.Filename()), "--format", dl.format, player, dl.url)
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
func (dl *AudioTrack) Play() {
	if dl.offset != 0 {
		offsetDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", dl.offset))
		dj.audioStream.Offset = offsetDuration
	}
	dj.audioStream.Source = gumble_ffmpeg.SourceFile(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, dl.Filename()))
	if err := dj.audioStream.Play(); err != nil {
		panic(err)
	} else {
		message := `<table><tr><td align="center"><img src="%s" width=150 /></td></tr><tr><td align="center"><b><a href="%s">%s</a> (%s)</b></td></tr><tr><td align="center">Added by %s</td></tr>`
		message = fmt.Sprintf(message, dl.thumbnail, dl.url, dl.title, dl.Duration().String(), dl.submitter.Name)
		if !isNil(dl.playlist) {
			message = fmt.Sprintf(message+`<tr><td align="center">From playlist "%s"</td></tr>`, dl.Playlist().Title())
		}
		if dj.conf.General.AnnounceNewTrack {
			dj.client.Self.Channel.Send(message+`</table>`, false)
		}
		go func() {
			dj.audioStream.Wait()
			dj.queue.OnSongFinished()
		}()
	}
}

// Delete deletes the song from ~/.mumbledj/songs if the cache is disabled.
func (dl *AudioTrack) Delete() error {
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
func (dl *AudioTrack) AddSkip(username string) error {
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
func (dl *AudioTrack) RemoveSkip(username string) error {
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
func (dl *AudioTrack) SkipReached(channelUsers int) bool {
	if float32(len(dl.skippers))/float32(channelUsers) >= dj.conf.General.SkipRatio {
		return true
	}
	return false
}

// Submitter returns the name of the submitter of the Song.
func (dl *AudioTrack) Submitter() string {
	return dl.submitter.Name
}

// Title returns the title of the Song.
func (dl *AudioTrack) Title() string {
	return dl.title
}

// ID returns the id of the Song.
func (dl *AudioTrack) ID() string {
	return dl.id
}

// Filename returns the filename of the Song.
func (dl *AudioTrack) Filename() string {
	return dl.id + "." + dl.format
}

// Duration returns duration for the Song.
func (dl *AudioTrack) Duration() time.Duration {
	timeDuration, _ := time.ParseDuration(strconv.Itoa(dl.duration) + "s")
	return timeDuration
}

// Thumbnail returns the thumbnail URL for the Song.
func (dl *AudioTrack) Thumbnail() string {
	return dl.thumbnail
}

// Playlist returns the playlist type for the Song (may be nil).
func (dl *AudioTrack) Playlist() Playlist {
	return dl.playlist
}

// DontSkip returns the DontSkip boolean value for the Song.
func (dl *AudioTrack) DontSkip() bool {
	return dl.dontSkip
}

// SetDontSkip sets the DontSkip boolean value for the Song.
func (dl *AudioTrack) SetDontSkip(value bool) {
	dl.dontSkip = value
}

// ----------------
// YOUTUBEDL PLAYLIST
// ----------------

// AddSkip adds a skip to the playlist's skippers slice.
func (p *AudioPlaylist) AddSkip(username string) error {
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
func (p *AudioPlaylist) RemoveSkip(username string) error {
	for i, user := range dj.playlistSkips[p.ID()] {
		if username == user {
			dj.playlistSkips[p.ID()] = append(dj.playlistSkips[p.ID()][:i], dj.playlistSkips[p.ID()][i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

// DeleteSkippers removes the skippers entry in dj.playlistSkips.
func (p *AudioPlaylist) DeleteSkippers() {
	delete(dj.playlistSkips, p.ID())
}

// SkipReached calculates the current skip ratio based on the number of users within MumbleDJ's
// channel and the number of usernames in the skippers slice. If the value is greater than or equal
// to the skip ratio defined in the config, the function returns true, and returns false otherwise.
func (p *AudioPlaylist) SkipReached(channelUsers int) bool {
	if float32(len(dj.playlistSkips[p.ID()]))/float32(channelUsers) >= dj.conf.General.PlaylistSkipRatio {
		return true
	}
	return false
}

// ID returns the id of the AudioPlaylist.
func (p *AudioPlaylist) ID() string {
	return p.id
}

// Title returns the title of the AudioPlaylist.
func (p *AudioPlaylist) Title() string {
	return p.title
}

// PerformGetRequest does all the grunt work for HTTPS GET request.
func PerformGetRequest(url string) (*jsonq.JsonQuery, error) {
	jsonString := ""

	if response, err := http.Get(url); err == nil {
		defer response.Body.Close()
		if response.StatusCode == 200 {
			if body, err := ioutil.ReadAll(response.Body); err == nil {
				jsonString = string(body)
				if jsonString[0] == '[' {
					jsonString = "{\"json\":" + jsonString + "}"
				}
			}
		} else {
			if response.StatusCode == 403 {
				return nil, errors.New("Invalid API key supplied.")
			}
			return nil, errors.New("Invalid ID supplied.")
		}
	} else {
		return nil, errors.New("An error occurred while receiving HTTP GET response.")
	}

	jsonData := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(jsonString))
	decoder.Decode(&jsonData)
	jq := jsonq.NewQuery(jsonData)

	return jq, nil
}
