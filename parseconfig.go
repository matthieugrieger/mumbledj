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
	"os/user"
)

type DjConfig struct {
	General struct {
		CommandPrefix string
		SkipRatio     float32
	}
	Volume struct {
		DefaultVolume float32
		LowestVolume  float32
		HighestVolume float32
	}
	Aliases struct {
		AddAlias       string
		SkipAlias      string
		AdminSkipAlias string
		VolumeAlias    string
		MoveAlias      string
		ReloadAlias    string
		KillAlias      string
	}
	Permissions struct {
		AdminsEnabled bool
		Admins        []string
		AdminAdd      bool
		AdminSkip     bool
		AdminVolume   bool
		AdminMove     bool
		AdminReload   bool
		AdminKill     bool
	}
}

func loadConfiguration() error {
	usr, err := user.Current()
	if err == nil {
		return gcfg.ReadFileInto(&dj.conf, fmt.Sprintf("%s/.mumbledj/config/mumbledj.gcfg", usr.HomeDir))
	} else {
		return errors.New("Configuration load failed.")
	}
}
