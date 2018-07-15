/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/toggleshuffle_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"testing"

	"reik.pl/mumbledj/bot"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ToggleShuffleCommandTestSuite struct {
	Command ToggleShuffleCommand
	suite.Suite
}

func (suite *ToggleShuffleCommandTestSuite) SetupSuite() {
	DJ = bot.NewMumbleDJ()

	viper.Set("commands.toggleshuffle.aliases", []string{"toggleshuffle", "ts"})
	viper.Set("commands.toggleshuffle.description", "toggleshuffle")
	viper.Set("commands.toggleshuffle.is_admin", true)
}

func (suite *ToggleShuffleCommandTestSuite) TestAliases() {
	suite.Equal([]string{"toggleshuffle", "ts"}, suite.Command.Aliases())
}

func (suite *ToggleShuffleCommandTestSuite) TestDescription() {
	suite.Equal("toggleshuffle", suite.Command.Description())
}

func (suite *ToggleShuffleCommandTestSuite) TestIsAdminCommand() {
	suite.True(suite.Command.IsAdminCommand())
}

func (suite *ToggleShuffleCommandTestSuite) TestExecuteWhenShuffleIsOff() {
	viper.Set("queue.automatic_shuffle_on", false)

	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.NotEqual("", message, "A message should be returned.")
	suite.False(isPrivateMessage, "This should not be a private message.")
	suite.Nil(err, "No error should be returned.")
	suite.True(viper.GetBool("queue.automatic_shuffle_on"), "Automatic shuffling should now be on.")
}

func (suite *ToggleShuffleCommandTestSuite) TestExecuteWhenShuffleIsOn() {
	viper.Set("queue.automatic_shuffle_on", true)

	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.NotEqual("", message, "A message should be returned.")
	suite.False(isPrivateMessage, "This should not be a private message.")
	suite.Nil(err, "No error should be returned.")
	suite.False(viper.GetBool("queue.automatic_shuffle_on"), "Automatic shuffling should now be off.")
}

func TestToggleShuffleCommandTestSuite(t *testing.T) {
	suite.Run(t, new(ToggleShuffleCommandTestSuite))
}
