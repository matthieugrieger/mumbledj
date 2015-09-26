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
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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
			song, err := yt.NewSong(user.Name, shortURL, startOffset, nil)
			if !isNil(song) {
<<<<<<< HEAD
				return song.Title(), err
=======
				return append(songArray, song), nil
>>>>>>> abf98ad... Song conditions are now checked in Service
			} else {
				return nil, err
			}
		}
	} else {
		return nil, err
	}
}

// NewSong gathers the metadata for a song extracted from a YouTube video, and returns the song.
func (yt YouTube) NewSong(user *gumble.User, id, offset string, playlist Playlist) (Song, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails&id=%s&key=%s", id, os.Getenv("YOUTUBE_API_KEY"))
	if apiResponse, err := PerformGetRequest(url); err == nil {
		title, _ := apiResponse.String("items", "0", "snippet", "title")
		thumbnail, _ := apiResponse.String("items", "0", "snippet", "thumbnails", "high", "url")
		duration, _ := apiResponse.String("items", "0", "contentDetails", "duration")

		song := &YouTubeSong{
			submitter: user,
			title:     title,
			id:        id,
			url:       "https://youtu.be/" + id,
			offset:    int(yt.parseTime(offset, `\?T\=(?P<days>\d+D)?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`).Seconds()),
			duration:  int(yt.parseTime(duration, `P(?P<days>\d+D)?T(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`).Seconds()),
			thumbnail: thumbnail,
			format:    "m4a",
			skippers:  make([]string, 0),
			playlist:  playlist,
			dontSkip:  false,
			service:   yt,
		}

<<<<<<< HEAD
		return song, nil
=======
// Download downloads the song via youtube-dl if it does not already exist on disk.
// All downloaded songs are stored in ~/.mumbledj/songs and should be automatically cleaned.
func (s *YouTubeSong) Download() error {
	if _, err := os.Stat(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, s.Filename())); os.IsNotExist(err) {
		cmd := exec.Command("youtube-dl", "--no-mtime", "--output", fmt.Sprintf(`~/.mumbledj/songs/%s`, s.Filename()), "--format", "m4a", "--", s.ID())
		if err := cmd.Run(); err == nil {
			if dj.conf.Cache.Enabled {
				dj.cache.CheckMaximumDirectorySize()
			}
			return nil
		}
		return errors.New("Song download failed.")
	}
	return nil
}

// Play plays the song. Once the song is playing, a notification is displayed in a text message that features the video
// thumbnail, URL, title, duration, and submitter.
func (s *YouTubeSong) Play() {
	if s.offset != 0 {
		offsetDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", s.offset))
		dj.audioStream.Offset = offsetDuration
	}
	dj.audioStream.Source = gumble_ffmpeg.SourceFile(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, s.Filename()))
	if err := dj.audioStream.Play(); err != nil {
		panic(err)
	} else {
		if s.Playlist() == nil {
			message := `
				<table>
					<tr>
						<td align="center"><img src="%s" width=150 /></td>
					</tr>
					<tr>
						<td align="center"><b><a href="http://youtu.be/%s">%s</a> (%s)</b></td>
					</tr>
					<tr>
						<td align="center">Added by %s</td>
					</tr>
				</table>
			`
			dj.client.Self.Channel.Send(fmt.Sprintf(message, s.Thumbnail(), s.ID(), s.Title(),
				s.Duration(), s.Submitter()), false)
		} else {
			message := `
				<table>
					<tr>
						<td align="center"><img src="%s" width=150 /></td>
					</tr>
					<tr>
						<td align="center"><b><a href="http://youtu.be/%s">%s</a> (%s)</b></td>
					</tr>
					<tr>
						<td align="center">Added by %s</td>
					</tr>
					<tr>
						<td align="center">From playlist "%s"</td>
					</tr>
				</table>
			`
			dj.client.Self.Channel.Send(fmt.Sprintf(message, s.Thumbnail(), s.ID(),
				s.Title(), s.Duration(), s.Submitter(), s.Playlist().Title()), false)
		}
		go func() {
			dj.audioStream.Wait()
			dj.queue.OnSongFinished()
		}()
	}
}

// Delete deletes the song from ~/.mumbledj/songs if the cache is disabled.
func (s *YouTubeSong) Delete() error {
	if dj.conf.Cache.Enabled == false {
		filePath := fmt.Sprintf("%s/.mumbledj/songs/%s.m4a", dj.homeDir, s.ID())
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err == nil {
				return nil
			}
			return errors.New("Error occurred while deleting audio file.")
		}
		return nil
	}
	return nil
}

// AddSkip adds a skip to the skippers slice. If the user is already in the slice, AddSkip
// returns an error and does not add a duplicate skip.
func (s *YouTubeSong) AddSkip(username string) error {
	for _, user := range s.skippers {
		if username == user {
			return errors.New("This user has already skipped the current song.")
		}
	}
	s.skippers = append(s.skippers, username)
	return nil
}

// RemoveSkip removes a skip from the skippers slice. If username is not in slice, an error is
// returned.
func (s *YouTubeSong) RemoveSkip(username string) error {
	for i, user := range s.skippers {
		if username == user {
			s.skippers = append(s.skippers[:i], s.skippers[i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

// SkipReached calculates the current skip ratio based on the number of users within MumbleDJ's
// channel and the number of usernames in the skippers slice. If the value is greater than or equal
// to the skip ratio defined in the config, the function returns true, and returns false otherwise.
func (s *YouTubeSong) SkipReached(channelUsers int) bool {
	if float32(len(s.skippers))/float32(channelUsers) >= dj.conf.General.SkipRatio {
		return true
>>>>>>> 2df3613... Fixed cache clearing earlier than expected
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
		if song, err := yt.NewSong(user, videoID, "", playlist); err == nil {
			songArray = append(songArray, song)
		}
	}
	return songArray, nil
}
