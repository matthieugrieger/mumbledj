package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/layeh/gumble/gumble_ffmpeg"
)

type YouTubeDL struct {
	id        string
	title     string
	thumbnail string
	submitter string
	duration  string
	playlist  Playlist
	skippers  []string
	dontSkip  bool
}

// Download downloads the song via youtube-dl if it does not already exist on disk.
// All downloaded songs are stored in ~/.mumbledj/songs and should be automatically cleaned.
func (dl *YouTubeDL) Download() error {

	// Checks to see if song is already downloaded
	if _, err := os.Stat(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, dl.id+".m4a")); os.IsNotExist(err) {
		cmd := exec.Command("youtube-dl", "--output", fmt.Sprintf(`~/.mumbledj/songs/%s`, dl.id+".m4a"), "--format", "m4a", "--", dl.url)
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

// Play plays the song. Once the song is playing, a notification is displayed in a text message that features the song
// thumbnail, URL, title, duration, and submitter.
func (dl *YouTubeDL) Play() {
	if dl.offset != 0 {
		offsetDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", dl.offset))
		dj.audioStream.Offset = offsetDuration
	}
	dj.audioStream.Source = gumble_ffmpeg.SourceFile(fmt.Sprintf("%s/.mumbledj/songs/%s.m4a", dj.homeDir, dl.id))
	if err := dj.audioStream.Play(); err != nil {
		panic(err)
	} else {
		message := `<table><tr>	<td align="center"><img src="%s" width=150 /></td></tr><tr><td align="center"><b><a href="%s">%s</a> (%s)</b></td></tr><tr><td align="center">Added by %s</td></tr>`
		message = fmt.Sprintf(message, dl.thumbnail, dl.url, dl.title, dl.duration, dl.submitter)
		if isNil(dl.playlist) {
			dj.client.Self.Channel.Send(message+`</table>`, false)
		} else {
			message += `<tr><td align="center">From playlist "%s"</td></tr></table>`
			dj.client.Self.Channel.Send(fmt.Sprintf(message, dl.playlist.Title()), false)
		}
		Verbose("Now playing " + dl.title)

		go func() {
			dj.audioStream.Wait()
			dj.queue.OnSongFinished()
		}()
	}
}

// Delete deletes the song from ~/.mumbledj/songs if the cache is disabled.
func (dl *YouTubeDL) Delete() error {
	if dj.conf.Cache.Enabled == false {
		filePath := fmt.Sprintf("%s/.mumbledj/songs/%s.m4a", dj.homeDir, dl.id)
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
func (dl *YouTubeDL) AddSkip(username string) error {
	for _, user := range dl.skippers {
		if username == user {
			return errors.New("This user has already skipped the current song.")
		}
	}
	dl.skippers = append(dl.skippers, username)
	return nil
}

// RemoveSkip removes a skip from the skippers slice. If username is not in slice, an error is
// returned.
func (dl *YouTubeDL) RemoveSkip(username string) error {
	for i, user := range dl.skippers {
		if username == user {
			dl.skippers = append(s.skippers[:i], s.skippers[i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

// SkipReached calculates the current skip ratio based on the number of users within MumbleDJ's
// channel and the number of usernames in the skippers slice. If the value is greater than or equal
// to the skip ratio defined in the config, the function returns true, and returns false otherwise.
func (dl *YouTubeDL) SkipReached(channelUsers int) bool {
	if float32(len(dl.skippers))/float32(channelUsers) >= dj.conf.General.SkipRatio {
		return true
	}
	return false
}

// Submitter returns the name of the submitter of the Song.
func (dl *YouTubeDL) Submitter() string {
	return dl.submitter
}

// Title returns the title of the Song.
func (dl *YouTubeDL) Title() string {
	return dl.title
}

// ID returns the id of the Song.
func (dl *YouTubeDL) ID() string {
	return dl.id
}

// Filename returns the filename of the Song.
func (dl *YouTubeDL) Filename() string {
	return dl.id + ".m4a"
}

// Duration returns the duration of the Song.
func (dl *YouTubeDL) Duration() string {
	return dl.duration
}

// Thumbnail returns the thumbnail URL for the Song.
func (dl *YouTubeDL) Thumbnail() string {
	return dl.thumbnail
}

// Playlist returns the playlist type for the Song (may be nil).
func (dl *YouTubeDL) Playlist() Playlist {
	return dl.playlist
}

// DontSkip returns the DontSkip boolean value for the Song.
func (dl *YouTubeDL) DontSkip() bool {
	return dl.dontSkip
}

// SetDontSkip sets the DontSkip boolean value for the Song.
func (dl *YouTubeDL) SetDontSkip(value bool) {
	dl.dontSkip = value
}
