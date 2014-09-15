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

-- Bot username
-- DEFAULT VALUE: "MumbleDJ"
config.BOT_USERNAME = "MumbleDJ"

-- Default channel
-- DEFAULT VALUE: "Bot Testing"
config.DEFAULT_CHANNEL = "Bot Testing"

-- Command prefix
-- DEFAULT VALUE: "!"
config.COMMAND_PREFIX = "!"

-- Show status output in console?
-- DEFAULT VALUE: true
config.OUTPUT = true

-- Number of users that, if reached, will pause the music until it is started again by a user.
-- This is to prevent against YouTube audio downloads when nobody is listening.
-- DEFAULT VALUE: 1
config.USER_SOUND_PAUSE_TARGET = 1


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
config.ADMINS = {}

-- Make start an admin command?
-- DEFAULT VALUE: false
config.ADMIN_START = false

-- Make play an admin command?
-- DEFAULT VALUE: false
config.ADMIN_PLAY = false

-- Make pause an admin command?
-- DEFAULT VALUE: false
config.ADMIN_PAUSE = false

-- Make add an admin command?
-- DEFAULT VALUE: false
config.ADMIN_ADD = false

-- Make skip an admin command?
-- DEFAULT VALUE: false
config.ADMIN_SKIP = false

-- Make volumeup an admin command?
-- DEFAULT VALUE: true
config.ADMIN_VOLUMEUP = true

-- Make volumedown an admin command?
-- DEFAULT VALUE: true
config.ADMIN_VOLUMEDOWN = true

-- Make move an admin command?
-- DEFAULT VALUE: true
config.ADMIN_MOVE = true

-- Make kill an admin command?
-- DEFAULT VALUE: true (I recommend never changing this to false)
config.ADMIN_KILL = true


-------------------------
-- STORAGE CONFIGURATION
-------------------------

-- Delete audio files after they have been played?
-- DEFAULT VALUE: true
config.DELETE_AUDIO = true

-- Delete thumbnails after they have been used?
config.DELETE_THUMBNAILS = true


----------------------
-- CHAT CONFIGURATION
----------------------

-- Enable/disable chat notifications for new songs (true = on, false = off)
-- DEFAULT VALUE: true
config.SHOW_NOTIFICATIONS = true

-- Enable/disable YouTube thumbnails (true = on, false = off)
-- DEFAULT VALUE: true
config.SHOW_THUMBNAILS = true


return config
