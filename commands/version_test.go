/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/version_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"testing"

	"reik.pl/mumbledj/bot"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type VersionCommandTestSuite struct {
	Command VersionCommand
	suite.Suite
}

func (suite *VersionCommandTestSuite) SetupSuite() {
	DJ = bot.NewMumbleDJ()

	viper.Set("commands.version.aliases", []string{"version", "v"})
	viper.Set("commands.version.description", "version")
	viper.Set("commands.version.is_admin", false)
	DJ.Version = "test"
}

func (suite *VersionCommandTestSuite) TestAliases() {
	suite.Equal([]string{"version", "v"}, suite.Command.Aliases())
}

func (suite *VersionCommandTestSuite) TestDescription() {
	suite.Equal("version", suite.Command.Description())
}

func (suite *VersionCommandTestSuite) TestIsAdminCommand() {
	suite.False(suite.Command.IsAdminCommand())
}

func (suite *VersionCommandTestSuite) TestExecute() {
	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.NotEqual("", message, "A message should be returned.")
	suite.Contains(message, "MumbleDJ", "The message should contain a MumbleDJ version string.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.Nil(err, "No error should be returned.")
}

func TestVersionCommandTestSuite(t *testing.T) {
	suite.Run(t, new(VersionCommandTestSuite))
}
