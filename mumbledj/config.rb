# MumbleDJ
# By Matthieu Grieger
# config.rb


# ------------------------
# CONNECTION CONFIGURATION
# ------------------------

# Bot username
# DEFAULT VALUE: "MumbleDJ"
BOT_USERNAME = "MumbleDJTest"

# Password to join Mumble server
# DEFAULT VALUE: "" (leave it as this value if no password is required)
MUMBLE_PASSWORD = ENV['MUMBLE_PW']

# Server address
# DEFAULT VALUE: "localhost"
MUMBLE_SERVER_ADDRESS = "matthieugrieger.com"

# Server port number
# DEFAULT VALUE: 64738
MUMBLE_SERVER_PORT = 64738


# ---------------------
# GENERAL CONFIGURATION
# ---------------------

# Default channel
# DEFAULT VALUE: "Music"
DEFAULT_CHANNEL = "Bot Testing"

# Command prefix
# DEFAULT VALUE: "!"
COMMAND_PREFIX = "!"

# Show status output in console?
# DEFAULT VALUE: true
OUTPUT_ENABLED = true

# Default volume
# DEFAULT VALUE: 0.2
VOLUME = 0.2

# Lowest volume allowed
# DEFAULT VALUE: 0.01
LOWEST_VOLUME = 0.01

# Highest volume allowed
# DEFAULT VALUE: 0.6
HIGHEST_VOLUME = 0.6

# Ratio that must be met or exceeded to trigger a song skip
# DEFAULT VALUE: 0.5
SKIP_RATIO = 0.5


# ---------------------
# COMMAND CONFIGURATION
# ---------------------

# Alias used for add command
# DEFAULT VALUE: "add"
ADD_ALIAS = "add"

# Alias used for skip command
# DEFAULT VALUE: "skip"
SKIP_ALIAS = "skip"

# Alias used for volume command
# DEFAULT VALUE: "volume"
VOLUME_ALIAS = "volume"

# Alias used for move command
# DEFAULT VALUE: "move"
MOVE_ALIAS = "move"

# Alias used for mute command
# DEFAULT VALUE: "mute"
MUTE_ALIAS = "mute"

# Alias used for unmute command
# DEFAULT VALUE: "unmute"
UNMUTE_ALIAS = "unmute"


# -------------------
# ADMIN CONFIGURATION
# -------------------

# Enable admins (true = on, false = off)
# DEFAULT VALUE: true
ENABLE_ADMINS = true

# List of admins
# NOTE: I recommend only giving users admin privileges if they are
# registered on the server. Otherwise people can just take their username
# and issue admin commands.
ADMINS = ["DrumZ"]

# Make add an admin command?
# DEFAULT VALUE: false
ADMIN_ADD = false

# Make skip an admin command?
# DEFAULT VALUE: false
ADMIN_SKIP = false

# Make volume an admin command?
# DEFAULT VALUE: true
ADMIN_VOLUME = true

# Make move an admin command?
# DEFAULT VALUE: true
ADMIN_MOVE = true

# Make mute an admin command?
# DEFAULT VALUE: true
ADMIN_MUTE = true

# Make unmute an admin command?
# DEFAULT VALUE: true
ADMIN_UNMUTE = true


