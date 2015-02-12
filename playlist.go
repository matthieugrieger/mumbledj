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
	id    string
	title string
}

// Returns a new Playlist type. Before returning the new type, the playlist's metadata is collected
// via the YouTube Gdata API.
func NewPlaylist(user, id string) (*Playlist, error) {
	jsonUrl := fmt.Sprintf("http://gdata.youtube.com/feeds/api/playlists/%s?v=2&alt=jsonc&maxresults=25", id)
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

	playlistTitle, _ := jq.String("data", "title")
	playlistItems, _ := jq.Int("data", "totalItems")
	if playlistItems > 25 {
		playlistItems = 25
	}

	playlist := &Playlist{
		id:    id,
		title: playlistTitle,
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
			duration:     songDuration,
			thumbnailUrl: songThumbnail,
			playlist:     playlist,
			dontSkip:     false,
		}
		dj.queue.AddSong(newSong)
	}

	return playlist, nil
}

// Adds a skip to the skippers slice for the current playlist.
func (p *Playlist) AddSkip(username string) error {
	for _, user := range dj.playlistSkips[p.id] {
		if username == user {
			return errors.New("This user has already skipped the current song.")
		}
	}
	dj.playlistSkips[p.id] = append(dj.playlistSkips[p.id], username)
	return nil
}

// Removes a skip from the skippers slice. If username is not in the slice, an error is
// returned.
func (p *Playlist) RemoveSkip(username string) error {
	for i, user := range dj.playlistSkips[p.id] {
		if username == user {
			dj.playlistSkips[p.id] = append(dj.playlistSkips[p.id][:i], dj.playlistSkips[p.id][i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

// Removes skippers entry in dj.playlistSkips.
func (p *Playlist) DeleteSkippers() {
	delete(dj.playlistSkips, p.id)
}

// Calculates current skip ratio based on number of users within MumbleDJ's channel and the
// amount of values in the skippers slice. If the value is greater than or equal to the skip ratio
// defined in mumbledj.gcfg, the function returns true. Returns false otherwise.
func (p *Playlist) SkipReached(channelUsers int) bool {
	if float32(len(dj.playlistSkips[p.id]))/float32(channelUsers) >= dj.conf.General.PlaylistSkipRatio {
		return true
	} else {
		return false
	}
}
