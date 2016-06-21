/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/config.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// SetDefaultConfig sets default values for all configuration options.
func SetDefaultConfig() {
	// API key defaults.
	viper.SetDefault("api_keys.youtube", "")
	viper.SetDefault("api_keys.soundcloud", "")

	// General defaults.
	viper.SetDefault("defaults.comment", "Hello! I am a bot. Type !help for a list of commands.")
	viper.SetDefault("defaults.channel", "")
	viper.SetDefault("defaults.player_command", "ffmpeg")

	// Queue defaults.
	viper.SetDefault("queue.track_skip_ratio", 0.5)
	viper.SetDefault("queue.playlist_skip_ratio", 0.5)
	viper.SetDefault("queue.max_track_duration", 0)
	viper.SetDefault("queue.max_tracks_per_playlist", 50)
	viper.SetDefault("queue.automatic_shuffle_on", false)
	viper.SetDefault("queue.announce_new_tracks", true)

	// Connection defaults.
	viper.SetDefault("connection.address", "127.0.0.1")
	viper.SetDefault("connection.port", 64738)
	viper.SetDefault("connection.password", "")
	viper.SetDefault("connection.username", "MumbleDJ")
	viper.SetDefault("connection.insecure", false)
	viper.SetDefault("connection.cert", "")
	viper.SetDefault("connection.key", "")
	viper.SetDefault("connection.access_tokens", "")
	viper.SetDefault("connection.retry_enabled", true)
	viper.SetDefault("connection.retry_attempts", 10)
	viper.SetDefault("connection.retry_interval", 5)

	// Cache defaults.
	viper.SetDefault("cache.enabled", false)
	viper.SetDefault("cache.maximum_size", "512MiB")
	viper.SetDefault("cache.expire_time", 24)
	viper.SetDefault("cache.check_interval", 5)
	viper.SetDefault("cache.directory", "$HOME/.cache/mumbledj")

	// Volume defaults.
	viper.SetDefault("volume.default", 0.2)
	viper.SetDefault("volume.lowest", 0.01)
	viper.SetDefault("volume.highest", 0.8)

	// Admins defaults.
	viper.SetDefault("admins.enabled", true)
	viper.SetDefault("admins.names", []string{"SuperUser"})

	// Command defaults.
	viper.SetDefault("commands.add.aliases", []string{"add", "a"})
	viper.SetDefault("commands.add.is_admin", false)
	viper.SetDefault("commands.add.description", "Adds a track or playlist from a media site to the queue.")

	viper.SetDefault("commands.addnext.aliases", []string{"addnext", "an"})
	viper.SetDefault("commands.addnext.is_admin", true)
	viper.SetDefault("commands.addnext.description", "Adds a track or playlist from a media site as the next item in the queue.")

	viper.SetDefault("commands.cachesize.aliases", []string{"cachesize", "cs"})
	viper.SetDefault("commands.cachesize.is_admin", true)
	viper.SetDefault("commands.cachesize.description", "Outputs the file size of the cache in MiB if caching is enabled.")

	viper.SetDefault("commands.currenttrack.aliases", []string{"currenttrack", "currentsong", "current"})
	viper.SetDefault("commands.currenttrack.is_admin", false)
	viper.SetDefault("commands.currenttrack.description", "Outputs information about the current track in the queue if one exists.")

	viper.SetDefault("commands.forceskip.aliases", []string{"forceskip", "fs"})
	viper.SetDefault("commands.forceskip.is_admin", true)
	viper.SetDefault("commands.forceskip.description", "Immediately skips the current track.")

	viper.SetDefault("commands.forceskipplaylist.aliases", []string{"forceskipplaylist", "fsp"})
	viper.SetDefault("commands.forceskipplaylist.is_admin", true)
	viper.SetDefault("commands.forceskipplaylist.description", "Immediately skips the current playlist.")

	viper.SetDefault("commands.help.aliases", []string{"help", "h"})
	viper.SetDefault("commands.help.is_admin", false)
	viper.SetDefault("commands.help.description", "Outputs this list of commands.")

	viper.SetDefault("commands.joinme.aliases", []string{"joinme", "join"})
	viper.SetDefault("commands.joinme.is_admin", true)
	viper.SetDefault("commands.joinme.description", "Moves MumbleDJ into your current channel if not playing audio to someone else.")

	viper.SetDefault("commands.kill.aliases", []string{"kill", "k"})
	viper.SetDefault("commands.kill.is_admin", true)
	viper.SetDefault("commands.kill.description", "Stops the bot and cleans its cache directory.")

	viper.SetDefault("commands.listtracks.aliases", []string{"listtracks", "listsongs", "list", "l"})
	viper.SetDefault("commands.listtracks.is_admin", false)
	viper.SetDefault("commands.listtracks.description", "Outputs a list of the tracks currently in the queue.")

	viper.SetDefault("commands.move.aliases", []string{"move", "m"})
	viper.SetDefault("commands.move.is_admin", true)
	viper.SetDefault("commands.move.description", "Moves the bot into the Mumble channel provided via argument.")

	viper.SetDefault("commands.nexttrack.aliases", []string{"nexttrack", "nextsong", "next"})
	viper.SetDefault("commands.nexttrack.is_admin", false)
	viper.SetDefault("commands.nexttrack.description", "Outputs information about the next track in the queue if one exists.")

	viper.SetDefault("commands.numcached.aliases", []string{"numcached", "nc"})
	viper.SetDefault("commands.numcached.is_admin", true)
	viper.SetDefault("commands.numcached.description", "Outputs the number of tracks cached on disk if caching is enabled.")

	viper.SetDefault("commands.numtracks.aliases", []string{"numtracks", "numsongs", "nt"})
	viper.SetDefault("commands.numtracks.is_admin", false)
	viper.SetDefault("commands.numtracks.description", "Outputs the number of tracks currently in the queue.")

	viper.SetDefault("commands.pause.aliases", []string{"pause"})
	viper.SetDefault("commands.pause.is_admin", false)
	viper.SetDefault("commands.pause.description", "Pauses audio playback.")

	viper.SetDefault("commands.reload.aliases", []string{"reload", "r"})
	viper.SetDefault("commands.reload.is_admin", true)
	viper.SetDefault("commands.reload.description", "Reloads the configuration file.")

	viper.SetDefault("commands.reset.aliases", []string{"reset", "re"})
	viper.SetDefault("commands.reset.is_admin", true)
	viper.SetDefault("commands.reset.description", "Resets the queue by removing all queue items.")

	viper.SetDefault("commands.resume.aliases", []string{"resume"})
	viper.SetDefault("commands.resume.is_admin", false)
	viper.SetDefault("commands.resume.description", "Resumes audio playback.")

	viper.SetDefault("commands.setcomment.aliases", []string{"setcomment", "comment", "sc"})
	viper.SetDefault("commands.setcomment.is_admin", true)
	viper.SetDefault("commands.setcomment.description", "Sets the comment displayed next to MumbleDJ's username in Mumble.")

	viper.SetDefault("commands.shuffle.aliases", []string{"shuffle", "shuf", "sh"})
	viper.SetDefault("commands.shuffle.is_admin", true)
	viper.SetDefault("commands.shuffle.description", "Randomizes the tracks currently in the queue.")

	viper.SetDefault("commands.skip.aliases", []string{"skip", "s"})
	viper.SetDefault("commands.skip.is_admin", false)
	viper.SetDefault("commands.skip.description", "Places a vote to skip the current track.")

	viper.SetDefault("commands.skipplaylist.aliases", []string{"skipplaylist", "sp"})
	viper.SetDefault("commands.skipplaylist.is_admin", false)
	viper.SetDefault("commands.skipplaylist.description", "Places a vote to skip the current playlist.")

	viper.SetDefault("commands.toggleshuffle.aliases", []string{"toggleshuffle", "toggleshuf", "togshuf", "tsh"})
	viper.SetDefault("commands.toggleshuffle.is_admin", true)
	viper.SetDefault("commands.toggleshuffle.description", "Toggles automatic track shuffling on/off.")

	viper.SetDefault("commands.version.aliases", []string{"version"})
	viper.SetDefault("commands.version.is_admin", false)
	viper.SetDefault("commands.version.description", "Outputs the current version of MumbleDJ.")

	viper.SetDefault("commands.volume.aliases", []string{"volume", "vol", "v"})
	viper.SetDefault("commands.volume.is_admin", false)
	viper.SetDefault("commands.volume.description", "Changes the volume if an argument is provided, outputs the current volume otherwise.")
}

// ReadConfigFile reads in the config file and updates the configuration accordingly.
func ReadConfigFile() error {
	logrus.Infoln("Reading config...")
	return viper.ReadInConfig()
}

// CheckForDuplicateAliases validates that all commands have unique aliases.
func CheckForDuplicateAliases() error {
	var aliases []string

	logrus.Infoln("Checking for duplicate aliases...")

	// It would be preferred to use viper.Sub("aliases") here, but there are some
	// nil pointer dereferencing issues.
	for _, setting := range viper.AllKeys() {
		if strings.HasSuffix(setting, "aliases") {
			aliases = append(aliases, viper.GetStringSlice(setting)...)
		}
	}

	// Sort the strings to allow us to fail faster in case there is a duplicate.
	sort.Strings(aliases)

	for i := 0; i < len(aliases)-1; i++ {
		if aliases[i] == aliases[i+1] {
			return fmt.Errorf("Duplicate alias found: %s", aliases[i])
		}
	}

	return nil
}
