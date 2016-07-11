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
	viper.SetDefault("connection.user_p12", "")
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
	viper.SetDefault("commands.prefix", "!")
	viper.SetDefault("commands.common_messages.no_tracks_error", "There are no tracks in the queue.")
	viper.SetDefault("commands.common_messages.caching_disabled_error", "Caching is currently disabled.")

	viper.SetDefault("commands.add.aliases", []string{"add", "a"})
	viper.SetDefault("commands.add.is_admin", false)
	viper.SetDefault("commands.add.description", "Adds a track or playlist from a media site to the queue.")
	viper.SetDefault("commands.add.messages.no_url_error", "A URL must be supplied with the add command.")
	viper.SetDefault("commands.add.messages.no_valid_tracks_error", "No valid tracks were found with the provided URL(s).")
	viper.SetDefault("commands.add.messages.tracks_too_long_error", "Your track(s) were either too long or an error occurred while processing them. No track(s) have been added.")
	viper.SetDefault("commands.add.messages.one_track_added", "<b>%s</b> added <b>1</b> track to the queue:<br><i>%s</i> from %s")
	viper.SetDefault("commands.add.messages.many_tracks_added", "<b>%s</b> added <b>%d</b> tracks to the queue.")
	viper.SetDefault("commands.add.messages.num_tracks_too_long", "<br><b>%d</b> tracks could not be added due to error or because they are too long.")

	viper.SetDefault("commands.addnext.aliases", []string{"addnext", "an"})
	viper.SetDefault("commands.addnext.is_admin", true)
	viper.SetDefault("commands.addnext.description", "Adds a track or playlist from a media site as the next item in the queue.")

	viper.SetDefault("commands.cachesize.aliases", []string{"cachesize", "cs"})
	viper.SetDefault("commands.cachesize.is_admin", true)
	viper.SetDefault("commands.cachesize.description", "Outputs the file size of the cache in MiB if caching is enabled.")
	viper.SetDefault("commands.cachesize.messages.current_size", "The current size of the cache is <b>%.2v MiB</b>.")

	viper.SetDefault("commands.currenttrack.aliases", []string{"currenttrack", "currentsong", "current"})
	viper.SetDefault("commands.currenttrack.is_admin", false)
	viper.SetDefault("commands.currenttrack.description", "Outputs information about the current track in the queue if one exists.")
	viper.SetDefault("commands.currenttrack.messages.current_track", "The current track is <i>%s</i>, added by <b>%s</b>.")

	viper.SetDefault("commands.forceskip.aliases", []string{"forceskip", "fs"})
	viper.SetDefault("commands.forceskip.is_admin", true)
	viper.SetDefault("commands.forceskip.description", "Immediately skips the current track.")
	viper.SetDefault("commands.forceskip.messages.track_skipped", "The current track has been forcibly skipped by <b>%s</b>.")

	viper.SetDefault("commands.forceskipplaylist.aliases", []string{"forceskipplaylist", "fsp"})
	viper.SetDefault("commands.forceskipplaylist.is_admin", true)
	viper.SetDefault("commands.forceskipplaylist.description", "Immediately skips the current playlist.")
	viper.SetDefault("commands.forceskipplaylist.messages.no_playlist_error", "The current track is not part of a playlist.")
	viper.SetDefault("commands.forceskipplaylist.messages.playlist_skipped", "The current playlist has been forcibly skipped by <b>%s</b>.")

	viper.SetDefault("commands.help.aliases", []string{"help", "h"})
	viper.SetDefault("commands.help.is_admin", false)
	viper.SetDefault("commands.help.description", "Outputs this list of commands.")
	viper.SetDefault("commands.help.messages.commands_header", "<br><b>Commands:</b><br>")
	viper.SetDefault("commands.help.messages.admin_commands_header", "<br><b>Admin Commands:</b><br>")

	viper.SetDefault("commands.joinme.aliases", []string{"joinme", "join"})
	viper.SetDefault("commands.joinme.is_admin", true)
	viper.SetDefault("commands.joinme.description", "Moves MumbleDJ into your current channel if not playing audio to someone else.")
	viper.SetDefault("commands.joinme.messages.others_are_listening_error", "Users in another channel are listening to me.")
	viper.SetDefault("commands.joinme.messages.in_your_channel", "I am now in your channel!")

	viper.SetDefault("commands.kill.aliases", []string{"kill", "k"})
	viper.SetDefault("commands.kill.is_admin", true)
	viper.SetDefault("commands.kill.description", "Stops the bot and cleans its cache directory.")

	viper.SetDefault("commands.listtracks.aliases", []string{"listtracks", "listsongs", "list", "l"})
	viper.SetDefault("commands.listtracks.is_admin", false)
	viper.SetDefault("commands.listtracks.description", "Outputs a list of the tracks currently in the queue.")
	viper.SetDefault("commands.listtracks.messages.invalid_integer_error", "An invalid integer was supplied.")
	viper.SetDefault("commands.listtracks.messages.track_listing", "<b>%d</b>: <i>%s</i>, added by <b>%s</b>.<br>")

	viper.SetDefault("commands.move.aliases", []string{"move", "m"})
	viper.SetDefault("commands.move.is_admin", true)
	viper.SetDefault("commands.move.description", "Moves the bot into the Mumble channel provided via argument.")
	viper.SetDefault("commands.move.messages.no_channel_provided_error", "A destination channel must be supplied to move the bot.")
	viper.SetDefault("commands.move.messages.channel_doesnt_exist_error", "The provided channel does not exist.")
	viper.SetDefault("commands.move.messages.move_successful", "You have successfully moved the bot to <b>%s</b>.")

	viper.SetDefault("commands.nexttrack.aliases", []string{"nexttrack", "nextsong", "next"})
	viper.SetDefault("commands.nexttrack.is_admin", false)
	viper.SetDefault("commands.nexttrack.description", "Outputs information about the next track in the queue if one exists.")
	viper.SetDefault("commands.nexttrack.messages.current_track_only_error", "The current track is the only track in the queue.")
	viper.SetDefault("commands.nexttrack.messages.next_track", "The next track is <i>%s</i>, added by <b>%s</b>.")

	viper.SetDefault("commands.numcached.aliases", []string{"numcached", "nc"})
	viper.SetDefault("commands.numcached.is_admin", true)
	viper.SetDefault("commands.numcached.description", "Outputs the number of tracks cached on disk if caching is enabled.")
	viper.SetDefault("commands.numcached.messages.num_cached", "There are currently <b>%d</b> items stored in the cache.")

	viper.SetDefault("commands.numtracks.aliases", []string{"numtracks", "numsongs", "nt"})
	viper.SetDefault("commands.numtracks.is_admin", false)
	viper.SetDefault("commands.numtracks.description", "Outputs the number of tracks currently in the queue.")
	viper.SetDefault("commands.numtracks.messages.one_track", "There is currently <b>1</b> track in the queue.")
	viper.SetDefault("commands.numtracks.messages.plural_tracks", "There are currently <b>%d</b> tracks in the queue.")

	viper.SetDefault("commands.pause.aliases", []string{"pause"})
	viper.SetDefault("commands.pause.is_admin", false)
	viper.SetDefault("commands.pause.description", "Pauses audio playback.")
	viper.SetDefault("commands.pause.messages.no_audio_error", "Either the audio is already paused, or there are no tracks in the queue.")
	viper.SetDefault("commands.pause.messages.paused", "<b>%s</b> has paused audio playback.")

	viper.SetDefault("commands.reload.aliases", []string{"reload", "r"})
	viper.SetDefault("commands.reload.is_admin", true)
	viper.SetDefault("commands.reload.description", "Reloads the configuration file.")
	viper.SetDefault("commands.reload.messages.reloaded", "The configuration file has been successfully reloaded.")

	viper.SetDefault("commands.reset.aliases", []string{"reset", "re"})
	viper.SetDefault("commands.reset.is_admin", true)
	viper.SetDefault("commands.reset.description", "Resets the queue by removing all queue items.")
	viper.SetDefault("commands.reset.messages.queue_reset", "<b>%s</b> has reset the queue.")

	viper.SetDefault("commands.resume.aliases", []string{"resume"})
	viper.SetDefault("commands.resume.is_admin", false)
	viper.SetDefault("commands.resume.description", "Resumes audio playback.")
	viper.SetDefault("commands.resume.messages.audio_error", "Either the audio is already playing, or there are no tracks in the queue.")
	viper.SetDefault("commands.resume.messages.resumed", "<b>%s</b> has resumed audio playback.")

	viper.SetDefault("commands.setcomment.aliases", []string{"setcomment", "comment", "sc"})
	viper.SetDefault("commands.setcomment.is_admin", true)
	viper.SetDefault("commands.setcomment.description", "Sets the comment displayed next to MumbleDJ's username in Mumble.")
	viper.SetDefault("commands.setcomment.messages.comment_removed", "The comment for the bot has been successfully removed.")
	viper.SetDefault("commands.setcomment.messages.comment_changed", "The comment for the bot has been successfully changed to the following: %s")

	viper.SetDefault("commands.shuffle.aliases", []string{"shuffle", "shuf", "sh"})
	viper.SetDefault("commands.shuffle.is_admin", true)
	viper.SetDefault("commands.shuffle.description", "Randomizes the tracks currently in the queue.")
	viper.SetDefault("commands.shuffle.messages.not_enough_tracks_error", "There are not enough tracks in the queue to execute a shuffle.")
	viper.SetDefault("commands.shuffle.messages.shuffled", "The audio queue has been shuffled.")

	viper.SetDefault("commands.skip.aliases", []string{"skip", "s"})
	viper.SetDefault("commands.skip.is_admin", false)
	viper.SetDefault("commands.skip.description", "Places a vote to skip the current track.")
	viper.SetDefault("commands.skip.messages.already_voted_error", "You have already voted to skip this track.")
	viper.SetDefault("commands.skip.messages.voted", "<b>%s</b> has voted to skip the current track.")

	viper.SetDefault("commands.skipplaylist.aliases", []string{"skipplaylist", "sp"})
	viper.SetDefault("commands.skipplaylist.is_admin", false)
	viper.SetDefault("commands.skipplaylist.description", "Places a vote to skip the current playlist.")
	viper.SetDefault("commands.skipplaylist.messages.no_playlist_error", "The current track is not part of a playlist.")
	viper.SetDefault("commands.skipplaylist.messages.already_voted_error", "You have already voted to skip this playlist.")
	viper.SetDefault("commands.skipplaylist.messages.voted", "<b>%s</b> has voted to skip the current playlist.")

	viper.SetDefault("commands.toggleshuffle.aliases", []string{"toggleshuffle", "toggleshuf", "togshuf", "tsh"})
	viper.SetDefault("commands.toggleshuffle.is_admin", true)
	viper.SetDefault("commands.toggleshuffle.description", "Toggles automatic track shuffling on/off.")
	viper.SetDefault("commands.toggleshuffle.messages.toggled_off", "Automatic shuffling has been toggled off.")
	viper.SetDefault("commands.toggleshuffle.messages.toggled_on", "Automatic shuffling has been toggled on.")

	viper.SetDefault("commands.version.aliases", []string{"version"})
	viper.SetDefault("commands.version.is_admin", false)
	viper.SetDefault("commands.version.description", "Outputs the current version of MumbleDJ.")
	viper.SetDefault("commands.version.messages.version", "MumbleDJ version: <b>%s</b>")

	viper.SetDefault("commands.volume.aliases", []string{"volume", "vol", "v"})
	viper.SetDefault("commands.volume.is_admin", false)
	viper.SetDefault("commands.volume.description", "Changes the volume if an argument is provided, outputs the current volume otherwise.")
	viper.SetDefault("commands.volume.messages.parsing_error", "The requested volume could not be parsed.")
	viper.SetDefault("commands.volume.messages.out_of_range_error", "Volumes must be between the values <b>%.2f</b> and <b>%.2f</b>.")
	viper.SetDefault("commands.volume.messages.current_volume", "The current volume is <b>%.2f</b>.")
	viper.SetDefault("commands.volume.messages.volume_changed", "<b>%s</b> has changed the volume to <b>%.2f</b>.")
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
