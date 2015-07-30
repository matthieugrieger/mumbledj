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
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/jsonq"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
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

// ------
// TYPES
// ------

// YouTube implements the Service interface
type YouTube struct{}

// YouTubeSong holds the metadata for a song extracted from a YouTube video.
type YouTubeSong struct {
	submitter string
	title     string
	id        string
	offset    int
	filename  string
	duration  string
	thumbnail string
	skippers  []string
	playlist  Playlist
	dontSkip  bool
}

// YouTubePlaylist holds the metadata for a YouTube playlist.
type YouTubePlaylist struct {
	id    string
	title string
}

// ---------------
// YOUTUBE SERVICE
// ---------------

// Name of the service
func (yt YouTube) ServiceName() string {
	return "Youtube"
}

// Checks to see if service will accept URL
func (yt YouTube) URLRegex(url string) bool {
	return RegexpFromURL(url, append(youtubeVideoPatterns, []string{youtubePlaylistPattern}...)) != nil
}

// Creates the requested song/playlist and adds to the queue
func (yt YouTube) NewRequest(user *gumble.User, url string) (string, error) {
	var shortURL, startOffset = "", ""
	if re, err := regexp.Compile(youtubePlaylistPattern); err == nil {
		if re.MatchString(url) {
			if dj.HasPermission(user.Name, dj.conf.Permissions.AdminAddPlaylists) {
				shortURL = re.FindStringSubmatch(url)[1]
				playlist, err := yt.NewPlaylist(user.Name, shortURL)
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
			song, err := yt.NewSong(user.Name, shortURL, startOffset, nil)
			return song.Title(), err
		}
	} else {
		return "", err
	}
}

// NewSong gathers the metadata for a song extracted from a YouTube video, and returns
// the song.
func (yt YouTube) NewSong(user, id, offset string, playlist *YouTubePlaylist) (*YouTubeSong, error) {
	var apiResponse *jsonq.JsonQuery
	var err error
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails&id=%s&key=%s",
		id, os.Getenv("YOUTUBE_API_KEY"))
	if apiResponse, err = yt.PerformGetRequest(url); err != nil {
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
			offset:    int((offsetDays * 86400) + (offsetHours * 3600) + (offsetMinutes * 60) + offsetSeconds),
			filename:  id + ".m4a",
			duration:  durationString,
			thumbnail: thumbnail,
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
func (yt YouTube) NewPlaylist(user, id string) (*YouTubePlaylist, error) {
	var apiResponse *jsonq.JsonQuery
	var err error
	// Retrieve title of playlist
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlists?part=snippet&id=%s&key=%s",
		id, os.Getenv("YOUTUBE_API_KEY"))
	if apiResponse, err = yt.PerformGetRequest(url); err != nil {
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
	if apiResponse, err = yt.PerformGetRequest(url); err != nil {
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

// ------------
// YOUTUBE SONG
// ------------

// Download downloads the song via youtube-dl if it does not already exist on disk.
// All downloaded songs are stored in ~/.mumbledj/songs and should be automatically cleaned.
func (s *YouTubeSong) Download() error {

	// Checks to see if song is already downloaded
	if _, err := os.Stat(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, s.Filename())); os.IsNotExist(err) {
		Verbose("Downloading " + s.Title())
		cmd := exec.Command("youtube-dl", "--output", fmt.Sprintf(`~/.mumbledj/songs/%s`, s.Filename()), "--format", "m4a", "--", s.ID())
		if err := cmd.Run(); err == nil {
			if dj.conf.Cache.Enabled {
				dj.cache.CheckMaximumDirectorySize()
			}
			Verbose(s.Title() + " downloaded")
			return nil
		}
		Verbose(s.Title() + " failed to download")
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
		if isNil(s.Playlist()) {
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
		Verbose("Now playing " + s.Title())

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
				Verbose("Deleted " + s.Title())
				return nil
			}
			Verbose("Failed to delete " + s.Title())
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
	}
	return false
}

// Submitter returns the name of the submitter of the YouTubeSong.
func (s *YouTubeSong) Submitter() string {
	return s.submitter
}

// Title returns the title of the YouTubeSong.
func (s *YouTubeSong) Title() string {
	return s.title
}

// ID returns the id of the YouTubeSong.
func (s *YouTubeSong) ID() string {
	return s.id
}

// Filename returns the filename of the YouTubeSong.
func (s *YouTubeSong) Filename() string {
	return s.filename
}

// Duration returns the duration of the YouTubeSong.
func (s *YouTubeSong) Duration() string {
	return s.duration
}

// Thumbnail returns the thumbnail URL for the YouTubeSong.
func (s *YouTubeSong) Thumbnail() string {
	return s.thumbnail
}

// Playlist returns the playlist type for the YouTubeSong (may be nil).
func (s *YouTubeSong) Playlist() Playlist {
	return s.playlist
}

// DontSkip returns the DontSkip boolean value for the YouTubeSong.
func (s *YouTubeSong) DontSkip() bool {
	return s.dontSkip
}

// SetDontSkip sets the DontSkip boolean value for the YouTubeSong.
func (s *YouTubeSong) SetDontSkip(value bool) {
	s.dontSkip = value
}

// ----------------
// YOUTUBE PLAYLIST
// ----------------

// AddSkip adds a skip to the playlist's skippers slice.
func (p *YouTubePlaylist) AddSkip(username string) error {
	for _, user := range dj.playlistSkips[p.ID()] {
		if username == user {
			return errors.New("This user has already skipped the current song.")
		}
	}
	dj.playlistSkips[p.ID()] = append(dj.playlistSkips[p.ID()], username)
	return nil
}

// RemoveSkip removes a skip from the playlist's skippers slice. If username is not in the slice
// an error is returned.
func (p *YouTubePlaylist) RemoveSkip(username string) error {
	for i, user := range dj.playlistSkips[p.ID()] {
		if username == user {
			dj.playlistSkips[p.ID()] = append(dj.playlistSkips[p.ID()][:i], dj.playlistSkips[p.ID()][i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

// DeleteSkippers removes the skippers entry in dj.playlistSkips.
func (p *YouTubePlaylist) DeleteSkippers() {
	delete(dj.playlistSkips, p.ID())
}

// SkipReached calculates the current skip ratio based on the number of users within MumbleDJ's
// channel and the number of usernames in the skippers slice. If the value is greater than or equal
// to the skip ratio defined in the config, the function returns true, and returns false otherwise.
func (p *YouTubePlaylist) SkipReached(channelUsers int) bool {
	if float32(len(dj.playlistSkips[p.ID()]))/float32(channelUsers) >= dj.conf.General.PlaylistSkipRatio {
		return true
	}
	return false
}

// ID returns the id of the YouTubePlaylist.
func (p *YouTubePlaylist) ID() string {
	return p.id
}

// Title returns the title of the YouTubePlaylist.
func (p *YouTubePlaylist) Title() string {
	return p.title
}

// -----------
// YOUTUBE API
// -----------

// PerformGetRequest does all the grunt work for a YouTube HTTPS GET request.
func (yt YouTube) PerformGetRequest(url string) (*jsonq.JsonQuery, error) {
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
