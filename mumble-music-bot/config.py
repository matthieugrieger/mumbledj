#----------------------#
#   Mumble Music Bot   #
#  By Matthieu Grieger #
#----------------------#-------------#
# config.py                          #
# Configuration options for the bot. #
#------------------------------------#

# CONNECTION CONFIGURATION
# ------------------------
# Server address
SERVER_ADDRESS = 'localhost'
# Server port (64738 is the default)
SERVER_PORT = '64738'
# Username (this will be the username of the bot as well)
SERVER_USERNAME = 
# Server password (leave blank if no password exists)
SERVER_PASSWORD = ''

# GENERAL CONFIGURATION
# ---------------------
# Default channel
DEFAULT_CHANNEL = 'Bot Testing'

# AUDIO CONFIGURATION
# -------------------
# Bitrate
BITRATE = 48000
# Number of users that, if reached, will pause the music until it is started again by a user.
# This is to prevent against YouTube audio downloads when nobody is listening.
USER_SOUND_PAUSE_LIMIT = 1
