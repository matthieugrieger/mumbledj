/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/config_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) SetupSuite() {
	DJ = NewMumbleDJ()
}

func (suite *ConfigTestSuite) SetupTest() {
	viper.Reset()
}

func (suite *ConfigTestSuite) TestCheckForDuplicateAliasesWhenNoDuplicatesExist() {
	viper.Set("commands.add.aliases", []string{"add", "a"})
	viper.Set("commands.addnext.aliases", []string{"addnext", "an"})
	viper.Set("commands.skip.aliases", []string{"skip", "s"})
	viper.Set("commands.skipplaylist.aliases", []string{"skipplaylist", "sp"})

	err := CheckForDuplicateAliases()

	suite.Nil(err, "No error should be returned as there are no duplicate aliases.")
}

func (suite *ConfigTestSuite) TestCheckForDuplicateAliasesWhenDuplicatesExist() {
	viper.Set("commands.add.aliases", []string{"add", "a"})
	viper.Set("commands.addnext.aliases", []string{"addnext", "an"})
	viper.Set("commands.skip.aliases", []string{"skip", "s"})
	viper.Set("commands.skipplaylist.aliases", []string{"skipplaylist", "sp"})
	viper.Set("commands.version.aliases", []string{"version", "v"})
	viper.Set("commands.volume.aliases", []string{"volume", "vol", "v"})

	err := CheckForDuplicateAliases()

	suite.NotNil(err, "An error should be returned as there are duplicate aliases.")
	suite.Contains(err.Error(), "v", "The error message should contain the duplicate alias.")
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
