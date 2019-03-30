/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/cache.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SortFilesByAge is a type that holds file information for cached items for
// sorting.
type SortFilesByAge []os.FileInfo

// Len returns the length of the file slice.
func (a SortFilesByAge) Len() int {
	return len(a)
}

// Swap swaps two elements in the file slice.
func (a SortFilesByAge) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less compares two file modification times to determine if one is less than
// the other. Returns true if the item in index i is older than the item in
// index j, false otherwise.
func (a SortFilesByAge) Less(i, j int) bool {
	return time.Since(a[i].ModTime()) < time.Since(a[j].ModTime())
}

// Cache keeps track of the filesize of the audio cache and
// provides methods for pruning the cache.
type Cache struct {
	NumAudioFiles int
	TotalFileSize int64
}

// NewCache creates an empty Cache and returns it.
func NewCache() *Cache {
	return &Cache{
		NumAudioFiles: 0,
		TotalFileSize: 0,
	}
}

// CheckDirectorySize checks the cache directory to determine if the filesize
// of the files within exceed the user-specified size limit. If so, the oldest
// files are cleared until it is no longer exceeding the limit.
func (c *Cache) CheckDirectorySize() {
	const bytesInMiB int = 1048576

	c.UpdateStatistics()
	for c.TotalFileSize > int64(viper.GetInt("cache.maximum_size")*bytesInMiB) {
		if err := c.DeleteOldest(); err != nil {
			break
		}
	}
}

// UpdateStatistics updates the statistics relevant to the cache (number of
// audio files cached, total current size of the cache).
func (c *Cache) UpdateStatistics() {
	c.NumAudioFiles, c.TotalFileSize = c.getCurrentStatistics()
	logrus.WithFields(logrus.Fields{
		"num_audio_files": c.NumAudioFiles,
		"total_file_size": c.TotalFileSize,
	}).Infoln("Updated cache statistics.")
}

// CleanPeriodically loops forever, deleting expired cached audio files as necessary.
func (c *Cache) CleanPeriodically() {
	for range time.Tick(time.Duration(viper.GetInt("cache.check_interval")) * time.Minute) {
		logrus.Infoln("Checking cache for expired files...")
		files, _ := ioutil.ReadDir(os.ExpandEnv(viper.GetString("cache.directory")))
		for _, file := range files {
			// It is safe to check the modification time because when audio files are
			// played their modification time is updated. This ensures that audio
			// files will not get deleted while they are playing, assuming a reasonable
			// expiry time is set in the configuration.
			hours := time.Since(file.ModTime()).Hours()
			if hours >= viper.GetFloat64("cache.expire_time") {
				logrus.WithFields(logrus.Fields{
					"expired_file": file.Name(),
				}).Infoln("Removing expired cache entry.")
				os.Remove(fmt.Sprintf("%s/%s", os.ExpandEnv(viper.GetString("cache.directory")), file.Name()))
			}
		}
	}
}

// DeleteOldest deletes the oldest file in the cache.
func (c *Cache) DeleteOldest() error {
	files, _ := ioutil.ReadDir(os.ExpandEnv(viper.GetString("cache.directory")))
	if len(files) > 0 {
		sort.Sort(SortFilesByAge(files))
		os.Remove(fmt.Sprintf("%s/%s", os.ExpandEnv(viper.GetString("cache.directory")), files[0].Name()))
		return nil
	}
	return errors.New("There are no files currently cached")
}

// DeleteAll deletes all cached audio files.
func (c *Cache) DeleteAll() error {
	dir, err := os.Open(os.ExpandEnv(viper.GetString("cache.directory")))
	if err != nil {
		return err
	}

	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	logrus.Infoln("Deleting all cached audio files...")
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(os.ExpandEnv(viper.GetString("cache.directory")), name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) getCurrentStatistics() (int, int64) {
	var totalSize int64
	files, _ := ioutil.ReadDir(os.ExpandEnv(viper.GetString("cache.directory")))
	for _, file := range files {
		totalSize += file.Size()
	}
	return len(files), totalSize
}
