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
SERVER_ADDRESS = 'matthieugrieger.com'

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


# ---------------------
# COMMAND CONFIGURATION
# ---------------------

# Allow users to start music queue
# DEFAULT VALUE: True
ALLOW_START = True

# Allow users to start music playback
# DEFAULT VALUE: True
ALLOW_PLAY = True

# Allow users to pause music playback
# DEFAULT VALUE: True
ALLOW_PAUSE = True

# Allow users to add music to queue
# DEFAULT VALUE: True
ALLOW_ADD = True

# Allow users to vote to skip tracks
# DEFAULT VALUE: True
ALLOW_SKIPS = True

# Allow users to raise volume
# DEFAULT VALUE: False
ALLOW_VOLUMEUP = False

# Allow users to lower volume
# DEFAULT VALUE: False
ALLOW_VOLUMEDOWN = False

# Allow users to move bot to another channel
# DEFAULT VALUE: False
ALLOW_MOVE = False

# Allow users to kill bot (this should rarely be used)
# DEFAULT VALUE: False
ALLOW_KILL = False

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
