/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/track_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TrackTestSuite struct {
	suite.Suite
	Track Track
}

func (suite *TrackTestSuite) SetupTest() {
	duration, _ := time.ParseDuration("1s")
	suite.Track = Track{
		ID:           "id",
		Title:        "title",
		Author:       "author",
		Submitter:    "submitter",
		Service:      "service",
		Filename:     "filename",
		ThumbnailURL: "thumbnailurl",
		Duration:     duration,
		Playlist:     new(Playlist),
	}
}

func (suite *TrackTestSuite) TestGetID() {
	suite.Equal("id", suite.Track.GetID())
}

func (suite *TrackTestSuite) TestGetTitle() {
	suite.Equal("title", suite.Track.GetTitle())
}

func (suite *TrackTestSuite) TestGetAuthor() {
	suite.Equal("author", suite.Track.GetAuthor())
}

func (suite *TrackTestSuite) TestGetSubmitter() {
	suite.Equal("submitter", suite.Track.GetSubmitter())
}

func (suite *TrackTestSuite) TestGetService() {
	suite.Equal("service", suite.Track.GetService())
}

func (suite *TrackTestSuite) TestGetFilenameWhenExists() {
	result := suite.Track.GetFilename()

	suite.Equal("filename", result)
}

func (suite *TrackTestSuite) TestGetFilenameWhenNotExists() {
	suite.Track.Filename = ""

	result := suite.Track.GetFilename()

	suite.Equal("", result)
}

func (suite *TrackTestSuite) TestGetThumbnailURLWhenExists() {
	result := suite.Track.GetThumbnailURL()

	suite.Equal("thumbnailurl", result)
}

func (suite *TrackTestSuite) TestGetThumbnailURLWhenNotExists() {
	suite.Track.ThumbnailURL = ""

	result := suite.Track.GetThumbnailURL()

	suite.Equal("", result)
}

func (suite *TrackTestSuite) TestGetDuration() {
	duration, _ := time.ParseDuration("1s")

	suite.Equal(duration, suite.Track.GetDuration())
}

func (suite *TrackTestSuite) TestGetPlaylistWhenExists() {
	result := suite.Track.GetPlaylist()

	suite.NotNil(result)
}

func (suite *TrackTestSuite) TestGetPlaylistWhenNotExists() {
	suite.Track.Playlist = nil

	result := suite.Track.GetPlaylist()

	suite.Nil(result)
}

func TestTrackTestSuite(t *testing.T) {
	suite.Run(t, new(TrackTestSuite))
}
