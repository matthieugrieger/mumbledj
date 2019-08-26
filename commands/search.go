/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/search.go
 * Copyright (c) 2019 Reikion (MIT License)
 */

package commands

import (
	"fmt"

	"strings"

	"github.com/spf13/viper"
	"go.reik.pl/mumbledj/interfaces"
	"layeh.com/gumble/gumble"
)

var (
//	errAnotherSteamActive = errors.New("Stream is playing already")
//	once                  sync.Once
)

// SearchCommand is a command that plays random Frieza laughs from Dragon Ball series
type SearchCommand struct {
}

// Aliases returns the current aliases for the command.
func (c *SearchCommand) Aliases() []string {
	return viper.GetStringSlice("commands.search.aliases")
}

// Description returns the description for the command.
func (c *SearchCommand) Description() string {
	return viper.GetString("commands.search.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *SearchCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.search.is_admin")
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
func (c *SearchCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if len(args) == 0 {
		return viper.GetString("commands.search.messages.not_enough_arguments"), true, nil
	}

	var service interfaces.Service
	for _, serv := range DJ.AvailableServices {
		if serv.GetReadableName() == viper.GetString("search.service") {
			service = serv
			break
		}
	}
	if service == nil {
		return viper.GetString("commands.search.messages.service_not_specified"), true, nil
	}

	track, err := service.SearchTrack(strings.Join(args, " "), user)
	if err != nil {
		return "", true, err
	}
	go DJ.Queue.AppendTrack(track)
	if err != nil {
		return "", true, err
	}
	return fmt.Sprintf(viper.GetString("commands.add.messages.one_track_added"),
		user.Name, track.GetTitle(), track.GetService()), false, nil
}
