/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/add_test.go
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

type AddCommandTestSuite struct {
	Command AddCommand
	suite.Suite
}

func (suite *AddCommandTestSuite) SetupSuite() {
	DJ = bot.NewMumbleDJ()
	bot.DJ = DJ

	// Trick the tests into thinking audio is already playing to avoid
	// attempting to play tracks that don't exist.
	DJ.AudioStream = new(gumbleffmpeg.Stream)

	viper.Set("commands.add.aliases", []string{"add", "a"})
	viper.Set("commands.add.description", "add")
	viper.Set("commands.add.is_admin", false)
}

func (suite *AddCommandTestSuite) SetupTest() {
	DJ.Queue = bot.NewQueue()
}

func (suite *AddCommandTestSuite) TestAliases() {
	suite.Equal([]string{"add", "a"}, suite.Command.Aliases())
}

func (suite *AddCommandTestSuite) TestDescription() {
	suite.Equal("add", suite.Command.Description())
}

func (suite *AddCommandTestSuite) TestIsAdminCommand() {
	suite.False(suite.Command.IsAdminCommand())
}

func (suite *AddCommandTestSuite) TestExecuteWithNoArgs() {
	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.Equal("", message, "No message should be returned since an error occurred.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.NotNil(err, "An error should be returned for attempting to add a track without providing a URL.")
}

// TODO: Implement this test.
func (suite *AddCommandTestSuite) TestExecuteWhenNoTracksFound() {

}

// TODO: Implement this test.
func (suite *AddCommandTestSuite) TestExecuteWhenTrackFound() {

}

// TODO: Implement this test.
func (suite *AddCommandTestSuite) TestExecuteWhenMultipleTracksFound() {

}

// TODO: Implement this test.
func (suite *AddCommandTestSuite) TestExecuteWithMultipleURLs() {

}

func TestAddCommandTestSuite(t *testing.T) {
	suite.Run(t, new(AddCommandTestSuite))
}
