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

type SongCache struct {
	NumSongs      int
	TotalFileSize int64
}

func NewSongCache() *SongCache {
	newCache := &SongCache{
		NumSongs:      0,
		TotalFileSize: 0,
	}
	return newCache
}

func (c *SongCache) GetNumSongs() int {
	songs, _ := ioutil.ReadDir(fmt.Sprintf("%s/.mumbledj/songs", dj.homeDir))
	return len(songs)
}

func (c *SongCache) GetCurrentTotalFileSize() int64 {
	var totalSize int64 = 0
	songs, _ := ioutil.ReadDir(fmt.Sprintf("%s/.mumbledj/songs", dj.homeDir))
	for _, song := range songs {
		totalSize += song.Size()
	}
	return totalSize
}

func (c *SongCache) CheckMaximumDirectorySize() {
	for c.GetCurrentTotalFileSize() > (dj.conf.Cache.MaximumSize * 1048576) {
		if err := c.ClearOldest(); err != nil {
			break
		}
	}
}

func (c *SongCache) Update() {
	c.NumSongs = c.GetNumSongs()
	c.TotalFileSize = c.GetCurrentTotalFileSize()
}

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

func (c *SongCache) ClearOldest() error {
	songs, _ := ioutil.ReadDir(fmt.Sprintf("%s/.mumbledj/songs", dj.homeDir))
	sort.Sort(ByAge(songs))
	if dj.queue.Len() > 0 {
		if (dj.queue.CurrentSong().Filename()) != songs[0].Name() {
			return os.Remove(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, songs[0].Name()))
		} else {
			return errors.New("Song is currently playing.")
		}
	} else {
		return os.Remove(fmt.Sprintf("%s/.mumbledj/songs/%s", dj.homeDir, songs[0].Name()))
	}
}
