/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/volume_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"fmt"
	"testing"

	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleffmpeg"
	"github.com/RichardNysater/mumbledj/bot"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type VolumeCommandTestSuite struct {
	Command VolumeCommand
	suite.Suite
}

func (suite *VolumeCommandTestSuite) SetupSuite() {
	DJ = bot.NewMumbleDJ()

	viper.Set("commands.volume.aliases", []string{"volume", "vol"})
	viper.Set("commands.volume.description", "volume")
	viper.Set("commands.volume.is_admin", false)
	viper.Set("volume.lowest", 0.2)
	viper.Set("volume.highest", 1)
	viper.Set("volume.default", 0.4)
}

func (suite *VolumeCommandTestSuite) SetupTest() {
	DJ.Volume = float32(viper.GetFloat64("volume.default"))
}

func (suite *VolumeCommandTestSuite) TestAliases() {
	suite.Equal([]string{"volume", "vol"}, suite.Command.Aliases())
}

func (suite *VolumeCommandTestSuite) TestDescription() {
	suite.Equal("volume", suite.Command.Description())
}

func (suite *VolumeCommandTestSuite) TestIsAdminCommand() {
	suite.False(suite.Command.IsAdminCommand())
}

func (suite *VolumeCommandTestSuite) TestExecuteWithNoArgs() {
	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.NotEqual("", message, "A message should be returned.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.Nil(err, "No error should be returned.")
	suite.Contains(message, "0.4", "The returned string should contain the current volume.")
}

func (suite *VolumeCommandTestSuite) TestExecuteWithValidArg() {
	dummyUser := &gumble.User{
		Name: "test",
	}
	message, isPrivateMessage, err := suite.Command.Execute(dummyUser, "0.6")

	suite.NotEqual("", message, "A message should be returned.")
	suite.False(isPrivateMessage, "This should not be a private message.")
	suite.Nil(err, "No error should be returned.")
	suite.Contains(message, "0.6", "The returned string should contain the new volume.")
	suite.Contains(message, "test", "The returned string should contain the username of whomever changed the volume.")
}

func (suite *VolumeCommandTestSuite) TestExecuteWithValidArgAndNonNilStream() {
	dummyUser := &gumble.User{
		Name: "test",
	}
	DJ.AudioStream = new(gumbleffmpeg.Stream)
	DJ.AudioStream.Volume = 0.2

	message, isPrivateMessage, err := suite.Command.Execute(dummyUser, "0.6")

	suite.NotEqual("", message, "A message should be returned.")
	suite.False(isPrivateMessage, "This should not be a private message.")
	suite.Nil(err, "No error should be returned.")
	suite.Contains(message, "0.6", "The returned string should contain the new volume.")
	suite.Contains(message, "test", "The returned string should contain the username of whomever changed the volume.")
	suite.Equal("0.60", fmt.Sprintf("%.2f", DJ.AudioStream.Volume), "The audio stream value should match the new volume.")
}

func (suite *VolumeCommandTestSuite) TestExecuteWithArgOutOfRange() {
	message, isPrivateMessage, err := suite.Command.Execute(nil, "1.4")

	suite.Equal("", message, "No message should be returned as an error occurred.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.NotNil(err, "An error should be returned since the provided argument was outside of the valid range.")
}

func (suite *VolumeCommandTestSuite) TestExecuteWithInvalidArg() {
	message, isPrivateMessage, err := suite.Command.Execute(nil, "test")

	suite.Equal("", message, "No message should be returned as an error occurred.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.NotNil(err, "An error should be returned as a non-floating-point argument was provided.")
}

func TestVolumeCommandTestSuite(t *testing.T) {
	suite.Run(t, new(VolumeCommandTestSuite))
}
