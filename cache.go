/*
 * MumbleDJ
 * By Matthieu Grieger
 * cache.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

// ByAge is a type that holds file information for the cache items.
type ByAge []os.FileInfo

func (a ByAge) Len() int {
	return len(a)
}
func (a ByAge) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByAge) Less(i, j int) bool {
	return time.Since(a[i].ModTime()) < time.Since(a[j].ModTime())
}

// SongCache is a struct that holds the number of songs currently cached and
// their combined file size.
type SongCache struct {
	NumSongs      int
	TotalFileSize int64
}

// NewSongCache creates an empty SongCache.
func NewSongCache() *SongCache {
	newCache := &SongCache{
		NumSongs:      0,
		TotalFileSize: 0,
	}
	return newCache
}

// GetNumSongs returns the number of songs currently cached.
func (c *SongCache) GetNumSongs() int {
	songs, _ := ioutil.ReadDir(fmt.Sprintf("%s/.mumbledj/songs", dj.homeDir))
	return len(songs)
}

// GetCurrentTotalFileSize calculates the total file size of the files within
// the cache and returns it.
func (c *SongCache) GetCurrentTotalFileSize() int64 {
	var totalSize int64
	songs, _ := ioutil.ReadDir(fmt.Sprintf("%s/.mumbledj/songs", dj.homeDir))
	for _, song := range songs {
		totalSize += song.Size()
	}
	return totalSize
}

// CheckMaximumDirectorySize checks the cache directory to determine if the filesize
// of the songs within exceed the user-specified size limit. If so, the oldest files
// get cleared until it is no longer exceeding the limit.
func (c *SongCache) CheckMaximumDirectorySize() {
	for c.GetCurrentTotalFileSize() > (dj.conf.Cache.MaximumSize * 1048576) {
		if err := c.ClearOldest(); err != nil {
			break
		}
	}
}

// Update updates the SongCache struct.
func (c *SongCache) Update() {
	c.NumSongs = c.GetNumSongs()
	c.TotalFileSize = c.GetCurrentTotalFileSize()
}

// ClearExpired clears cache items that are older than the cache period set within
// the user configuration.
func (c *SongCache) ClearExpired() {
	for range time.Tick(5 * time.Minute) {
		songs, _ := ioutil.ReadDir(fmt.Sprintf("%s/.mumbledj/songs", dj.homeDir))
		for _, song := range songs {
			hours := time.Since(song.ModTime()).Hours()
			if hours >= dj.conf.Cache.ExpireTime {
				if dj.queue.Len() > 0 {
					if (dj.queue.CurrentSong().Filename()) != song.Name() {
						os.Remove(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, song.Name()))
					}
				} else {
					os.Remove(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, song.Name()))
				}
			}
		}
	}
}

// ClearOldest deletes the oldest item in the cache.
func (c *SongCache) ClearOldest() error {
	songs, _ := ioutil.ReadDir(fmt.Sprintf("%s/.mumbledj/songs", dj.homeDir))
	sort.Sort(ByAge(songs))
	if dj.queue.Len() > 0 {
		if (dj.queue.CurrentSong().Filename()) != songs[0].Name() {
			return os.Remove(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, songs[0].Name()))
		}
		return errors.New("Song is currently playing.")
	}
	return os.Remove(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, songs[0].Name()))
}
