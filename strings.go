/*
 * MumbleDJ
 * By Matthieu Grieger
 * strings.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

// Message shown to users when they do not have permission to execute a command.
const NO_PERMISSION_MSG = "You do not have permission to execute that command."

// Message shown to users when they try to execute a command that doesn't exist.
const COMMAND_DOESNT_EXIST_MSG = "The command you entered does not exist."

// Message shown to users when they try to move the bot to a non-existant channel.
const CHANNEL_DOES_NOT_EXIST_MSG = "The channel you specified does not exist."

// Message shown to users when they attempt to add an invalid URL to the queue.
const INVALID_URL_MSG = "The URL you submitted does not match the required format."

// Message shown to users when they attempt to perform an action on a song when
// no song is playing.
const NO_MUSIC_PLAYING_MSG = "There is no music playing at the moment."

// Message shown to users when they issue a command that requires an argument and one was not supplied.
const NO_ARGUMENT_MSG = "The command you issued requires an argument and you did not provide one."

// Message shown to users when they try to change the volume to a value outside the volume range.
const NOT_IN_VOLUME_RANGE_MSG = "Out of range. The volume must be between %g and %g."

// Message shown to users when they successfully change the volume.
const VOLUME_SUCCESS_MSG = "You have successfully changed the volume to the following: %g."

// Message shown to user when a successful configuration reload finishes.
const CONFIG_RELOAD_SUCCESS_MSG = "The configuration has been successfully reloaded."

// Message shown to a channel when a new song starts playing.
const NOW_PLAYING_HTML = `
	<table>
		<tr>
			<td align="center"><img src="%s" width=150 /></td>
		</tr>
		<tr>
			<td align="center"><b><a href="http://youtu.be/%s">%s</a> (%s)</b></td>
		</tr>
		<tr>
			<td align="center">Added by %s</td>
		</tr>
	</table>
`

// Message shown to channel when a song is added to the queue by a user.
const SONG_ADDED_HTML = `
	<b>%s</b> has voted to skip this song.
`

// Message shown to channel when a song has been skipped.
const SONG_SKIPPED_HTML = `
	The number of votes required for a skip has been met. <b>Skipping song!</b>
`

// Message shown to users when they ask for the current volume (volume command without argument)
const CUR_VOLUME_HTML = `
	The current volume is <b>%g</b>.
`
