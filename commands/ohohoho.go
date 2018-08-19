/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/add.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"sync"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleffmpeg"
	"github.com/spf13/viper"
)

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
		return waitForRandomOhohoho()
	} else {
		howMany, err := strconv.Atoi(args[0])
		if err != nil || howMany < 1 || howMany > 10 {
			return "", true, errors.New(viper.GetString("commands.ohohoho.messages.how_many_times_error"))
		}

		for i := 0; i < howMany; i++ {
			msg, pub, err := waitForRandomOhohoho()
			if err != nil {
				if err == errAnotherSteamActive {
					// it's ok, don't inform user, that it's request is interrupted
					return "", true, nil
				}
				return msg, pub, err
			}
		}
	}

	return "", true, nil
}
func waitForRandomOhohoho() (string, bool, error) {
	mutex.Lock()
	// ensure that Ohohoho hasn't started in the meantime
	if DJ.AudioStream != nil && DJ.AudioStream.State() != gumbleffmpeg.StateStopped {
		// oh no, it's started already
		mutex.Unlock()
		return "", true, errAnotherSteamActive
	}
	err := playRandomOhohoho()
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

func playRandomOhohoho() error {
	// there are 56 files hohoho files from 1.flac to 56.flac
	chosenRandom := strconv.Itoa(rand.Intn(55) + 1)
	ohohohoSample, err := Assets.Open(path.Join("hohoho", chosenRandom+".flac"))
	if err != nil {
		return fmt.Errorf(viper.GetString("commands.ohohoho.messages.internal_sample_error"), chosenRandom)
	}

	source := gumbleffmpeg.SourceReader(ohohohoSample)
	DJ.AudioStream = gumbleffmpeg.New(DJ.Client, source)
	DJ.AudioStream.Volume = DJ.Volume
	DJ.AudioStream.Play()
	return nil
}
