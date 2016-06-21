/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/listtracks_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"testing"

	"github.com/layeh/gumble/gumbleffmpeg"
	"github.com/matthieugrieger/mumbledj/bot"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ListTracksCommandTestSuite struct {
	Command ListTracksCommand
	suite.Suite
}

func (suite *ListTracksCommandTestSuite) SetupSuite() {
	DJ = bot.NewMumbleDJ()
	bot.DJ = DJ

	// Trick the tests into thinking audio is already playing to avoid
	// attempting to play tracks that don't exist.
	DJ.AudioStream = new(gumbleffmpeg.Stream)

	viper.Set("commands.listtracks.aliases", []string{"listtracks", "list"})
	viper.Set("commands.listtracks.description", "listtracks")
	viper.Set("commands.listtracks.is_admin", false)
}

func (suite *ListTracksCommandTestSuite) SetupTest() {
	DJ.Queue = bot.NewQueue()
}

func (suite *ListTracksCommandTestSuite) TestAliases() {
	suite.Equal([]string{"listtracks", "list"}, suite.Command.Aliases())
}

func (suite *ListTracksCommandTestSuite) TestDescription() {
	suite.Equal("listtracks", suite.Command.Description())
}

func (suite *ListTracksCommandTestSuite) TestIsAdminCommand() {
	suite.False(suite.Command.IsAdminCommand())
}

func (suite *ListTracksCommandTestSuite) TestExecuteWithNoTracks() {
	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.Equal("", message, "No message should be returned.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.NotNil(err, "An error should be returned as there are no tracks to list.")
}

func (suite *ListTracksCommandTestSuite) TestExecuteWithNoArg() {
	track := new(bot.Track)
	track.Title = "title"
	track.Submitter = "test"

	DJ.Queue.AppendTrack(track)

	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.NotEqual("", message, "A message containing track information should be returned.")
	suite.Contains(message, "title", "The returned message should contain the track title.")
	suite.Contains(message, "test", "The returned message should contain the track submitter.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.Nil(err, "No error should be returned.")
}

func (suite *ListTracksCommandTestSuite) TestExecuteWithValidArg() {
	track1 := new(bot.Track)
	track1.Title = "first"
	track1.Submitter = "test"

	track2 := new(bot.Track)
	track2.Title = "second"
	track2.Submitter = "test"

	track3 := new(bot.Track)
	track3.Title = "third"
	track3.Submitter = "test"

	DJ.Queue.AppendTrack(track1)
	DJ.Queue.AppendTrack(track2)
	DJ.Queue.AppendTrack(track3)

	message, isPrivateMessage, err := suite.Command.Execute(nil, "2")

	suite.NotEqual("", message, "A message containing track information should be returned.")
	suite.Contains(message, "first", "The returned message should contain the first track.")
	suite.Contains(message, "second", "The returned message should contain the second track.")
	suite.NotContains(message, "third", "The returned message should not contain the third track.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.Nil(err, "No error should be returned.")
}

func (suite *ListTracksCommandTestSuite) TestExecuteWithArgLargerThanQueueLength() {
	track := new(bot.Track)
	track.Title = "track"
	track.Submitter = "test"

	DJ.Queue.AppendTrack(track)

	message, isPrivateMessage, err := suite.Command.Execute(nil, "2")

	suite.NotEqual("", message, "A message containing track information should be returned.")
	suite.Contains(message, "1", "The returned message should contain the first track.")
	suite.NotContains(message, "2", "The returned message should not contain any further tracks.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.Nil(err, "No error should be returned.")
}

func (suite *ListTracksCommandTestSuite) TestExecuteWithInvalidArg() {
	track := new(bot.Track)
	track.Title = "track"
	track.Submitter = "test"

	DJ.Queue.AppendTrack(track)

	message, isPrivateMessage, err := suite.Command.Execute(nil, "test")

	suite.Equal("", message, "No message should be returned.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.NotNil(err, "An error should be returned due to an invalid argument being supplied.")
}

func TestListTracksCommandTestSuite(t *testing.T) {
	suite.Run(t, new(ListTracksCommandTestSuite))
}
