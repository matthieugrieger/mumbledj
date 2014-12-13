/*
 * MumbleDJ
 * By Matthieu Grieger
 * parseconfig.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"
	"github.com/BurntSushi/toml"
)

type djConfig struct {
	title string
	general generalConfig
	volume volumeConfig
	aliases aliasConfig `toml:"command-aliases"`
	permissions permissionsConfig
}

type generalConfig struct {
	defaultChannel string `toml:"default_channel"`
	commandPrefix string `toml:"command_prefix"`
	skipRatio float32 `toml:"skip_ratio"`
}

type volumeConfig struct {
	defaultVolume float32 `toml:"default_volume"`
	lowestVolume float32 `toml:"lowest_volume"`
	highestVolume float32 `toml:"highest_volume"`
}

type aliasConfig struct {
	addAlias string `toml:"add_alias"`
	skipAlias string `toml:"skip_alias"`
	adminSkipAlias string `toml:"admin_skip_alias"`
	volumeAlias string `toml:"volume_alias"`
	moveAlias string `toml:"move_alias"`
	reloadAlias string `toml:"reload_alias"`
	killAlias string `toml:"kill_alias"`
}

type permissionsConfig struct {
	adminsEnabled bool `toml:"enable_admins"`
	adminList []string `toml:"admins"`
	adminAdd bool `toml:"admin_add"`
	adminSkip bool `toml:"admin_skip"`
	adminVolume bool `toml:"admin_volume"`
	adminMove bool `toml:"admin_move"`
	adminReload bool `toml:"admin_reload"`
	adminKill bool `toml:"admin_kill"`
}

func loadConfiguration() (djConfig, error) {
	var conf djConfig
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		return conf, errors.New("Configuration load failed.")
	}
	return conf, nil
}
