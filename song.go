/*
 * MumbleDJ
 * By Matthieu Grieger
 * song.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

import (
	//"github.com/layeh/gumble/gumble_ffmpeg"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
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
	response, err := http.Get(jsonUrl)
	jsonString := ""
	if err == nil {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err == nil {
			jsonString = string(body)
		}
	}

	jsonData := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(jsonString))
	decoder.Decode(&jsonData)
	jq := jsonq.NewQuery(jsonData)

	videoTitle, _ := jq.String("data", "title")
	videoThumbnail, _ := jq.String("data", "thumbnail", "sqDefault")
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

func (s *Song) Download() bool {
	err := exec.Command(fmt.Sprintf("youtube-dl --output \"~/.mumbledj/songs/%(id)s.%(ext)s\" --quiet --format m4a %s", s.youtubeId))
	if err == nil {
		return true
	} else {
		return false
	}
}

func (s *Song) Play() bool {
	return false
}

func (s *Song) Delete() bool {
	usr, err := user.Current()
	if err == nil {
		filePath := fmt.Sprintf("%s/.mumbledj/songs/%s.m4a", usr.HomeDir, s.youtubeId)
		if _, err := os.Stat(filePath); err == nil {
			err := os.Remove(filePath)
			if err == nil {
				return true
			} else {
				return false
			}
		} else {
			return true
		}
	} else {
		return false
	}
}

func (s *Song) AddSkip(username string) bool {
	for _, user := range s.skippers {
		if username == user {
			return false
		}
	}
	s.skippers = append(s.skippers, username)
	return true
}

func (s *Song) RemoveSkip(username string) bool {
	for i, user := range s.skippers {
		if username == user {
			s.skippers = append(s.skippers[:i], s.skippers[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Song) SkipReached(channelUsers int) bool {
	return false
}
