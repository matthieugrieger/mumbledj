/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/youtube/api.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/jmoiron/jsonq"
)

// Video holds the metadata for a YouTube video.
type Video struct {
	id              string
	title           string
	duration        string
	secondsDuration string
	thumbnail       string
}

// Playlist holds the metadata for a YouTube playlist.
type Playlist struct {
	id       string
	title    string
	videoIds []string
}

// GetYouTubeVideo retrieves the metadata for a new YouTube video, and creates and returns a
// Video type.
func GetYouTubeVideo(id string) (*YouTubeVideo, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails&id=%s&key=%s",
		id, os.Getenv("YOUTUBE_API_KEY"))
	if response, err := PerformGetRequest(url); err != nil {
		return nil, err
	}

	title, _ := response.String("items", "0", "snippet", "title")
	thumbnail, _ := response.String("items", "0", "snippet", "thumbnails", "high", "url")
	duration, _ := response.String("items", "0", "contentDetails", "duration")

	minutes := int(duration[2:strings.Index(duration, "M")])
	seconds := int(duration[strings.Index(duration, "M")+1 : len(duration)-1])
	totalSeconds := (minutes * 60) + seconds
	durationString := fmt.Sprintf("%d:%d", minutes, seconds)

	video := &YoutubeVideo{
		id:              id,
		title:           title,
		duration:        durationString,
		secondsDuration: totalSeconds,
		thumbnail:       thumbnail,
	}
	return video, nil
}

// GetYouTubePlaylist retrieves the metadata for a new YouTube playlist, and creates and returns
// a Playlist type.
func GetYouTubePlaylist(id string) (*YouTubePlaylist, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=25&playlistId=%s&key=%s",
		id, os.Getenv("YOUTUBE_API_KEY"))
	if response, err := PerformGetRequest(url); err != nil {
		return nil, err
	}

	title, _ := response.String("items", "0")

}

// PerformGetRequest does all the grunt work for a Youtube HTTPS GET request.
func PerformGetRequest(url string) (*jsonq.JsonQuery, error) {
	jsonString := ""

	if response, err := http.Get(url); err == nil {
		defer response.Body.Close()
		if response.StatusCode == 200 {
			if body, err := ioutil.ReadAll(response.Body); err == nil {
				jsonString = string(body)
			}
		} else {
			if response.StatusCode == 403 {
				return nil, errors.New("Invalid API key supplied.")
			}
			return nil, errors.New("Invalid YouTube ID supplied.")
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
