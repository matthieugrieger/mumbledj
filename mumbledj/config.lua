-------------------------
--      MumbleDJ       --
-- By Matthieu Grieger --
-------------------------------------------------
-- config.lua                                  --
-- This is where all the configuration options --
-- for the bot can be set.                     --
-------------------------------------------------

local config = {}
-------------------------
-- GENERAL CONFIGURATION
-------------------------

-- Default channel
-- DEFAULT VALUE: "Music"
config.DEFAULT_CHANNEL = "Music"

-- Command prefix
-- DEFAULT VALUE: "!"
config.COMMAND_PREFIX = "!"

-- Show status output in console?
-- DEFAULT VALUE: true
config.OUTPUT = true

-- Default volume (256 being normal volume)
-- DEFAULT VALUE: 32
config.VOLUME = 32

-- Lowest volume allowed
-- DEFAULT VALUE: 16
config.LOWEST_VOLUME = 16

-- Highest volume allowed
-- DEFAULT VALUE: 512
config.HIGHEST_VOLUME = 512

-- Ratio that must be met or exceeded to trigger a song skip.
-- DEFAULT VALUE: 0.5
config.SKIP_RATIO = 0.5


-------------------------
-- COMMAND CONFIGURATION
-------------------------

-- Alias used for add command.
-- DEFAULT VALUE: "add"
config.ADD_ALIAS = "add"

-- Alias used for skip command.
-- DEFAULT VALUE: "skip"
config.SKIP_ALIAS = "skip"

-- Alias used for volume command.
-- DEFAULT VALUE: "volume"
config.VOLUME_ALIAS = "volume"

-- Alias used for move command.
-- DEFAULT VALUE: "move"
config.MOVE_ALIAS = "move"

-- Alias used for kill command.
-- DEFAULT VALUE: "kill"
config.KILL_ALIAS = "kill"


-----------------------
-- ADMIN CONFIGURATION
-----------------------

-- Enable admins (true = on, false = off)
-- DEFAULT VALUE: true
config.ENABLE_ADMINS = true

-- List of admins
-- NOTE: I recommend only giving users admin privileges if they are registered
-- on the server. Otherwise people can just take their username and issue admin
-- commands.
-- EXAMPLE:
-- 	config.ADMINS = {"Matt", "Matthieu"}
config.ADMINS = {"Matt"}

-- Make add an admin command?
-- DEFAULT VALUE: false
config.ADMIN_ADD = false

-- Make skip an admin command?
-- DEFAULT VALUE: false
config.ADMIN_SKIP = false

-- Make volume an admin command?
-- DEFAULT VALUE: true
config.ADMIN_VOLUME = true

-- Make move an admin command?
-- DEFAULT VALUE: true
config.ADMIN_MOVE = true

-- Make kill an admin command?
-- DEFAULT VALUE: true (I recommend never changing this to false)
config.ADMIN_KILL = true


----------------------
-- CHAT CONFIGURATION
----------------------

-- Enable/disable chat notifications for new songs (true = on, false = off)
-- DEFAULT VALUE: true
config.SHOW_NOTIFICATIONS = true

-------------------------
-- MESSAGE CONFIGURATION
-------------------------

-- Message shown to users when they do not have permission to execute a command.
config.NO_PERMISSION_MSG = "You do not have permission to execute that command."

-- Message shown to users when they try to move the bot to a non-existant channel.
config.CHANNEL_DOES_NOT_EXIST_MSG = "The channel you specified does not exist."

-- Message shown to users when they attempt to add an invalid URL to the queue.
config.INVALID_URL_MSG = "The URL you submitted does not match the required format. Either you did not provide a YouTube URL, or an error occurred during the downloading & encoding process."

-- Message shown to users when they attempt to use the stop command when no music is playing.
config.NO_MUSIC_PLAYING_MSG = "There is no music playing at the moment."

-- Message shown to users when they issue a command that requires an argument and one was not supplied.
config.NO_ARGUMENT = "The command you issued requires an argument and you did not provide one. Make sure a space exists between the command and the argument."

-- Message shown to users when they try to change the volume to a value outside the volume range.
config.NOT_IN_VOLUME_RANGE = "The volume you tried to supply is not in the allowed volume range. The value must be between " .. config.LOWEST_VOLUME .. " and " .. config.HIGHEST_VOLUME .. "."


----------------------
-- HTML CONFIGURATION
----------------------

-- Message shown to channel when a new song starts playing.
config.NOW_PLAYING_HTML = [[
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
]]

-- Message shown to channel when a song is added to the queue by a user.
config.SONG_ADDED_HTML = [[
	<b>%s</b> has added "%s" to the queue.
]]

-- Message shown to channel when a user votes to skip a song.
config.USER_SKIP_HTML = [[
	<b>%s</b> has voted to skip this song.
]]

-- Message shown to channel when a song has been skipped.
config.SONG_SKIPPED_HTML = [[
	The number of votes required for a skip has been met. <b>Skipping song!</b>
]]

return config
