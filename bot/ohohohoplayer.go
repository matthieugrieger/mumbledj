/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/ohohoho.go
 * Copyright (c) 2019 Reikion (MIT License)
 */

package bot

import (
	"errors"
	"math/rand"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"sync"
	"time"

	"go.reik.pl/mumbledj/assets"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"layeh.com/gumble/gumbleffmpeg"
)

// Assets embedded in binary
var Assets = assets.Assets

// After init() samplesList contains folders name and number of files in folder
var samplesList = map[string]int{}

// GetSampleList returns map[string]int with folders name and number of files in folder
func GetSampleList() map[string]int {
	return samplesList
}

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
	errAnotherSteamActive  = errors.New("Stream is playing already")
	errSampleNotFound      = errors.New("Sample not found")
	errInternalSampleError = errors.New("Internal sample error")
	once                   sync.Once
)

// trackStreamInfo is for boxing source and offset for channel usage
type trackStreamInfo struct {
	source gumbleffmpeg.Source
	offset time.Duration
}

// OhohohoPlayer is a command that plays random Frieza laughs from Dragon Ball series
type OhohohoPlayer struct {
	mutex            sync.Mutex
	restorePrevious  bool // informs OhohohoPlayer we should restore track from queue
	ohohohoPlaying   bool // we are playing sample
	stopPlaying      chan trackStreamInfo
	stopSamplePlayer chan struct{}
}

// NewOhohohoPlayer returns new instance of OhohohoPlayer
func NewOhohohoPlayer() *OhohohoPlayer {
	return &OhohohoPlayer{
		stopPlaying:      make(chan trackStreamInfo),
		stopSamplePlayer: make(chan struct{}, 1),
	}
}

func (c *OhohohoPlayer) Stop() error {
	c.stopPlaying <- trackStreamInfo{}
	return nil
}

func (c *OhohohoPlayer) EmptyStop() error {
	c.stopPlaying <- trackStreamInfo{}
	return nil
}

// IsInterrupting informs if Queue command should remove track from queue. If it's true, track should remain on list
func (c *OhohohoPlayer) IsInterrupting() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.restorePrevious
}

func (c *OhohohoPlayer) prepareSample(sampleSetName string, howMany int) (gumbleffmpeg.Source, time.Duration, chan error) {

	c.mutex.Lock()

	var (
		source gumbleffmpeg.Source
		offset time.Duration
	)

	// needed for PlaySample method signalling
	// if goroutine playing sample finished its work, nil is sent via done channel
	// else if error occured work is interrupted and err is sent via done channel
	// it isn't used if new sample or stop has been requested by uesr
	done := make(chan error)

	if DJ.AudioStream != nil {
		if c.ohohohoPlaying {
			// we're currently playing our sample, but not track, so stop previous one
			c.stopPlaying <- trackStreamInfo{}
			// gather offset and source of track from previous goroutin
			// we're playing, so we need to gather offset and source from previous goroutine
			logrus.Debugln("Waiting for offset and source from previous goroutine")
			t := <-c.stopPlaying
			logrus.Debugln("Response obtained")
			if c.restorePrevious == true {
				source = t.source
				offset = t.offset
			}
			// inform goroutine that it needs to end playing sample
			//logrus.Debugln("Informing that samplePlayer should stop its work")
			// needed to inform that previous loop iteration of samplePlayer ended, because user has stopped it externally
			//c.stopSamplePlayer <- struct{}{}
		} else {
			// looks like track from queue is playing
			lastTrack, err := DJ.Queue.CurrentTrack()
			if err == nil {
				// it's playing track from queue, so it interrupts original playlist
				// get information about track to resume it in the future
				filepath := os.ExpandEnv(viper.GetString("cache.directory") + "/" + lastTrack.GetFilename())
				source = gumbleffmpeg.SourceFile(filepath)
				offset = DJ.AudioStream.Elapsed()
				logrus.Infoln(offset, source)
				c.restorePrevious = true
			}
		}
		DJ.AudioStream.Stop()
		DJ.AudioStream = nil
	}

	c.ohohohoPlaying = true

	go func() {
		sampleName := sampleSetName
		for i := 0; i < howMany; i++ {
			select {
			case <-c.stopSamplePlayer:
				logrus.Debugln("Stopping sample player")
				return
			default:
				// there was no signal, continue normal work
				logrus.Debugln("Continuing playing next sample")
			}

			sample, err := c.openSample(sampleName)
			if err != nil {
				done <- err
				return
			}
			err = c.waitForRandomOhohoho(sample)
			if err != nil {
				done <- err
				return
			}
		}

		done <- nil
	}()

	c.mutex.Unlock()

	return source, offset, done

}

