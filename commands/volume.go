/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/volume.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"
	"strconv"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// VolumeCommand is a command that changes the volume of the audio output.
type VolumeCommand struct{}

// Aliases returns the current aliases for the command.
func (c *VolumeCommand) Aliases() []string {
	return viper.GetStringSlice("commands.volume.aliases")
}

// Description returns the description for the command.
func (c *VolumeCommand) Description() string {
	return viper.GetString("commands.volume.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *VolumeCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.volume.is_admin")
}

// Execute executes the command with the given user and arguments.
// Return value descriptions:
//    string: A message to be returned to the user upon successful execution.
//    bool:   Whether the message should be private or not. true = private,
//            false = public (sent to whole channel).
//    error:  An error message to be returned upon unsuccessful execution.
//            If no error has occurred, pass nil instead.
// Example return statement:
//    return "This is a private message!", true, nil
func (c *VolumeCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if len(args) == 0 {
		// Send the user the current volume level.
		return fmt.Sprintf(viper.GetString("commands.volume.messages.current_volume"), DJ.Volume), true, nil
	}

	newVolume, err := strconv.ParseFloat(args[0], 32)
	if err != nil {
		return "", true, errors.New(viper.GetString("commands.volume.messages.parsing_error"))
	}

	if newVolume <= viper.GetFloat64("volume.lowest") || newVolume >= viper.GetFloat64("volume.highest") {
		return "", true, fmt.Errorf(viper.GetString("commands.volume.messages.out_of_range_error"),
			viper.GetFloat64("volume.lowest"), viper.GetFloat64("volume.highest"))
	}

	newVolume32 := float32(newVolume)

	if DJ.AudioStream != nil {
		DJ.AudioStream.Volume = newVolume32
	}
	DJ.Volume = newVolume32

	return fmt.Sprintf(viper.GetString("commands.volume.messages.volume_changed"),
		user.Name, newVolume32), false, nil
}
