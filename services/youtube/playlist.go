/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/youtube/playlist.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package youtube

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Playlist holds the metadata for a YouTube playlist.
type Playlist struct {
	id    string
	title string
}

// NewPlaylist gathers the metadata for a YouTube playlist and returns it.
func NewPlaylist(user, id string) (*Playlist, error) {
	// Retrieve title of playlist
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlists?part=snippet&id=%s&key=%s",
		id, os.Getenv("YOUTUBE_API_KEY"))
	if response, err := PerformGetRequest(url); err != nil {
		return nil, err
	}
	title, _ := response.String("items", "0")

	playlist := &Playlist{
		id:    id,
		title: title,
	}

	// Retrieve items in playlist
	url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=25&playlistId=%s&key=%s",
		id, os.Getenv("YOUTUBE_API_KEY"))
	if response, err = PerformGetRequest(url); err != nil {
		return nil, err
	}
	numVideos := response.Int("pageInfo", "totalResults")
	if numVideos > 25 {
		numVideos = 25
	}

	for i := 0; i < numVideos; i++ {
		index := strconv.Itoa(i)
		videoTitle, _ := response.String("items", index, "snippet", "title")
		videoID, _ := response.String("items", index, "snippet", "resourceId", "videoId")
		videoThumbnail, _ := response.String("items", index, "snippet", "thumbnails", "high", "url")

		// A completely separate API call just to get the duration of a video in a
		// playlist? WHY GOOGLE, WHY?!
		url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=contentDetails&id=%s&key=%s",
			videoID, os.Getenv("YOUTUBE_API_KEY"))
		if response, err = PerformGetRequest(url); err != nil {
			return nil, err
		}
		videoDuration, _ := response.String("items", "0", "contentDetails", "duration")
		minutes := int(videoDuration[2:strings.Index(videoDuration, "M")])
		seconds := int(videoDuration[strings.Index(videoDuration, "M")+1 : len(videoDuration)-1])
		totalSeconds := (minutes * 60) + seconds
		durationString := fmt.Sprintf("%d:%d", minutes, seconds)

		if dj.conf.General.MaxSongDuration == 0 || totalSeconds <= dj.conf.General.MaxSongDuration {
			playlistSong := &Song{
				submitter: user,
				title:     videoTitle,
				id:        videoID,
				duration:  durationString,
				thumbnail: videoThumbnail,
				skippers:  make([]string, 0),
				playlist:  playlist,
				dontSkip:  false,
			}
			dj.queue.AddSong(playlistSong)
		}
	}
	return playlist, nil
}

// AddSkip adds a skip to the playlist's skippers slice.
func (p *Playlist) AddSkip(username string) error {
	for _, user := range dj.playlistSkips[p.id] {
		if username == user {
			return errors.New("This user has already skipped the current song.")
		}
	}
	dj.playlistSkips[p.id] = append(dj.playlistSkips[p.id], username)
	return nil
}

// RemoveSkip removes a skip from the playlist's skippers slice. If username is not in the slice
// an error is returned.
func (p *Playlist) RemoveSkip(username string) error {
	for i, user := range dj.playlistSkips[p.id] {
		if username == user {
			dj.playlistSkips[p.id] = append(dj.playlistSkips[p.id][:i], dj.playlistSkips[p.id][i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

// DeleteSkippers removes the skippers entry in dj.playlistSkips.
func (p *Playlist) DeleteSkippers() {
	delete(dj.playlistSkips, p.id)
}

// SkipReached calculates the current skip ratio based on the number of users within MumbleDJ's
// channel and the number of usernames in the skippers slice. If the value is greater than or equal
// to the skip ratio defined in the config, the function returns true, and returns false otherwise.
func (p *Playlist) SkipReached(channelUsers int) bool {
	if float32(len(dj.playlistSkips[p.id]))/float32(channelUsers) >= dj.conf.General.PlaylistSkipRatio {
		return true
	}
	return false
}
