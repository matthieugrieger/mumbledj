MumbleDJ Changelog
==================

### September 16, 2014
* Removed volumeup/volumedown commands, replaced with just volume.
* Added deque.lua (thanks Pierre Chapuis!).
* Added song_queue.lua.

### September 15, 2014
* Added command parsing.
* Added placeholder functions for various commands.
* Added admin, message, and other miscellaneous config options.
* Removed start command.
* Added a permissions/admin system.

### September 14, 2014
* Changed the base for the project from pymumble to piepan.
* Entire codebase is now written in Lua instead of Python.
* Re-implemented some of the config in config.lua.
* Implemented code to connect the bot to the server and move it into the Bot Testing channel.

### September 13, 2014
* Added song.py, a file that houses the Song class.
* Added command & storage options to config.py.

### September 12, 2014
* mumble-music-bot repository renamed to mumbledj.
* Renamed all mentions of mumble-music-bot or musicbot to mumbledj.
* Restructured project for easier imports.
* Added .gitignore for pymumble.
* Now successfully connects to Mumble servers.
* Added command parsing.

### September 11, 2014
* mumble-music-bot repository created.
* Added config.py with some basic configuration options.
* Put placeholder methods within the MusicBot object.
* Add run_bot.py.
