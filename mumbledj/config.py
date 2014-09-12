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
SERVER_ADDRESS = 'localhost'

# Server port (64738 is the default)
SERVER_PORT = 64738

# Username (this will be the username of the bot as well)
SERVER_USERNAME = 'MumbleDJ'

# Server password (leave blank if no password exists)
SERVER_PASSWORD = ''


# ---------------------
# GENERAL CONFIGURATION
# ---------------------

# Default channel
DEFAULT_CHANNEL = 'Bot Testing'

# Debugging mode (True = on, False = off)
DEBUG = False

# Command prefix (this is the character that designates a command)
COMMAND_PREFIX = '!'


# ------------------
# CHAT CONFIGURATION
# ------------------

# Enable/disable chat notifications
SHOW_CHAT_NOTIFICATIONS = True

# Enable/disable YouTube thumbnails (only has an effect if SHOW_CHAT_NOTIFICATIONS is True)
SHOW_YT_THUMBNAILS = True


# -------------------
# AUDIO CONFIGURATION
# -------------------

# Bitrate
BITRATE = 48000

# Number of users that, if reached, will pause the music until it is started again by a user.
# This is to prevent against YouTube audio downloads when nobody is listening.
USER_SOUND_PAUSE_TARGET = 1
