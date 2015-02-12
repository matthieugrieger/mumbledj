/*
 * MumbleDJ
 * By Matthieu Grieger
 * song.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Song type declaration.
type Song struct {
	submitter    string
	title        string
	youtubeId    string
	duration     string
	thumbnailUrl string
	skippers     []string
	playlist     *Playlist
	dontSkip     bool
}

// Returns a new Song type. Before returning the new type, the song's metadata is collected
// via the YouTube Gdata API.
func NewSong(user, id string, playlist *Playlist) (*Song, error) {
	jsonUrl := fmt.Sprintf("http://gdata.youtube.com/feeds/api/videos/%s?v=2&alt=jsonc", id)
	jsonString := ""

	if response, err := http.Get(jsonUrl); err == nil {
		defer response.Body.Close()
		if response.StatusCode != 400 && response.StatusCode != 404 {
			if body, err := ioutil.ReadAll(response.Body); err == nil {
				jsonString = string(body)
			}
		} else {
			return nil, errors.New("Invalid YouTube ID supplied.")
		}
	} else {
		return nil, errors.New("An error occurred while receiving HTTP GET request.")
	}

	jsonData := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(jsonString))
	decoder.Decode(&jsonData)
	jq := jsonq.NewQuery(jsonData)

	videoTitle, _ := jq.String("data", "title")
	videoThumbnail, _ := jq.String("data", "thumbnail", "hqDefault")
	duration, _ := jq.Int("data", "duration")
	videoDuration := fmt.Sprintf("%d:%02d", duration/60, duration%60)

	song := &Song{
		submitter:    user,
		title:        videoTitle,
		youtubeId:    id,
		duration:     videoDuration,
		thumbnailUrl: videoThumbnail,
		playlist:     playlist,
		dontSkip:     false,
	}
	return song, nil
}

// Downloads the song via youtube-dl. All downloaded songs are stored in ~/.mumbledj/songs and should be automatically cleaned.
func (s *Song) Download() error {
	cmd := exec.Command("youtube-dl", "--output", fmt.Sprintf(`~/.mumbledj/songs/%s.m4a`, s.youtubeId), "--format", "m4a", s.youtubeId)
	if err := cmd.Run(); err == nil {
		return nil
	} else {
		return errors.New("Song download failed.")
	}
}

// Plays the song. Once the song is playing, a notification is displayed in a text message that features the video thumbnail, URL, title,
// duration, and submitter.
func (s *Song) Play() {
	if err := dj.audioStream.Play(fmt.Sprintf("%s/.mumbledj/songs/%s.m4a", dj.homeDir, s.youtubeId), dj.queue.OnSongFinished); err != nil {
		panic(err)
	} else {
		if s.playlist == nil {
			dj.client.Self.Channel.Send(fmt.Sprintf(NOW_PLAYING_HTML, s.thumbnailUrl, s.youtubeId, s.title,
				s.duration, s.submitter), false)
		} else {
			dj.client.Self.Channel.Send(fmt.Sprintf(NOW_PLAYING_PLAYLIST_HTML, s.thumbnailUrl, s.youtubeId,
				s.title, s.duration, s.submitter, s.playlist.title), false)
		}
	}
}

// Deletes the song from ~/.mumbledj/songs.
func (s *Song) Delete() error {
	filePath := fmt.Sprintf("%s/.mumbledj/songs/%s.m4a", dj.homeDir, s.youtubeId)
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err == nil {
			return nil
		} else {
			return errors.New("Error occurred while deleting audio file.")
		}
	} else {
		return nil
	}
}

// Adds a skip to the skippers slice. If the user is already in the slice AddSkip() returns
// an error and does not add a duplicate skip.
func (s *Song) AddSkip(username string) error {
	for _, user := range s.skippers {
		if username == user {
			return errors.New("This user has already skipped the current song.")
		}
	}
	s.skippers = append(s.skippers, username)
	return nil
}

// Removes a skip from the skippers slice. If username is not in the slice, an error is
// returned.
func (s *Song) RemoveSkip(username string) error {
	for i, user := range s.skippers {
		if username == user {
			s.skippers = append(s.skippers[:i], s.skippers[i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

// Calculates current skip ratio based on number of users within MumbleDJ's channel and the
// amount of values in the skippers slice. If the value is greater than or equal to the skip ratio
// defined in mumbledj.gcfg, the function returns true. Returns false otherwise.
func (s *Song) SkipReached(channelUsers int) bool {
	if float32(len(s.skippers))/float32(channelUsers) >= dj.conf.General.SkipRatio {
		return true
	} else {
		return false
	}
}
