MumbleDJ Changelog
==================

### December 27, 2014 -- `v2.0.0`
* Reached feature parity with old version of MumbleDJ.
* Bot is now written completely in Golang instead of Lua and Python.
* Now uses [`gumble`](https://github.com/layeh/gumble) for interacting with Mumble instead of [`piepan`](https://github.com/layeh/piepan).
* Stability improved in many areas.
* Audio quality is slightly better due to using higher bitrate m4a files instead of Ogg Vorbis.
* All YouTube URLs should be supported now.
* Added an admin skip command that allows an admin to force skip a song.
* Added `mumbledj.gcfg`, where all configuration options are now stored.
* Added a reload command, used to reload the configuration when a change is made.

### December 8, 2014
* Switched from Ruby to Go, using `gumble` instead of `mumble-ruby` now.

### November 15, 2014
* Created "v2" branch for Ruby rewrite.

### November 9, 2014
* Fixed volume changed message showing wrong value.

### October 24, 2014
* Switched volume change method. The volume is now changed directly through `piepan` instead of `ffmpeg`.
* Fixed another bug with volume changing.

### October 23, 2014
* Fixed a bug that would not allow audio encoding to complete successfully while on Debian.
* Fixed a stupid typo that broke the `!volume` command.
* Updated `SETUP.md` with instructions on installing MumbleDJ on Debian.
* Added missing commands in `SETUP.md`.

### October 18, 2014
* Fixed a crash when an error occurs during the audio downloading & encoding process.
* Fixed a crash that occurs when the bot tries to join a default channel that does not exist. If the default channel does not exist, the bot will just move itself 
to the root of the server instead.

### October 13, 2014
* Added `SETUP.md` which contains a guide on how to configure MumbleDJ and install its dependencies.
* Deleted song_queue.lua and moved all contents to mumbledj.lua. In the end this will make the script simpler.
* Fixed song skipping.

### October 7, 2014
* Made user skip message show even when the target number of skips has been reached.
* Made "Music" the default Mumble channel.

### September 30, 2014
* Fixed skips not working correctly.
* Fixed a crash related to private messages.

### September 26, 2014
* Removed play and pause commands. There were issues with these commands, and they both serve functions that can be done within Mumble per-user.
* Removed all play and pause configuration options from config.lua.
* Updated .gitignore to ignore .ogg files.

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
