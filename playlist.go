/*
 * MumbleDJ
 * By Matthieu Grieger
 * playlist.go
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
	"strconv"
	"strings"
)

// Playlist type declaration.
type Playlist struct {
	songs     *SongQueue
	youtubeId string
	title     string
	submitter string
	skippers  []string
	skipped   bool
}

// Returns a new Playlist type. Before returning the new type, the playlist's metadata is collected
// via the YouTube Gdata API.
func NewPlaylist(user, id string) *Playlist {
	queue := NewSongQueue()
	jsonUrl := fmt.Sprintf("http://gdata.youtube.com/feeds/api/playlists/%s?v=2&alt=jsonc&maxresults=25", id)
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

	playlistTitle, _ := jq.String("data", "title")
	playlistItems, _ := jq.Int("data", "totalItems")
	if playlistItems > 25 {
		playlistItems = 25
	}

	for i := 0; i < playlistItems; i++ {
		index := strconv.Itoa(i)
		songTitle, _ := jq.String("data", "items", index, "video", "title")
		songId, _ := jq.String("data", "items", index, "video", "id")
		songThumbnail, _ := jq.String("data", "items", index, "video", "thumbnail", "hqDefault")
		duration, _ := jq.Int("data", "items", index, "video", "duration")
		songDuration := fmt.Sprintf("%d:%02d", duration/60, duration%60)
		newSong := &Song{
			submitter:    user,
			title:        songTitle,
			youtubeId:    songId,
			playlistId:   id,
			duration:     songDuration,
			thumbnailUrl: songThumbnail,
		}
		queue.AddItem(newSong)
	}

	playlist := &Playlist{
		songs:     queue,
		youtubeId: id,
		title:     playlistTitle,
		submitter: user,
		skipped:   false,
	}
	return playlist
}

// Adds a skip to the skippers slice. If the user is already in the slice AddSkip() returns
// an error and does not add a duplicate skip.
func (p *Playlist) AddSkip(username string) error {
	for _, user := range p.skippers {
		if username == user {
			return errors.New("This user has already skipped the current song.")
		}
	}
	p.skippers = append(p.skippers, username)
	return nil
}

// Removes a skip from the skippers slice. If username is not in the slice, an error is
// returned.
func (p *Playlist) RemoveSkip(username string) error {
	for i, user := range p.skippers {
		if username == user {
			p.skippers = append(p.skippers[:i], p.skippers[i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

// Calculates current skip ratio based on number of users within MumbleDJ's channel and the
// amount of values in the skippers slice. If the value is greater than or equal to the skip ratio
// defined in mumbledj.gcfg, the function returns true. Returns false otherwise.
func (p *Playlist) SkipReached(channelUsers int) bool {
	if float32(len(p.skippers))/float32(channelUsers) >= dj.conf.General.PlaylistSkipRatio {
		return true
	} else {
		return false
	}
}

// Returns "playlist" as the item type. Used for differentiating Songs from Playlists.
func (p *Playlist) ItemType() string {
	return "playlist"
}
