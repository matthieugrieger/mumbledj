/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/help_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"testing"

	"github.com/matthieugrieger/mumbledj/bot"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type HelpCommandTestSuite struct {
	Command HelpCommand
	suite.Suite
}

func (suite *HelpCommandTestSuite) SetupSuite() {
	DJ = bot.NewMumbleDJ()

	viper.Set("commands.help.aliases", []string{"help", "h"})
	viper.Set("commands.help.description", "help")
	viper.Set("commands.help.is_admin", false)
}

func (suite *HelpCommandTestSuite) SetupTest() {
	viper.Set("admins.enabled", true)
}

func (suite *HelpCommandTestSuite) TestAliases() {
	suite.Equal([]string{"help", "h"}, suite.Command.Aliases())
}

func (suite *HelpCommandTestSuite) TestDescription() {
	suite.Equal("help", suite.Command.Description())
}

func (suite *HelpCommandTestSuite) TestIsAdminCommand() {
	suite.False(suite.Command.IsAdminCommand())
}

// TODO: Implement this test.
func (suite *HelpCommandTestSuite) TestExecuteWhenPermissionsEnabledAndUserIsNotAdmin() {

}

// TODO: Implement this test.
func (suite *HelpCommandTestSuite) TestExecuteWhenPermissionsEnabledAndUserIsAdmin() {

}

func (suite *HelpCommandTestSuite) TestExecuteWhenPermissionsDisabled() {
	viper.Set("admins.enabled", false)

	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.NotEqual("", message, "A message should be returned.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.Nil(err, "No error should be returned.")
	suite.Contains(message, "help", "The returned message should contain command descriptions.")
	suite.Contains(message, "add", "The returned message should contain command descriptions.")
	suite.Contains(message, "Admin Commands", "The returned message should contain admin command descriptions.")
}

func TestHelpCommandTestSuite(t *testing.T) {
	suite.Run(t, new(HelpCommandTestSuite))
}