// PlaySample plays random file from folder given by user as argument, which is located in asset directory
func (c *OhohohoPlayer) PlaySample(sampleName string, howMany int) error {
	err := c.isSampleSetExisting(sampleName)
	if err != nil {
		return err
	}

	source, offset, done := c.prepareSample(sampleName, howMany)
	select {
	case <-c.stopPlaying:
		logrus.Debugln("Informing that samplePlayer should stop its work")
		// we need non-blocking request prepared, because at start of every function go scheduler can make context switch
		c.stopSamplePlayer <- struct{}{}
		logrus.Debugln("Stopping previous sample")
		DJ.AudioStream.Stop()
		// block until sample Player goroutine receive signal
		select {
		case c.stopSamplePlayer <- struct{}{}:
			// it hasn't blocked, so sample Player goroutine received message already
			// so clear channel
			<-c.stopSamplePlayer
		default:
			// it will block, so block
			c.stopSamplePlayer <- struct{}{}
			// we need to consume what we produced
			<-c.stopSamplePlayer
		}

		logrus.Debugln("Stopped previous sample")
		// ping prepareSample that stream is stopped and return needed data for restorePreviousTrack
		logrus.Debugln("Informing about previous track for unblocking")
		c.stopPlaying <- trackStreamInfo{source, offset}
		return nil
	case err = <-done:
		if err != nil {
			switch err {
			case errInternalSampleError:
				logrus.WithField("err", errInternalSampleError).Errorln("Critical error, check mumbledj source code")
				return err
			default:
				logrus.Debug("Done_err", err)
				return err
			}
		}
		logrus.Debugln("Done, sample finished")
		DJ.AudioStream = nil
	}

	if c.restorePrevious {
		c.restorePreviousTrack(source, offset)
	}

	//if c.ohohohoPlaying && !c.restorePrevious {
	//	DJ.AudioStream = nil
	//}

	c.ohohohoPlaying = false
	return nil
}

// IsSampleSetExisting checks if dir with samples exist in bundled assets
func (c *OhohohoPlayer) isSampleSetExisting(sampleName string) error {
	if _, ok := samplesList[sampleName]; !ok {
		return errSampleNotFound
	}
	return nil
}

// OpenSample try to open sample and returns opened file or nil, err if error occured
func (c *OhohohoPlayer) openSample(sampleName string) (http.File, error) {

	noOfSamples := samplesList[sampleName]
	// rand rands from [0;n), so we need +1 to scale to [1;n]
	chosenRandom := strconv.Itoa(rand.Intn(noOfSamples) + 1)
	sample, err := Assets.Open(path.Join(sampleName, chosenRandom+".flac"))
	if err != nil {
		return nil, errInternalSampleError
	}

	return sample, nil
}

func (c *OhohohoPlayer) waitForRandomOhohoho(sample http.File) error {
	// ensure that Ohohoho hasn't started in the meantime
	if DJ.AudioStream != nil && DJ.AudioStream.State() != gumbleffmpeg.StateStopped {
		logrus.Debugln(errAnotherSteamActive)
		return errAnotherSteamActive
	}

	c.playRandomOhohoho(sample)
	DJ.AudioStream.Wait()

	return nil
}

func (c *OhohohoPlayer) playRandomOhohoho(assetFile http.File) {
	source := gumbleffmpeg.SourceReader(assetFile)
	DJ.AudioStream = gumbleffmpeg.New(DJ.Client, source)
	DJ.AudioStream.Volume = DJ.Volume
	DJ.AudioStream.Play()
}

func (c *OhohohoPlayer) restorePreviousTrack(source gumbleffmpeg.Source, offset time.Duration) {
	// restore previous track
	// c.mutex.Lock()
	// defer c.mutex.Unlock()
	logrus.Infoln("Restoring previous track")
	DJ.AudioStream = gumbleffmpeg.New(DJ.Client, source)
	DJ.AudioStream.Offset = offset
	DJ.AudioStream.Volume = DJ.Volume
	DJ.AudioStream.Play()
	c.restorePrevious = false
	go func() {
		DJ.AudioStream.Wait()
		if DJ.Ohohoho.IsInterrupting() {
			// do not skip item from queue
			return
		}
		DJ.Queue.Skip()
	}()
}
