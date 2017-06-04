/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/mixcloud.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package services

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/antonholmquist/jason"
	"layeh.com/gumble/gumble"
	"github.com/RichardNysater/mumbledj/bot"
	"github.com/RichardNysater/mumbledj/interfaces"
	"sync"
)

// Mixcloud is a wrapper around the Mixcloud API.
// https://www.mixcloud.com/developers/
type Mixcloud struct {
	*GenericService
}

// NewMixcloudService returns an initialized Mixcloud service object.
func NewMixcloudService() *Mixcloud {
	return &Mixcloud{
		&GenericService{
			ReadableName: "Mixcloud",
			Format:       "m4a",
			TrackRegex: []*regexp.Regexp{
				regexp.MustCompile(`https?:\/\/(www\.)?mixcloud\.com\/([\w-]+)\/([\w-]+)(#t=\n\n?(:\n\n)*)?`),
			},
			// Playlists are currently unsupported by Mixcloud's API.
			PlaylistRegex: nil,
		},
	}
}

// CheckAPIKey performs a test API call with the API key
// provided in the configuration file to determine if the
// service should be enabled.
func (mc *Mixcloud) CheckAPIKey() error {
	// Mixcloud (at the moment) does not require an API key,
	// so we can just return nil.
	return nil
}

// GetTracks uses the passed URL to find and return
// tracks associated with the URL. An error is returned
// if any error occurs during the API call.
func (mc *Mixcloud) GetTracks(url string, submitter *gumble.User) ([]interfaces.Track, error) {
	var (
		apiURL string
		err    error
		resp   *http.Response
		v      *jason.Object
		tracks []interfaces.Track
	)

	apiURL = strings.Replace(url, "www", "api", 1)

	// Track playback offset is not present in Mixcloud URLs,
	// so we can safely assume that users will not request
	// a playback offset in the URL.
	offset, _ := time.ParseDuration("0s")

	resp, err = http.Get(apiURL)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	v, err = jason.NewObjectFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	id, _ := v.GetString("slug")
	trackURL, _ := v.GetString("url")
	title, _ := v.GetString("name")
	author, _ := v.GetString("user", "username")
	authorURL, _ := v.GetString("user", "url")
	durationSecs, _ := v.GetInt64("audio_length")
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", durationSecs))
	thumbnail, err := v.GetString("pictures", "large")
	if err != nil {
		// Track has no artwork, using profile avatar instead.
		thumbnail, _ = v.GetString("user", "pictures", "large")
	}
	var wg sync.WaitGroup
	track := bot.Track{
		ID:             id,
		URL:            trackURL,
		Title:          title,
		Author:         author,
		AuthorURL:      authorURL,
		Submitter:      submitter.Name,
		Service:        mc.ReadableName,
		ThumbnailURL:   thumbnail,
		Filename:       id + ".track",
		Duration:       duration,
		PlaybackOffset: offset,
		Playlist:       nil,
		WaitGroup:	wg,
	}

	tracks = append(tracks, track)

	return tracks, nil
}
