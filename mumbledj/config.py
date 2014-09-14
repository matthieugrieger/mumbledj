#----------------------#
#       MumbleDJ       #
#  By Matthieu Grieger #
#----------------------#-------------#
# config.py                          #
# Configuration options for the bot. #
#------------------------------------#

# ------------------------
# CONNECTION CONFIGURATION
# ------------------------

# Server address
# DEFAULT VALUE: 'localhost'
SERVER_ADDRESS = 'localhost'

# Server port
# DEFAULT VALUE: 64738
SERVER_PORT = 64738

# Username (this will be the username of the bot as well)
# DEFAULT VALUE: 'MumbleDJ'
SERVER_USERNAME = 'MumbleDJ'

# Server password (leave blank if no password exists)
# DEFAULT VALUE: ''
SERVER_PASSWORD = ''


# ---------------------
# GENERAL CONFIGURATION
# ---------------------

# Default channel
# DEFAULT VALUE: 'MumbleDJ'
DEFAULT_CHANNEL = 'Bot Testing'

# Debugging mode (True = on, False = off)
# DEFAULT VALUE: False
DEBUG = False

# Command prefix (this is the character that designates a command)
# NOTE: This must be one character!
# DEFAULT VALUE: '!'
COMMAND_PREFIX = '!'


# -------------------
# ADMIN CONFIGURATION
# -------------------

# Enable/disable admin-only commands
# DEFAULT VALUE: True
ENABLE_ADMIN_ONLY_COMMANDS = True

# List of approved ADMINS. Add all usernames who should receive admin
# privileges here.
# NOTE: I recommend only adding admins who are registered users on your server.
# Otherwise other people can use the username and get access to the admin commands.
# EXAMPLE:
# 	ADMINS = ['matthieu', 'matt']
ADMINS = ['Matt', 'DrumZ']

# Make start command admin-only
# DEFAULT VALUE: False
START_ADMIN_ONLY = False

# Make play command admin-only
# DEFAULT VALUE: False
PLAY_ADMIN_ONLY = False

# Make pause command admin-only
# DEFAULT VALUE: False
PAUSE_ADMIN_ONLY = False

# Make add command admin-only
# DEFAULT VALUE: False
ADD_ADMIN_ONLY = False

# Make skip command admin-only
# DEFAULT VALUE: False
SKIP_ADMIN_ONLY = False

# Make volumeup command admin-only
# DEFAULT VALUE: True
VOLUMEUP_ADMIN_ONLY = True

# Make volumedown command admin-only
# DEFAULT VALUE: True
VOLUMEDOWN_ADMIN_ONLY = True

# Make move command admin-only
# DEFAULT VALUE: True
MOVE_ADMIN_ONLY = True

# Make kill command admin-only (I really don't recommend changing this to False...)
# DEFAULT VALUE: True
KILL_ADMIN_ONLY = True


# ---------------------
# STORAGE CONFIGURATION
# ---------------------

# Delete audio files after they have been played.
# DEFAULT VALUE: True
DELETE_AUDIO = True

# Delete thumbnails after they have been used.
# DEFAULT VALUE: True
DELETE_THUMBNAILS = True


# ------------------
# CHAT CONFIGURATION
# ------------------

# Enable/disable chat notifications
# DEFAULT VALUE: True
SHOW_CHAT_NOTIFICATIONS = True

# Enable/disable YouTube thumbnails (only has an effect if SHOW_CHAT_NOTIFICATIONS is True)
# DEFAULT VALUE: True
SHOW_YT_THUMBNAILS = True


# -------------------
# AUDIO CONFIGURATION
# -------------------

# Bitrate
# DEFAULT VALUE: 48000
BITRATE = 48000

# Number of users that, if reached, will pause the music until it is started again by a user.
# This is to prevent against YouTube audio downloads when nobody is listening.
# DEFAULT VALUE: 1
USER_SOUND_PAUSE_TARGET = 1
