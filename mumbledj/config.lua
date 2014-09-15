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

return config
