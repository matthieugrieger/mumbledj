/*
 * MumbleDJ
 * By Matthieu Grieger
 * parseconfig.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"
	"fmt"

	"github.com/scalingdata/gcfg"
)

// DjConfig is a Golang struct representation of mumbledj.gcfg file structure for parsing.
type DjConfig struct {
	General struct {
		CommandPrefix      string
		SkipRatio          float32
		PlaylistSkipRatio  float32
		DefaultComment     string
		MaxSongDuration    int
		MaxSongPerPlaylist int
		AutomaticShuffleOn bool
		AnnounceNewTrack   bool
		PlayerCommand      string
	}
	Cache struct {
		Enabled     bool
		MaximumSize int64
		ExpireTime  float64
	}
	Volume struct {
		DefaultVolume float32
		LowestVolume  float32
		HighestVolume float32
	}
	Aliases struct {
		AddAlias               string
		SearchAlias            string
		AddNextAlias           string
		SkipAlias              string
		SkipPlaylistAlias      string
		AdminSkipAlias         string
		AdminSkipPlaylistAlias string
		HelpAlias              string
		VolumeAlias            string
		MoveAlias              string
		ReloadAlias            string
		ResetAlias             string
		NumSongsAlias          string
		NextSongAlias          string
		CurrentSongAlias       string
		SetCommentAlias        string
		NumCachedAlias         string
		CacheSizeAlias         string
		KillAlias              string
		ShuffleAlias           string
		ShuffleOnAlias         string
		ShuffleOffAlias        string
		ListSongsAlias         string
		VersionAlias           string
	}
	Permissions struct {
		AdminsEnabled      bool
		Admins             []string
		AdminAdd           bool
		AdminSearch        bool
		AdminAddNext       bool
		AdminAddPlaylists  bool
		AdminSkip          bool
		AdminHelp          bool
		AdminVolume        bool
		AdminMove          bool
		AdminReload        bool
		AdminReset         bool
		AdminNumSongs      bool
		AdminNextSong      bool
		AdminCurrentSong   bool
		AdminSetComment    bool
		AdminNumCached     bool
		AdminCacheSize     bool
		AdminKill          bool
		AdminShuffle       bool
		AdminShuffleToggle bool
		AdminListSongs     bool
		AdminVersion       bool
	}
	ServiceKeys struct {
		Youtube    string
		SoundCloud string
	}
}

// Loads mumbledj.gcfg into dj.conf, a variable of type DjConfig.
func loadConfiguration() error {
	if gcfg.ReadFileInto(&dj.conf, fmt.Sprintf("%s/.mumbledj/config/mumbledj.gcfg", dj.homeDir)) == nil {
		return nil
	}
	fmt.Printf("%s/.mumbledj/config/mumbledj.gcfg\n", dj.homeDir)
	return errors.New("Configuration load failed.")
}
