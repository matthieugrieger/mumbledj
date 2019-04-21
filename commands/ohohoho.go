/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/ohohoho.go
 * Copyright (c) 2019 Reikion (MIT License)
 */

package commands

import (
	"errors"
	"fmt"
	"regexp"
	"reik.pl/mumbledj/bot"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"layeh.com/gumble/gumble"
)

var samplesList = map[string]int{}

func init() {
	assetsDirs := Assets.List()
	//match dirs
	// ex. ohohoho/1.flac and ohohoho is [0][1] submatch
	reg := regexp.MustCompile("(.+?)/.*")

	for _, el := range assetsDirs {
		matches := reg.FindAllStringSubmatch(el, -1)
		if matches != nil && Assets.HasDir(matches[0][1]) {
			// count files in folder by the way
			samplesList[matches[0][1]]++
		}
	}
}

var (
	errAnotherSteamActive = errors.New("Stream is playing already")
	once                  sync.Once
)

// OhohohoCommand is a command that plays random Frieza laughs from Dragon Ball series
type OhohohoCommand struct {
}

// Aliases returns the current aliases for the command.
func (c *OhohohoCommand) Aliases() []string {
	return viper.GetStringSlice("commands.ohohoho.aliases")
}

// Description returns the description for the command.
func (c *OhohohoCommand) Description() string {
	return viper.GetString("commands.ohohoho.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *OhohohoCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.ohohoho.is_admin")
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
func (c *OhohohoCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if len(args) == 0 {
		var sb strings.Builder
		for k := range bot.GetSampleList() {
			sb.WriteString("<br>")
			sb.WriteString(" - ")
			sb.WriteString(k)
		}
		logrus.Println(fmt.Sprintf(viper.GetString("commands.ohohoho.messages.available_samples"), sb.String()))
		return fmt.Sprintf(viper.GetString("commands.ohohoho.messages.available_samples"), sb.String()), true, nil
	}

	if len(args) == 1 {
		err := DJ.Ohohoho.PlaySample(args[0], 1)
		if err != nil {
			return "", true, err
		}
	}

	if len(args) == 2 {
		if args[0] == "s" {
			if args[1] == "stop" {
				logrus.Debugln("Stopping via cmd...")
				DJ.Ohohoho.Stop()
				return "Stopping...", true, nil
			}
			if args[1] == "empty" {
				logrus.Debugln("Emptying via cmd...")
				DJ.Ohohoho.EmptyStop()
				return "Emptying stop...", true, nil
			}
		}
		howMany, err := strconv.Atoi(args[1])
		// second argument is empty, probably space, play once
		if args[1] == "" {
			err := DJ.Ohohoho.PlaySample(args[0], 1)
			if err != nil {
				return "", true, err
			}
		} else if err != nil || howMany < 1 || howMany > 10 {
			return "", true, errors.New(viper.GetString("commands.ohohoho.messages.how_many_times_error"))
		} else {
			err := DJ.Ohohoho.PlaySample(args[0], howMany)
			//msg, pub, err := c.waitForRandomOhohoho(args[0])
			if err != nil {
				return "", true, err
			}
		}
	}
	return "", true, nil
}
