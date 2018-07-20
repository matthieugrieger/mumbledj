/*
 * MumbleDJ
 * By Matthieu Grieger
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 * commands/reset_test.go
 */

package commands

import (
	"testing"

	"github.com/layeh/gumble/gumbleffmpeg"
	"github.com/matthieugrieger/mumbledj/bot"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ResetCommandTestSuite struct {
	Command ResetCommand
	suite.Suite
}

func (suite *ResetCommandTestSuite) SetupSuite() {
	DJ = bot.NewMumbleDJ()
	bot.DJ = DJ

	// Trick the tests into thinking audio is already playing to avoid
	// attempting to play tracks that don't exist.
	DJ.AudioStream = new(gumbleffmpeg.Stream)

	viper.Set("commands.reset.aliases", []string{"reset", "re"})
	viper.Set("commands.reset.description", "Resets the queue by removing all queue items.")
	viper.Set("commands.reset.is_admin", true)
}

func (suite *ResetCommandTestSuite) SetupTest() {
	DJ.Queue = bot.NewQueue()
}

func (suite *ResetCommandTestSuite) TestAliases() {
	suite.Equal([]string{"reset", "re"}, suite.Command.Aliases())
}

func (suite *ResetCommandTestSuite) TestDescription() {
	suite.Equal("reset", suite.Command.Description())
}

func (suite *ResetCommandTestSuite) TestIsAdminCommand() {
	suite.False(suite.Command.IsAdminCommand())
}

func (suite *ResetCommandTestSuite) TestResetWorksOnEmpty() {
	// TODO: Assuming the Queue is currently empty, is that the case?
	suite.Command.Execute(nil)
	suite.Zero(DJ.Queue.Length())
}

func (suite *ResetCommandTestSuite) TestResetWorksOneTrack() {
	track := new(bot.Track)
	track.Submitter = "test"
	track.Title = "test"

	DJ.Queue.AppendTrack(track)
	suite.Command.Execute(nil)
	suite.Zero(DJ.Queue.Length())
}

func (suite *ResetCommandTestSuite) TestResetWorksMultipleTracks() {
	track1 := new(bot.Track)
	track1.Submitter = "test"
	track1.Title = "test"
	track2 := new(bot.Track)
	track2.Submitter = "test"
	track2.Title = "test"

	DJ.Queue.AppendTrack(track1)
	DJ.Queue.AppendTrack(track2)
	suite.Command.Execute(nil)
	suite.Zero(DJ.Queue.Length())
}

func TestResetCommandTestSuite(t *testing.T) {
	suite.Run(t, new(ResetCommandTestSuite))
}
