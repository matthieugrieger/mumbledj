/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/playlist_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PlaylistTestSuite struct {
	Playlist Playlist
	suite.Suite
}

func (suite *PlaylistTestSuite) SetupTest() {
	suite.Playlist = Playlist{
		ID:        "id",
		Title:     "title",
		Submitter: "submitter",
		Service:   "service",
	}
}

func (suite *PlaylistTestSuite) TestGetID() {
	suite.Equal("id", suite.Playlist.GetID())
}

func (suite *PlaylistTestSuite) TestGetTitle() {
	suite.Equal("title", suite.Playlist.GetTitle())
}

func (suite *PlaylistTestSuite) TestGetSubmitter() {
	suite.Equal("submitter", suite.Playlist.GetSubmitter())
}

func (suite *PlaylistTestSuite) TestGetService() {
	suite.Equal("service", suite.Playlist.GetService())
}

func TestPlaylistTestSuite(t *testing.T) {
	suite.Run(t, new(PlaylistTestSuite))
}
