/*
 * MumbleDJ
 * By Matthieu Grieger
 * strings.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

// Message shown to users when the bot has an invalid API key.
const INVALID_API_KEY = "MumbleDJ does not have a valid %s API key."

// Message shown to users when they do not have permission to execute a command.
const NO_PERMISSION_MSG = "You do not have permission to execute that command."

// Message shown to users when they try to add a playlist to the queue and do not have permission to do so.
const NO_PLAYLIST_PERMISSION_MSG = "You do not have permission to add playlists to the queue."

// Message shown to users when they try to execute a command that doesn't exist.
const COMMAND_DOESNT_EXIST_MSG = "The command you entered does not exist."

// Message shown to users when they try to move the bot to a non-existant channel.
const CHANNEL_DOES_NOT_EXIST_MSG = "The channel you specified does not exist."

// Message shown to users when they attempt to add an invalid URL to the queue.
const INVALID_URL_MSG = "The URL you submitted does not match the required format."

// Message shown to users when they attempt to add a video that's too long
const TRACK_TOO_LONG_MSG = "The %s you submitted exceeds the duration allowed by the server."

// Message shown to users when they attempt to perform an action on a song when
// no song is playing.
const NO_MUSIC_PLAYING_MSG = "There is no music playing at the moment."

// Message shown to users when they attempt to skip a playlist when there is no playlist playing.
const NO_PLAYLIST_PLAYING_MSG = "There is no playlist playing at the moment."

// Message shown to users when they attempt to use the nextsong command when there is no song coming up.
const NO_SONG_NEXT_MSG = "There are no songs queued at the moment."

// Message shown to users when they issue a command that requires an argument and one was not supplied.
const NO_ARGUMENT_MSG = "The command you issued requires an argument and you did not provide one."

// Message shown to users when they try to change the volume to a value outside the volume range.
const NOT_IN_VOLUME_RANGE_MSG = "Out of range. The volume must be between %f and %f."

// Message shown to user when a successful configuration reload finishes.
const CONFIG_RELOAD_SUCCESS_MSG = "The configuration has been successfully reloaded."

// Message shown to users when an admin skips a song.
const ADMIN_SONG_SKIP_MSG = "An admin has decided to skip the current song."

// Message shown to users when an admin skips a playlist.
const ADMIN_PLAYLIST_SKIP_MSG = "An admin has decided to skip the current playlist."

// Message shown to users when the audio for a video could not be downloaded.
const AUDIO_FAIL_MSG = "The audio download for this video failed. %s has likely not generated the audio files for this %s yet. Skipping to the next song!"

// Message shown to users when they supply an URL that does not contain a valid ID.
const INVALID_ID_MSG = "The %s URL you supplied did not contain a valid ID."

// Message shown to user when they successfully update the bot's comment.
const COMMENT_UPDATED_MSG = "The comment for the bot has successfully been updated."

// Message shown to user when they request to see the number of songs cached on disk.
const NUM_CACHED_MSG = "There are currently %d songs cached on disk."

// Message shown to user when they request to see the total size of the cache.
const CACHE_SIZE_MSG = "The cache is currently %g MB in size."

// Message shown to user when they attempt to issue a cache-related command when caching is not enabled.
const CACHE_NOT_ENABLED_MSG = "The cache is not currently enabled."

// Message shown to channel when a song is added to the queue by a user.
const SONG_ADDED_HTML = `
	<b>%s</b> has added "%s" to the queue.
`

// Message shown to channel when a playlist is added to the queue by a user.
const PLAYLIST_ADDED_HTML = `
	<b>%s</b> has added the %s "%s" to the queue.
`

// Message shown to channel when a song has been skipped.
const SONG_SKIPPED_HTML = `
	The number of votes required for a skip has been met. <b>Skipping song!</b>
`

// Message shown to channel when a playlist has been skipped.
const PLAYLIST_SKIPPED_HTML = `
	The number of votes required for a skip has been met. <b>Skipping playlist!</b>
`

// Message shown to display bot commands.
const HELP_HTML = `<br/>
	<b>User Commands:</b>
	<p><b>!help</b> - Displays this help.</p>
	<p><b>!add</b> - Adds songs/playlists to queue.</p>
	<p><b>!volume</b> - Either tells you the current volume or sets it to a new volume.</p>
	<p><b>!skip</b> - Casts a vote to skip the current song</p>
	<p> <b>!skipplaylist</b> - Casts a vote to skip over the current playlist.</p>
	<p><b>!numsongs</b> - Shows how many songs are in queue.</p>
	<p><b>!nextsong</b> - Shows the title and submitter of the next queue item if it exists.</p>
	<p><b>!currentsong</b> - Shows the title and submitter of the song currently playing.</p>
	<p style="-qt-paragraph-type:empty"><br/></p>
	<p><b>Admin Commands:</b></p>
	<p><b>!reset</b> - An admin command that resets the song queue. </p>
	<p><b>!forceskip</b> - An admin command that forces a song skip. </p>
	<p><b>!forceskipplaylist</b> - An admin command that forces a playlist skip. </p>
	<p><b>!move </b>- Moves MumbleDJ into channel if it exists.</p>
	<p><b>!reload</b> - Reloads mumbledj.gcfg configuration settings.</p>
	<p><b>!setcomment</b> - Sets the comment for the bot.</p>
	<p><b>!numcached</b></p> - Outputs the number of songs cached on disk.</p>
	<p><b>!cachesize</b></p> - Outputs the total file size of the cache in MB.</p>
	<p><b>!kill</b> - Safely cleans the bot environment and disconnects from the server.</p>
`

// Message shown to users when they ask for the current volume (volume command without argument)
const CUR_VOLUME_HTML = `
	The current volume is <b>%.2f</b>.
`

// Message shown to users when another user votes to skip the current song.
const SKIP_ADDED_HTML = `
	<b>%s</b> has voted to skip the current song.
`

// Message shown to users when the submitter of a song decides to skip their song.
const SUBMITTER_SKIP_HTML = `
	The current song has been skipped by <b>%s</b>, the submitter.
`

// Message shown to users when another user votes to skip the current playlist.
const PLAYLIST_SKIP_ADDED_HTML = `
	<b>%s</b> has voted to skip the current %s.
`

// Message shown to users when the submitter of a song decides to skip their song.
const PLAYLIST_SUBMITTER_SKIP_HTML = `
	The current %s has been skipped by <b>%s</b>, the submitter.
`

// Message shown to users when they successfully change the volume.
const VOLUME_SUCCESS_HTML = `
	<b>%s</b> has changed the volume to <b>%.2f</b>.
`

// Message shown to users when a user successfully resets the SongQueue.
const QUEUE_RESET_HTML = `
	<b>%s</b> has cleared the song queue.
`

// Message shown to users when a user asks how many songs are in the queue.
const NUM_SONGS_HTML = `
	There are currently <b>%d</b> song(s) in the queue.
`

// Message shown to users when they issue the nextsong command.
const NEXT_SONG_HTML = `
	The next song in the queue is "%s", added by <b>%s</b>.
`

// Message shown to users when they issue the currentsong command.
const CURRENT_SONG_HTML = `
	The song currently playing is "%s", added by <b>%s</b>.
`

// Message shown to users when the currentsong command is issued when a song from a
// playlist is playing.
const CURRENT_SONG_PLAYLIST_HTML = `
	The %s currently playing is "%s", added <b>%s</b> from the %s "%s".
`
