/*
 * MumbleDJ
 * By Matthieu Grieger
 * song.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

import (
	//"github.com/layeh/gumble/gumble_ffmpeg"
	"os/exec"
	"fmt"
)

type Song struct {
	submitter string
	title string
	youtubeId string
	duration string
	thumbnailUrl string
	skippers []string
}

func NewSong(user, id string) *Song {
	song := &Song{
		submitter: user,
		youtubeId: id,
	}
	return song
}

func (s *Song) Download() bool {
	err := exec.Command(fmt.Sprintf("youtube-dl --output ~/.mumbledj/songs/%(id)s.%(ext)s --quiet --format bestaudio --audio-format vorbis --prefer-ffmpeg https://youtube.com/watch?v=%s", s.youtubeId))
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
	return false
}

func (s *Song) AddSkip(username string) bool {
	for _,user := range s.skippers {
		if username == user {
			return false
		}
	}
	s.skippers = append(s.skippers, username)
	return true
}

func (s *Song) RemoveSkip(username string) bool {
	for i,user := range s.skippers {
		if username == user {
			s.skippers = append(s.skippers[:i], s.skippers[i+1:]...)
			return true
		}
	}
	return false
}

