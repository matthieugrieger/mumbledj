/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/cachesize.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"

	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// CacheSizeCommand is a command that outputs the current size of the cache.
type CacheSizeCommand struct{}

// Aliases returns the current aliases for the command.
func (c *CacheSizeCommand) Aliases() []string {
	return viper.GetStringSlice("commands.cachesize.aliases")
}

// Description returns the description for the command.
func (c *CacheSizeCommand) Description() string {
	return viper.GetString("commands.cachesize.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *CacheSizeCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.cachesize.is_admin")
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
func (c *CacheSizeCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	const bytesInMiB = 1048576

	if !viper.GetBool("cache.enabled") {
		return "", true, errors.New(viper.GetString("commands.common_messages.caching_disabled_error"))
	}

	DJ.Cache.UpdateStatistics()
	return fmt.Sprintf(viper.GetString("commands.cachesize.messages.current_size"), DJ.Cache.TotalFileSize/bytesInMiB), true, nil
}
