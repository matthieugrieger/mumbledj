/*
 * MumbleDJ
 * By Matthieu Grieger
 * service_youtube.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"math"

	"github.com/jmoiron/jsonq"
	"github.com/layeh/gumble/gumble"
)

// Regular expressions for youtube urls
var youtubePlaylistPattern = `https?:\/\/www\.youtube\.com\/playlist\?list=([\w-]+)`
var youtubeVideoPatterns = []string{
	`https?:\/\/www\.youtube\.com\/watch\?v=([\w-]+)(\&t=\d*m?\d*s?)?`,
	`https?:\/\/youtube\.com\/watch\?v=([\w-]+)(\&t=\d*m?\d*s?)?`,
	`https?:\/\/youtu.be\/([\w-]+)(\?t=\d*m?\d*s?)?`,
	`https?:\/\/youtube.com\/v\/([\w-]+)(\?t=\d*m?\d*s?)?`,
	`https?:\/\/www.youtube.com\/v\/([\w-]+)(\?t=\d*m?\d*s?)?`,
}

// SearchService name
var youtubeSearchServiceName = "yt"

// SearchService video URL prefix
var youtubeVideoURLprefix = "https://www.youtube.com/watch?v="

// YouTube implements the Service interface
type YouTube struct{}

// ServiceName is the human readable version of the service name
func (yt YouTube) ServiceName() string {
	return "YouTube"
}

// TrackName is the human readable version of the service name
func (yt YouTube) TrackName() string {
	return "Video"
}

// URLRegex checks to see if service will accept URL
func (yt YouTube) URLRegex(url string) bool {
	return RegexpFromURL(url, append(youtubeVideoPatterns, []string{youtubePlaylistPattern}...)) != nil
}

// SearchRegex checks to see if service will accept the searchString 
func (yt YouTube) SearchRegex(searchService string) bool {
	return searchService == youtubeSearchServiceName
}

// NewRequest creates the requested song/playlist and adds to the queue
func (yt YouTube) NewRequest(user *gumble.User, url string) ([]Song, error) {
	var songArray []Song
	var shortURL, startOffset = "", ""
	if re, err := regexp.Compile(youtubePlaylistPattern); err == nil {
		if re.MatchString(url) {
			shortURL = re.FindStringSubmatch(url)[1]
			return yt.NewPlaylist(user, shortURL)
		} else {
			re = RegexpFromURL(url, youtubeVideoPatterns)
			matches := re.FindAllStringSubmatch(url, -1)
			shortURL = matches[0][1]
			if len(matches[0]) == 3 {
				startOffset = matches[0][2]
			}
			song, err := yt.NewSong(user, shortURL, startOffset, nil)
			if !isNil(song) {
				return append(songArray, song), nil
			} else {
				return nil, err
			}
		}
	} else {
		return nil, err
	}
}

// SearchSong searches for a Song and returns the Songs URL
func (yt YouTube) SearchSong(searchString string) (string, error) {
	var returnString, apiStringValue string
	searchURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&key=%s&maxResults=1&type=video", searchString, dj.conf.ServiceKeys.Youtube)

	if apiResponse, err := PerformGetRequest(searchURL); err == nil {
		apiStringValue, _ = apiResponse.String("items", "0", "id", "videoId")
		returnString = youtubeVideoURLprefix + apiStringValue
		return returnString, nil
	}
	return "", errors.New(fmt.Sprintf(INVALID_API_KEY, yt.ServiceName()))
}

// NewSong gathers the metadata for a song extracted from a YouTube video, and returns the song.
func (yt YouTube) NewSong(user *gumble.User, id, offset string, playlist Playlist) (Song, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails&id=%s&key=%s", id, dj.conf.ServiceKeys.Youtube)
	if apiResponse, err := PerformGetRequest(url); err == nil {
		title, _ := apiResponse.String("items", "0", "snippet", "title")
		thumbnail, _ := apiResponse.String("items", "0", "snippet", "thumbnails", "high", "url")
		duration, _ := apiResponse.String("items", "0", "contentDetails", "duration")

		song := &AudioTrack{
			submitter: user,
			title:     title,
			id:        id,
			url:       "https://youtu.be/" + id,
			offset:    int(yt.parseTime(offset, `T\=(?P<days>\d+D)?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`).Seconds()),
			duration:  int(yt.parseTime(duration, `P(?P<days>\d+D)?T(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`).Seconds()),
			thumbnail: thumbnail,
			format:    "bestaudio",
			skippers:  make([]string, 0),
			playlist:  playlist,
			dontSkip:  false,
			service:   yt,
		}

		return song, nil
	}
	return nil, errors.New(fmt.Sprintf(INVALID_API_KEY, yt.ServiceName()))
}

// parseTime converts from the string youtube returns to a time.Duration
func (yt YouTube) parseTime(duration, regex string) time.Duration {
	var days, hours, minutes, seconds, totalSeconds int64
	if duration != "" {
		timestampExp := regexp.MustCompile(regex)
		timestampMatch := timestampExp.FindStringSubmatch(strings.ToUpper(duration))
		timestampResult := make(map[string]string)
		for i, name := range timestampExp.SubexpNames() {
			if i < len(timestampMatch) {
				timestampResult[name] = timestampMatch[i]
			}
		}

		if timestampResult["days"] != "" {
			days, _ = strconv.ParseInt(strings.TrimSuffix(timestampResult["days"], "D"), 10, 32)
		}
		if timestampResult["hours"] != "" {
			hours, _ = strconv.ParseInt(strings.TrimSuffix(timestampResult["hours"], "H"), 10, 32)
		}
		if timestampResult["minutes"] != "" {
			minutes, _ = strconv.ParseInt(strings.TrimSuffix(timestampResult["minutes"], "M"), 10, 32)
		}
		if timestampResult["seconds"] != "" {
			seconds, _ = strconv.ParseInt(strings.TrimSuffix(timestampResult["seconds"], "S"), 10, 32)
		}

		totalSeconds = int64((days * 86400) + (hours * 3600) + (minutes * 60) + seconds)
	} else {
		totalSeconds = 0
	}
	output, _ := time.ParseDuration(strconv.Itoa(int(totalSeconds)) + "s")
	return output
}

// NewPlaylist gathers the metadata for a YouTube playlist and returns it.
func (yt YouTube) NewPlaylist(user *gumble.User, id string) ([]Song, error) {
	var apiResponse *jsonq.JsonQuery
	var songArray []Song
	var err error
	// Retrieve title of playlist
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlists?part=snippet&id=%s&key=%s", id, dj.conf.ServiceKeys.Youtube)
	if apiResponse, err = PerformGetRequest(url); err != nil {
		return nil, err
	}
	title, _ := apiResponse.String("items", "0", "snippet", "title")

	playlist := &AudioPlaylist{
		id:    id,
		title: title,
	}


	maxSongs := math.MaxInt32
	if (dj.conf.General.MaxSongPerPlaylist > 0){
		maxSongs = dj.conf.General.MaxSongPerPlaylist
	}
	pageToken := ""
	for len(songArray) < maxSongs{	//Iterate over the pages

		// Retrieve items in this page of the playlist
		url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=50&playlistId=%s&key=%s&pageToken=%s",
			id, dj.conf.ServiceKeys.Youtube, pageToken)
		if apiResponse, err = PerformGetRequest(url); err != nil {
			return nil, err
		}

		songs, _ := apiResponse.Array("items")
		for j := 0; j < len(songs) && len(songArray) < maxSongs ; j++ {
			index := strconv.Itoa(j)
			videoID, _ := apiResponse.String("items", index, "snippet", "resourceId", "videoId")
			if song, err := yt.NewSong(user, videoID, "", playlist); err == nil {
				songArray = append(songArray, song)
			}
		}
		if pageToken, err = apiResponse.String("nextPageToken"); err != nil {
			break
		}
	}
	return songArray, nil
}
