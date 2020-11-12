/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/skiptracker_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"testing"

	"layeh.com/gumble/gumble"
	"github.com/stretchr/testify/suite"
)

type SkipTrackerTestSuite struct {
	suite.Suite
	Skips *SkipTracker
	User1 *gumble.User
	User2 *gumble.User
}

func (suite *SkipTrackerTestSuite) SetupSuite() {
	suite.User1 = new(gumble.User)
	suite.User1.Name = "User1"
	suite.User2 = new(gumble.User)
	suite.User2.Name = "User2"
}

func (suite *SkipTrackerTestSuite) SetupTest() {
	suite.Skips = NewSkipTracker()
}

func (suite *SkipTrackerTestSuite) TestNewSkipTracker() {
	suite.Zero(suite.Skips.NumTrackSkips(), "The track skip slice should be empty upon initialization.")
	suite.Zero(suite.Skips.NumPlaylistSkips(), "The playlist skip slice should be empty upon initialization.")
}

// TODO: Fix these tests.
/*func (suite *SkipTrackerTestSuite) TestAddTrackSkip() {
	err := suite.Skips.AddTrackSkip(suite.User1)

	suite.Equal(1, suite.Skips.NumTrackSkips(), "There should now be one user in the track skip slice.")
	suite.Zero(0, suite.Skips.NumPlaylistSkips(), "The playlist skip slice should be unaffected.")
	suite.Nil(err, "No error should be returned.")

	err = suite.Skips.AddTrackSkip(suite.User2)

	suite.Equal(2, suite.Skips.NumTrackSkips(), "There should now be two users in the track skip slice.")
	suite.Zero(0, suite.Skips.NumPlaylistSkips(), "The playlist skip slice should be unaffected.")
	suite.Nil(err, "No error should be returned.")

	err = suite.Skips.AddTrackSkip(suite.User1)

	suite.Equal(2, suite.Skips.NumTrackSkips(), "This is a duplicate skip, so the track skip slice should still only have two users.")
	suite.Zero(0, suite.Skips.NumPlaylistSkips(), "The playlist skip slice should be unaffected.")
	suite.NotNil(err, "An error should be returned since this user has already voted to skip the current track.")
}

func (suite *SkipTrackerTestSuite) TestAddPlaylistSkip() {
	err := suite.Skips.AddPlaylistSkip(suite.User1)

	suite.Zero(suite.Skips.NumTrackSkips(), "The track skip slice should be unaffected.")
	suite.Equal(1, suite.Skips.NumPlaylistSkips(), "There should now be one user in the playlist skip slice.")
	suite.Nil(err, "No error should be returned.")

	err = suite.Skips.AddPlaylistSkip(suite.User2)

	suite.Zero(suite.Skips.NumTrackSkips(), "The track skip slice should be unaffected.")
	suite.Equal(2, suite.Skips.NumPlaylistSkips(), "There should now be two users in the playlist skip slice.")
	suite.Nil(err, "No error should be returned.")

	err = suite.Skips.AddPlaylistSkip(suite.User1)

	suite.Zero(suite.Skips.NumTrackSkips(), "The track skip slice should be unaffected.")
	suite.Equal(2, suite.Skips.NumPlaylistSkips(), "This is a duplicate skip, so the playlist skip slice should still only have two users.")
	suite.NotNil(err, "An error should be returned since this user has already voted to skip the current playlist.")
}

func (suite *SkipTrackerTestSuite) TestRemoveTrackSkip() {
	suite.Skips.AddTrackSkip(suite.User1)
	err := suite.Skips.RemoveTrackSkip(suite.User2)

	suite.Equal(1, suite.Skips.NumTrackSkips(), "User2 has not skipped the track so the track skip slice should be unaffected.")
	suite.Zero(suite.Skips.NumPlaylistSkips(), "The playlist skip slice should be unaffected.")
	suite.NotNil(err, "An error should be returned since User2 has not skipped the track yet.")

	err = suite.Skips.RemoveTrackSkip(suite.User1)

	suite.Zero(suite.Skips.NumTrackSkips(), "User1 skipped the track, so their skip should be removed.")
	suite.Zero(suite.Skips.NumPlaylistSkips(), "The playlist skip slice should be unaffected.")
	suite.Nil(err, "No error should be returned.")
}

func (suite *SkipTrackerTestSuite) TestRemovePlaylistSkip() {
	suite.Skips.AddPlaylistSkip(suite.User1)
	err := suite.Skips.RemovePlaylistSkip(suite.User2)

	suite.Zero(suite.Skips.NumTrackSkips(), "The track skip slice should be unaffected.")
	suite.Equal(1, suite.Skips.NumPlaylistSkips(), "User2 has not skipped the playlist so the playlist skip slice should be unaffected.")
	suite.NotNil(err, "An error should be returned since User2 has not skipped the playlist yet.")

	err = suite.Skips.RemovePlaylistSkip(suite.User1)

	suite.Zero(suite.Skips.NumTrackSkips(), "The track skip slice should be unaffected.")
	suite.Zero(suite.Skips.NumPlaylistSkips(), "User1 skipped the playlist, so their skip should be removed.")
	suite.Nil(err, "No error should be returned.")
}

func (suite *SkipTrackerTestSuite) TestResetTrackSkips() {
	suite.Skips.AddTrackSkip(suite.User1)
	suite.Skips.AddTrackSkip(suite.User2)
	suite.Skips.AddPlaylistSkip(suite.User1)
	suite.Skips.AddPlaylistSkip(suite.User2)

	suite.Equal(2, suite.Skips.NumTrackSkips(), "There should be two users in the track skip slice.")
	suite.Equal(2, suite.Skips.NumPlaylistSkips(), "There should be two users in the playlist skip slice.")

	suite.Skips.ResetTrackSkips()

	suite.Zero(suite.Skips.NumTrackSkips(), "The track skip slice has been reset, so the length should be zero.")
	suite.Equal(2, suite.Skips.NumPlaylistSkips(), "The playlist skip slice should be unaffected.")
}

func (suite *SkipTrackerTestSuite) TestResetPlaylistSkips() {
	suite.Skips.AddTrackSkip(suite.User1)
	suite.Skips.AddTrackSkip(suite.User2)
	suite.Skips.AddPlaylistSkip(suite.User1)
	suite.Skips.AddPlaylistSkip(suite.User2)

	suite.Equal(2, suite.Skips.NumTrackSkips(), "There should be two users in the track skip slice.")
	suite.Equal(2, suite.Skips.NumPlaylistSkips(), "There should be two users in the playlist skip slice.")

	suite.Skips.ResetPlaylistSkips()

	suite.Equal(2, suite.Skips.NumTrackSkips(), "The track skip slice should be unaffected.")
	suite.Zero(suite.Skips.NumPlaylistSkips(), "The playlist skip slice has been reset, so the length should be zero.")
}*/

func TestSkipTrackerTestSuite(t *testing.T) {
	suite.Run(t, new(SkipTrackerTestSuite))
}
