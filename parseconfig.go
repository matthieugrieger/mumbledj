/*
 * MumbleDJ
 * By Matthieu Grieger
 * parseconfig.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

import (
	"code.google.com/p/gcfg"
	"errors"
	"fmt"
)

// Golang struct representation of mumbledj.gcfg file structure for parsing.
type DjConfig struct {
	General struct {
		CommandPrefix     string
		SkipRatio         float32
		PlaylistSkipRatio float32
	}
	Volume struct {
		DefaultVolume float32
		LowestVolume  float32
		HighestVolume float32
	}
	Aliases struct {
		AddAlias               string
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
		KillAlias              string
	}
	Permissions struct {
		AdminsEnabled     bool
		Admins            []string
		AdminAdd          bool
		AdminAddPlaylists bool
		AdminSkip         bool
		AdminHelp         bool
		AdminVolume       bool
		AdminMove         bool
		AdminReload       bool
		AdminReset        bool
		AdminNumSongs     bool
		AdminNextSong     bool
		AdminKill         bool
	}
}

// Loads mumbledj.gcfg into dj.conf, a variable of type DjConfig.
func loadConfiguration() error {
	if gcfg.ReadFileInto(&dj.conf, fmt.Sprintf("%s/.mumbledj/config/mumbledj.gcfg", dj.homeDir)) == nil {
		return nil
	} else {
		return errors.New("Configuration load failed.")
	}
}
