MumbleDJ Changelog
==================

### September 26, 2014
* Removed play and pause commands. There were issues with these commands, and they both serve functions that can be done within Mumble per-user.
* Removed all play and pause configuration options from config.lua.

### September 25, 2014
* Forced ffmpeg to use libvorbis codec.

### September 23, 2014
* Bot now seems to be working!
* Skipping songs works now.
* Second audio track now plays directly after the first one.
* Silenced ffmpeg output again.

### September 18, 2014
* Added command alias options to config.
* Moved most of the skip-related code to song_queue.lua.
* Commented more thoroughly the code, mostly pointing out what each function does.
* Made more progress toward a working song queue. It only seems to play the first song in the queue at the moment.

### September 17, 2014
* Removed USERNAME field from config.lua. It wasn't needed and introduced situations that may cause problems.
* Fixed download_audio.py. It now seems to reliably download/encode audio. ffmpeg output has been silenced.
* Volume is now set during encode by ffmpeg, since the volume option in piepan's play() does not seem to work.

### September 16, 2014
* Removed volumeup/volumedown commands, replaced with just volume.
* Added deque.lua (thanks Pierre Chapuis!).
* Added song_queue.lua.
* Made significant progress toward a working song system.

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
