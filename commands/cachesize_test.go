/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/cachesize_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"testing"

	"reik.pl/mumbledj/bot"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type CacheSizeCommandTestSuite struct {
	Command CacheSizeCommand
	suite.Suite
}

func (suite *CacheSizeCommandTestSuite) SetupSuite() {
	DJ = bot.NewMumbleDJ()

	viper.Set("commands.cachesize.aliases", []string{"cachesize", "cs"})
	viper.Set("commands.cachesize.description", "cachesize")
	viper.Set("commands.cachesize.is_admin", true)
}

func (suite *CacheSizeCommandTestSuite) TestAliases() {
	suite.Equal([]string{"cachesize", "cs"}, suite.Command.Aliases())
}

func (suite *CacheSizeCommandTestSuite) TestDescription() {
	suite.Equal("cachesize", suite.Command.Description())
}

func (suite *CacheSizeCommandTestSuite) TestIsAdminCommand() {
	suite.True(suite.Command.IsAdminCommand())
}

func (suite *CacheSizeCommandTestSuite) TestExecuteWhenCachingIsDisabled() {
	viper.Set("cache.enabled", false)

	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.Equal("", message, "An error occurred so no message should be returned.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.NotNil(err, "An error should be returned because caching is disabled.")
}

// TODO: Implement this test.
func (suite *CacheSizeCommandTestSuite) TestExecuteWhenCachingIsEnabled() {

}

func TestCacheSizeCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CacheSizeCommandTestSuite))
}
