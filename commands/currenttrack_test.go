/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/currenttrack_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"testing"

	"github.com/layeh/gumble/gumbleffmpeg"
	"reik.pl/mumbledj/bot"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type CurrentTrackCommandTestSuite struct {
	Command CurrentTrackCommand
	suite.Suite
}

func (suite *CurrentTrackCommandTestSuite) SetupSuite() {
	DJ = bot.NewMumbleDJ()

	// Trick the tests into thinking audio is already playing to avoid
	// attempting to play tracks that don't exist.
	DJ.AudioStream = new(gumbleffmpeg.Stream)

	viper.Set("commands.currenttrack.aliases", []string{"currenttrack", "current"})
	viper.Set("commands.currenttrack.description", "currenttrack")
	viper.Set("commands.currenttrack.is_admin", false)
}

func (suite *CurrentTrackCommandTestSuite) SetupTest() {
	DJ.Queue = bot.NewQueue()
}

func (suite *CurrentTrackCommandTestSuite) TestAliases() {
	suite.Equal([]string{"currenttrack", "current"}, suite.Command.Aliases())
}

func (suite *CurrentTrackCommandTestSuite) TestDescription() {
	suite.Equal("currenttrack", suite.Command.Description())
}

func (suite *CurrentTrackCommandTestSuite) TestIsAdminCommand() {
	suite.False(suite.Command.IsAdminCommand())
}

func (suite *CurrentTrackCommandTestSuite) TestExecuteWhenQueueIsEmpty() {
	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.Equal("", message, "No message should be returned since an error occurred.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.NotNil(err, "An error should be returned since the queue is empty.")
}

func (suite *CurrentTrackCommandTestSuite) TestExecuteWhenQueueNotEmpty() {
	track := new(bot.Track)
	track.Submitter = "test"
	track.Title = "test"

	DJ.Queue.AppendTrack(track)

	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.NotEqual("", message, "A message should be returned with the current track information.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.Nil(err, "No error should be returned.")
}

func TestCurrentTrackCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CurrentTrackCommandTestSuite))
}
