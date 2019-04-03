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
	"math/rand"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleffmpeg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

var errAnotherSteamActive = errors.New("Stream is playing already")
var mutex sync.Mutex

// OhohohoCommand is a command that plays random Frieza laughs from Dragon Ball series
type OhohohoCommand struct{}

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

	mutex.Lock()
	if DJ.AudioStream != nil {
		DJ.AudioStream.Stop()
		DJ.AudioStream = nil
	}
	mutex.Unlock()

	if len(args) == 0 {
		var sb strings.Builder
		for k := range samplesList {
			sb.WriteString("<br>")
			sb.WriteString(" - ")
			sb.WriteString(k)
		}
		logrus.Println(fmt.Sprintf(viper.GetString("commands.ohohoho.messages.available_samples"), sb.String()))
		return fmt.Sprintf(viper.GetString("commands.ohohoho.messages.available_samples"), sb.String()), true, nil
	} else if len(args) == 1 {
		msg, pub, err := waitForRandomOhohoho(args[0])
		if err != nil {
			return msg, pub, err
		}
	} else if len(args) == 2 {

		howMany, err := strconv.Atoi(args[1])
		if err != nil || howMany < 1 || howMany > 10 {
			return "", true, errors.New(viper.GetString("commands.ohohoho.messages.how_many_times_error"))
		}

		for i := 0; i < howMany; i++ {
			msg, pub, err := waitForRandomOhohoho(args[0])
			if err != nil {
				if err == errAnotherSteamActive {
					// it's ok, don't inform user, that its request is interrupted
					return "", true, nil
				}
				return msg, pub, err
			}
		}
	}

	return "", true, nil
}
func waitForRandomOhohoho(whichSample string) (string, bool, error) {
	mutex.Lock()
	// ensure that Ohohoho hasn't started in the meantime
	if DJ.AudioStream != nil && DJ.AudioStream.State() != gumbleffmpeg.StateStopped {
		// oh no, it's started already
		mutex.Unlock()
		return "", true, errAnotherSteamActive
	}
	err := playRandomOhohoho(whichSample)
	if err != nil {
		mutex.Unlock()
		return "", true, err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	if DJ.AudioStream != nil {
		go func() {
			DJ.AudioStream.Wait()
			wg.Done()
		}()
	}
	mutex.Unlock()
	wg.Wait()

	return "", true, nil
}

func playRandomOhohoho(whichSample string) error {
	//check is sample exists
	if _, ok := samplesList[whichSample]; !ok {
		return fmt.Errorf(viper.GetString("commands.ohohoho.messages.sample_not_exists_error"), whichSample)
	}

	noOfSamples := samplesList[whichSample]
	// rand rands from [0;n), so we need +1 to scale to [1;n]
	chosenRandom := strconv.Itoa(rand.Intn(noOfSamples) + 1)
	ohohohoSample, err := Assets.Open(path.Join(whichSample, chosenRandom+".flac"))
	if err != nil {
		return fmt.Errorf(viper.GetString("commands.ohohoho.messages.internal_sample_error"), chosenRandom)
	}

	source := gumbleffmpeg.SourceReader(ohohohoSample)
	DJ.AudioStream = gumbleffmpeg.New(DJ.Client, source)
	DJ.AudioStream.Volume = DJ.Volume
	DJ.AudioStream.Play()
	return nil
}
