/*
 * MumbleDJ
 * By Matthieu Grieger
 * service_youtube.go
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
	"regexp"
	"strconv"
	"strings"

	"github.com/jmoiron/jsonq"
	"github.com/layeh/gumble/gumble"
)

// Regular expressions for youtube urls
var youtubePlaylistPattern = `https?:\/\/www\.youtube\.com\/playlist\?list=([\w-]+)`
var youtubeVideoPatterns = []string{
	`https?:\/\/(www\.)?youtube\.com\/watch\?v=([\w-]+)(\&t=\d*m?\d*s?)?`,
	`https?:\/\/(www\.)?youtube\.com\/v\/([\w-]+)(\?t=\d*m?\d*s?)?`,
	`https?:\/\/youtu.be\/([\w-]+)(\?t=\d*m?\d*s?)?`,
}

// ------
// TYPES
// ------

// YouTube implements the Service interface
type YouTube struct{}

// ---------------
// YOUTUBE SERVICE
// ---------------

// URLRegex checks to see if service will accept URL
func (yt YouTube) URLRegex(url string) bool {
	return RegexpFromURL(url, append(youtubeVideoPatterns, []string{youtubePlaylistPattern}...)) != nil
}

// NewRequest creates the requested song/playlist and adds to the queue
func (yt YouTube) NewRequest(user *gumble.User, url string) (string, error) {
	var shortURL, startOffset = "", ""
	if re, err := regexp.Compile(youtubePlaylistPattern); err == nil {
		if re.MatchString(url) {
			if dj.HasPermission(user.Name, dj.conf.Permissions.AdminAddPlaylists) {
				shortURL = re.FindStringSubmatch(url)[1]
				playlist, err := yt.NewPlaylist(user, shortURL)
				return playlist.Title(), err
			} else {
				return "", errors.New("NO_PLAYLIST_PERMISSION")
			}
		} else {
			re = RegexpFromURL(url, youtubeVideoPatterns)
			matches := re.FindAllStringSubmatch(url, -1)
			shortURL = matches[0][1]
			if len(matches[0]) == 3 {
				startOffset = matches[0][2]
			}
			song, err := yt.NewSong(user, shortURL, startOffset, nil)
			if !isNil(song) {
				return song.Title(), nil
			} else {
				Verbose("youtube.NewRequest: " + err.Error())
				return "", err
			}
		}
	} else {
		return "", err
	}
}

// NewSong gathers the metadata for a song extracted from a YouTube video, and returns the song.
func (yt YouTube) NewSong(user *gumble.User, id, offset string, playlist Playlist) (Song, error) {
	var apiResponse *jsonq.JsonQuery
	var err error
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails&id=%s&key=%s",
		id, os.Getenv("YOUTUBE_API_KEY"))
	if apiResponse, err = PerformGetRequest(url); err != nil {
		return nil, errors.New(INVALID_API_KEY)
	}

	var offsetDays, offsetHours, offsetMinutes, offsetSeconds int64
	if offset != "" {
		offsetExp := regexp.MustCompile(`t\=(?P<days>\d+d)?(?P<hours>\d+h)?(?P<minutes>\d+m)?(?P<seconds>\d+s)?`)
		offsetMatch := offsetExp.FindStringSubmatch(offset)
		offsetResult := make(map[string]string)
		for i, name := range offsetExp.SubexpNames() {
			if i < len(offsetMatch) {
				offsetResult[name] = offsetMatch[i]
			}
		}

		if offsetResult["days"] != "" {
			offsetDays, _ = strconv.ParseInt(strings.TrimSuffix(offsetResult["days"], "d"), 10, 32)
		}
		if offsetResult["hours"] != "" {
			offsetHours, _ = strconv.ParseInt(strings.TrimSuffix(offsetResult["hours"], "h"), 10, 32)
		}
		if offsetResult["minutes"] != "" {
			offsetMinutes, _ = strconv.ParseInt(strings.TrimSuffix(offsetResult["minutes"], "m"), 10, 32)
		}
		if offsetResult["seconds"] != "" {
			offsetSeconds, _ = strconv.ParseInt(strings.TrimSuffix(offsetResult["seconds"], "s"), 10, 32)
		}
	}

	title, _ := apiResponse.String("items", "0", "snippet", "title")
	thumbnail, _ := apiResponse.String("items", "0", "snippet", "thumbnails", "high", "url")
	duration, _ := apiResponse.String("items", "0", "contentDetails", "duration")

	var days, hours, minutes, seconds int64
	timestampExp := regexp.MustCompile(`P(?P<days>\d+D)?T(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`)
	timestampMatch := timestampExp.FindStringSubmatch(duration)
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

	totalSeconds := int((days * 86400) + (hours * 3600) + (minutes * 60) + seconds)
	var durationString string
	if hours != 0 {
		if days != 0 {
			durationString = fmt.Sprintf("%d:%02d:%02d:%02d", days, hours, minutes, seconds)
		} else {
			durationString = fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
		}
	} else {
		durationString = fmt.Sprintf("%d:%02d", minutes, seconds)
	}

	if dj.conf.General.MaxSongDuration == 0 || totalSeconds <= dj.conf.General.MaxSongDuration {
		song := &YouTubeSong{
			submitter: user,
			title:     title,
			id:        id,
			url:       "https://youtu.be/" + id,
			offset:    int((offsetDays * 86400) + (offsetHours * 3600) + (offsetMinutes * 60) + offsetSeconds),
			duration:  durationString,
			thumbnail: thumbnail,
			format:    "m4a",
			skippers:  make([]string, 0),
			playlist:  playlist,
			dontSkip:  false,
		}
		dj.queue.AddSong(song)
		Verbose(song.Submitter() + " added track " + song.Title())

		return song, nil
	}
	return nil, errors.New(VIDEO_TOO_LONG_MSG)
}

// NewPlaylist gathers the metadata for a YouTube playlist and returns it.
func (yt YouTube) NewPlaylist(user *gumble.User, id string) (Playlist, error) {
	var apiResponse *jsonq.JsonQuery
	var err error
	// Retrieve title of playlist
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlists?part=snippet&id=%s&key=%s", id, os.Getenv("YOUTUBE_API_KEY"))
	if apiResponse, err = PerformGetRequest(url); err != nil {
		return nil, err
	}
	title, _ := apiResponse.String("items", "0", "snippet", "title")

	playlist := &YouTubePlaylist{
		id:    id,
		title: title,
	}

	// Retrieve items in playlist
	url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=50&playlistId=%s&key=%s",
		id, os.Getenv("YOUTUBE_API_KEY"))
	if apiResponse, err = PerformGetRequest(url); err != nil {
		return nil, err
	}
	numVideos, _ := apiResponse.Int("pageInfo", "totalResults")
	if numVideos > 50 {
		numVideos = 50
	}

	for i := 0; i < numVideos; i++ {
		index := strconv.Itoa(i)
		videoID, _ := apiResponse.String("items", index, "snippet", "resourceId", "videoId")
		yt.NewSong(user, videoID, "", playlist)
	}
	return playlist, nil
}

// PerformGetRequest does all the grunt work for a YouTube HTTPS GET request.
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
