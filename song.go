/*
 * MumbleDJ
 * By Matthieu Grieger
 * song.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
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

type Song struct {
	submitter    string
	title        string
	youtubeId    string
	duration     string
	thumbnailUrl string
	skippers     []string
}

func NewSong(user, id string) *Song {
	jsonUrl := fmt.Sprintf("http://gdata.youtube.com/feeds/api/videos/%s?v=2&alt=jsonc", id)
	jsonString := ""

	if response, err := http.Get(jsonUrl); err == nil {
		defer response.Body.Close()
		if body, err := ioutil.ReadAll(response.Body); err == nil {
			jsonString = string(body)
		}
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
	}
	return song
}

func (s *Song) Download() error {
	cmd := exec.Command("youtube-dl", "--output", fmt.Sprintf(`~/.mumbledj/songs/%s.m4a`, s.youtubeId), "--format", "m4a", s.youtubeId)
	if err := cmd.Run(); err == nil {
		return nil
	} else {
		return errors.New("Song download failed.")
	}
}

func (s *Song) Play() {
	dj.audioStream.Play(fmt.Sprintf("%s/.mumbledj/songs/%s.m4a", dj.homeDir, s.youtubeId))
	dj.client.Self().Channel().Send(fmt.Sprintf(NOW_PLAYING_HTML, s.thumbnailUrl, s.youtubeId, s.title, s.duration, s.submitter), false)
}

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

func (s *Song) AddSkip(username string) error {
	for _, user := range s.skippers {
		if username == user {
			return errors.New("This user has already skipped the current song.")
		}
	}
	s.skippers = append(s.skippers, username)
	return nil
}

func (s *Song) RemoveSkip(username string) error {
	for i, user := range s.skippers {
		if username == user {
			s.skippers = append(s.skippers[:i], s.skippers[i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

func (s *Song) SkipReached(channelUsers int) bool {
	if float32(len(s.skippers))/float32(channelUsers) >= dj.conf.General.SkipRatio {
		return true
	} else {
		return false
	}
}
