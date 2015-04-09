/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/youtube/song.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package youtube

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Song holds the metadata for a song extracted from a YouTube video.
type Song struct {
	submitter string
	title     string
	id        string
	duration  string
	thumbnail string
	skippers  []string
	playlist  *Playlist
	dontSkip  bool
}

// NewSong gathers the metadata for a song extracted from a YouTube video, and returns
// the Song.
func NewSong(user, id string, playlist *Playlist) (*Song, error) {
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

	if dj.conf.General.MaxSongDuration == 0 || totalSeconds <= dj.conf.General.MaxSongDuration {
		song := &Song{
			submitter:       user,
			title:           title,
			id:              id,
			duration:        durationString,
			secondsDuration: totalSeconds,
			thumbnail:       thumbnail,
			skippers:        make([]string, 0),
			playlist:        nil,
			dontSkip:        false,
		}
		dj.queue.AddSong(song)
		return song, nil
	}
	return nil, errors.New("Song exceeds the maximum allowed duration.")
}

// Download downloads the song via youtube-dl if it does not already exist on disk.
// All downloaded songs are stored in ~/.mumbledj/songs and should be automatically cleaned.
func (s *Song) Download() error {
	if _, err := os.Stat(fmt.Sprintf("%s/.mumbledj/songs/%s.m4a", dj.homeDir, s.id)); os.IsNotExist(err) {
		cmd := exec.Command("youtube-dl", "--output", fmt.Sprintf(`~/.mumbledj/songs/%s.m4a`, s.id), "--format", "m4a", "--", s.id)
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
func (s *Song) Play() {
	if err := dj.audioStream.Play(fmt.Sprintf("%s/.mumbledj/songs/%s.m4a", dj.homeDir, s.id), dj.queue.OnSongFinished); err != nil {
		panic(err)
	} else {
		if s.playlist == nil {
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
			dj.client.Self.Channel.Send(fmt.Sprintf(message, s.thumbnail, s.id, s.title,
				s.duration, s.submitter), false)
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
			dj.client.Self.Channel.Send(fmt.Sprintf(message, s.thumbnail, s.id,
				s.title, s.duration, s.submitter, s.playlist.title), false)
		}
	}
}

// Delete deletes the song from ~/.mumbledj/songs if the cache is disabled.
func (s *Song) Delete() error {
	if dj.conf.Cache.Enabled == false {
		filePath := fmt.Sprintf("%s/.mumbledj/songs/%s.m4a", dj.homeDir, s.id)
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
func (s *Song) AddSkip(username string) error {
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
func (s *Song) RemoveSkip(username string) error {
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
func (s *Song) SkipReached(channelUsers int) bool {
	if float32(len(s.skippers))/float32(channelUsers) >= dj.conf.General.SkipRatio {
		return true
	}
	return false
}
