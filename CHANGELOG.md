MumbleDJ Changelog
==================

### June 17, 2016 -- `v2.10.0`
* Added `!joinme` command (thanks [@azlux](https://github.com/azlux)).

### May 24, 2016 -- `v2.9.1`
* Fixed player command configuration setting not being applied to youtube-dl calls.

### April 9, 2016 -- `v2.9.0`
* Added support for Mixcloud (thanks [@benklett](https://github.com/benklett)).

### February 14, 2016 -- `v2.8.15`
* Fixed an incorrectly formatted error message (thanks [@GabrielPlassard](https://github.com/GabrielPlassard)).

### February 12, 2016 -- `v2.8.14`
* Audio is now downloaded using the `bestaudio` format. This prevents situations in which some audio would not play because an `.m4a` file was not available (thanks [@mpacella88](https://github.com/mpacella88)).

### February 6, 2016 -- `v2.8.13`
* Added `!version` command to display the version of MumbleDJ that is running (thanks [@zeblau](https://github.com/zeblau)).
* Added `version` commandline argument to display the version of MumbleDJ that is running (thanks [@zeblau](https://github.com/zeblau)).

### January 26, 2016 -- `v2.8.12`
* Temporarily fixed discontinued code.google.com imports.

### January 14, 2016 -- `v2.8.11`
* Fixed: Unable to use offsets if it's formatted as &t vs ?t in the URL (thanks [@fiveofeight](https://github.com/fiveofeight)).

### January 11, 2016 -- `v2.8.10`
* Created a new configuration value in the General section called PlayerCommand. This allows the user to change between "ffmpeg" and "avconv" for playing audio files.
* Added check for valid PlayerCommand value. If the value is invalid the bot will default to `ffmpeg`.

### December 26, 2015 -- `v2.8.9`
* Fixed an incorrect `!currentsong` message for songs within playlists.

### December 21, 2015 -- `v2.8.8`
* Fixed a typo in song list HTML (thanks [@mkody](https://github.com/mkody)).

### December 19, 2015 -- `v2.8.7`
* Added AnnounceNewTracks config option (thanks [@HowIChrgeLazer](https://github.com/HowIChrgeLazer)).

### December 16, 2015 -- `v2.8.6`
* Added !addnext command (thanks [@nkhoit](https://github.com/nkhoit)).
* Added argument to !listsongs command to specify how many songs to list (thanks [@nkhoit](https://github.com/nkhoit)).

### December 14, 2015 -- `v2.8.5`
* Added !listsongs command (thanks [@nkhoit](https://github.com/nkhoit)).

### December 7, 2015 -- `v2.8.4`
* YouTube and SoundCloud API keys are now stored in the configuration file instead of environment variables. Existing installations with API keys in environment variables will automatically be migrated to the configuration file (thanks [@Gamah](https://github.com/Gamah)).

### October 16, 2015 -- `v2.8.3`
* Playlists can now be over 50 songs in length (thanks [@GabrielPlassard](https://github.com/GabrielPlassard)).
* Added MaxSongPerPlaylist configuration option.

### October 14, 2015 -- `v2.8.2`
* Fixed possible index out of range panic when auto shuffle is on (thanks [@GabrielPlassard](https://github.com/GabrielPlassard)).

### October 12, 2015 -- `v2.8.1`
* Added !shuffle, !shuffleon, and !shuffleoff commands (thanks [@GabrielPlassard](https://github.com/GabrielPlassard)).

### October 1, 2015 -- `v2.8.0`
* Added Soundcloud support (thanks [@MichaelOultram](https://github.com/MichaelOultram)).

### August 12, 2015 -- `v2.7.5`
* Fixed cache clearing earlier than expected (thanks [@CMahaff](https://github.com/CMahaff)).

### May 19, 2015 -- `v2.7.4`
* Fixed a panic that occurred when certain YouTube playlists were added to the queue.

### May 14, 2015 -- `v2.7.3`
* Fixed `!move` not working for subchannels (thanks [@mkbwong](https://github.com/mkbwong)).
* Fixed MumbleDJ showing invalid YouTube ID error message in chat when an invalid YouTube API key is supplied (thanks [@fiveofeight](https://github.com/fiveofeight)).
* Fixed MumbleDJ showing invalid YouTube ID error message in chat when a song exceeds the allowed time duration.

### May 12, 2015 -- `v2.7.2`
* Fixed incorrect values shown in timestamp for videos over an hour long.
* Reworked timestamp parsing.

### May 9, 2015 -- `v2.7.1`
* Added support for YouTube offsets. This means that YouTube URLs with the `t` parameter will start at the time specified in the URL instead of the beginning.
* Cleaned up comments in some files and removed some unnecessary code.
* Fixed a bug in which a duration of 0:00 was shown for songs that were less than a minute long.

### April 17, 2015 -- `v2.7.0`
* Migrated all YouTube API calls to YouTube Data API v3. This means that you **MUST** follow the instructions in the following link if you were using a previous version of MumbleDJ: https://github.com/matthieugrieger/mumbledj#youtube-api-keys.
* Made the SongQueue much more flexible. These changes will allow easy addition of support for other music services.

### March 28, 2015 -- `v2.6.10`
* Fixed a crash that would occur when the last song of a playlist was skipped.

### March 27, 2015 -- `v2.6.9`
* Fixed a race condition that would sometimes cause the bot to crash (thanks [dylanetaft](https://github.com/dylanetaft)!).

### March 26, 2015 -- `v2.6.8`
* Renamed `mumbledj.gcfg` to `config.gcfg`. However, please note that it will still be called `mumbledj.gcfg` in your `~/.mumbledj` directory. Hopefully this will avoid any ambiguity when referring to the
config files.
* Tweaked the `Makefile` to handle situations where `go build` creates an executable with an appended version number.

### March 20, 2015 -- `v2.6.7`
* Fixed a typo in `mumbledj.gcfg`.
* Songs and playlists are now skipped immediately if the submitter submits a skip command.
* `SONG_SKIPPED_HTML` and `PLAYLIST_SKIPPED_HTML` are no longer shown if the submitter or admin skips a song/playlist.

### March 7, 2015 -- `v2.6.6`
* Added missing AdminSkipPlaylistAlias option to `mumbledj.gcfg`.

### February 25, 2015 -- `v2.6.5`
* Added automatic connection retries if the bot loses connection to the server. The bot will attempt to reconnect to the server every 30 seconds for a period of 15 minutes, then exit if a connection cannot be made.

### February 20, 2015 -- `v2.6.4`
* Fixed failed audio downloads for YouTube videos with IDs beginning with "-".

### February 19, 2015 -- `v2.6.3`
* Added `gumbleutil.CertificateLockFile()` for more secure connections.
* Added `-insecure` boolean commandline flag to allow MumbleDJ to connect to a server without overwriting `~/.mumbledj/cert.lock`.

### February 18, 2015 -- `v2.6.2`
* Fixed bot crashing after 5 minutes if there is nothing in the song queue.
* Fixed queue freezing up if the download of the first song in queue fails.

### February 17, 2015 -- `v2.6.0, v2.6.1`
* Added caching system to MumbleDJ.
* Added configuration variables in `mumbledj.gcfg` for caching related settings (please note that caching is off by default).
* Added `!numcached` and `!cachesize` commands for admins.
* Added optional song length limit (thanks [jakexks](https://github.com/jakexks)!)

### February 12, 2015 -- `v2.5.0`
* Updated dependencies and fixed code to match `gumble` API changes.
* Greatly simplified the song queue data structure. Some new bugs could potentially have arisen. Let me know if you find any!

### February 9, 2015 -- `v2.4.3`
* Added configuration option in `mumbledj.gcfg` for default bot comment.
* Fixed text messages only containing images crashing the bot.

### February 7, 2015 -- `v2.4.2`
* Updated `gumble` and `gumbleutil` dependencies.
* Removed `sanitize` dependency.
* Reworked `Makefile` slightly.
* Now uses `gumbleutil.PlainText` for removing HTML tags instead of `sanitize`.
* Added `!setcomment` which allows admin users to set the comment for the bot.
* Made "Now Playing" notification and `!currentsong` show the playlist title of the song if it is included in a playlist.
* Added ability to connect to Mumble server using a PEM cert/key pair. Use the commandline flags `cert` and `key` to make use of this.
* Added an easier to read error message upon unsuccessful connection to server.

### February 3, 2015 -- `v2.4.1`
* Made it possible to place MumbleDJ binary in `~/bin` instead of `/usr/local/bin` if the folder exists.

### February 2, 2015 -- `v2.3.4, v2.3.5, v2.3.6, v2.3.7, v2.4.0`
* Added panic on audio play fail for debugging purposes.
* Fixed '!' being recognized as '!skipplaylist'.
* Fixed !reset crash when there is no audio playing.
* Fixed newlines after YouTube URL messing up !add commands.
* Fixed empty song/playlist entry being added upon !add with invalid YouTube ID.
* Fixed go build issues.
* Added `goop` dependency management. Make sure you have `openal` installed, or it won't work right!
* Fixed crash on invalid playlist URL.

### January 30, 2015 -- `v2.3.3`
* Fixed private messages crashing the bot when the target user switches channels or disconnects.

### January 26, 2015 -- `v2.3.2`
* Fixed !nextsong showing incorrect information about the next song in the queue.

### January 25, 2015 -- `v2.3.0, v2.3.1`
* Added !currentsong command, which displays information about the song currently playing.
* MumbleDJ now removes disconnected users from skiplists for playlists and songs within the SongQueue.
* Fixed crash when a user disconnects when no song is playing.

### January 19, 2015 -- `v2.2.11`
* Fixed not being able to use the move command with channels with spaces in their name.

### January 14, 2015 -- `v2.2.9, v2.2.10`
* Set AudioEncoder Application to `gopus.Audio` instead of `gopus.Voice` for hopefully better sound quality.
* Added some commands to the !help string that were missing.

### January 12, 2015 -- `v2.2.8`
* Added !nextsong command, which outputs some information about the next song in the queue if it exists.

### January 10, 2015 -- `v2.2.6, v2.2.7`
* Fixed type mismatch error when building MumbleDJ.
* Added traversal function to SongQueue.
* Added !numsongs command, which outputs how many songs are currently in the SongQueue.
* Added !help command, which displays a list of valid commands in Mumble chat.

### January 9, 2015 -- `v2.2.5`
* Fixed some YouTube playlist URLs crashing the bot and not retrieving metadata correctly.

### January 8, 2015 -- `v2.2.4`
* Fixed a crash caused by a user trying to skip the same song more than once.

### January 7, 2015 -- `v2.2.3`
* Fixed a crash caused by entering a skip request when no song is currently playing.

### January 5, 2015 -- `v2.2.1, v2.2.2`
* Attached `gumbleutil.AutoBitrate` EventListener to client. This should hopefully fix the issues with audio cutting in and out.
* Moved dependency installation to default `make` command to better enforce new updates.
* Added `make build` to `Makefile` to allow previous functionality of the default `make` command.
* Hopefully fixed a situation that would cause the song queue to stop working.
* Added `!reset` command. Use this to reset the song queue.

### January 3, 2015 -- `v2.2.0`
* Added ability to add YouTube playlists to the queue. Note that the max size of a playlist is 25 songs, anything larger will only use the first 25 songs in the playlist.
* Fixed a crash while attempting to add URLs to the queue.
* Re-made the song queue using my own "queue-like" structure using slices.

### December 30, 2014 -- `v2.1.3`
* Fixed YouTube URL parsing not working for some forms of YouTube URLs.
* Now recovers more gracefully if an audio download fails. Instead of panicking, the bot will send a message to the user who added the URL, telling them the audio download failed.

### December 29, 2014 -- `v2.1.2`
* Fixed skip messages not being displayed in chat.

### December 27, 2014 -- `v2.0.0, v2.1.0, v2.1.1`
* Reached feature parity with old version of MumbleDJ.
* Bot is now written completely in Golang instead of Lua and Python.
* Now uses [`gumble`](https://github.com/layeh/gumble) for interacting with Mumble instead of [`piepan`](https://github.com/layeh/piepan).
* Stability improved in many areas.
* Audio quality is slightly better due to using higher bitrate m4a files instead of Ogg Vorbis.
* All YouTube URLs should be supported now.
* Added an admin skip command that allows an admin to force skip a song.
* Added `mumbledj.gcfg`, where all configuration options are now stored.
* Added a reload command, used to reload the configuration when a change is made.
* Implemented volume control. Now changes volume while audio is playing!
* Code is now more thoroughly commented.
* Fixed char comparison with dj.conf.General.CommandPrefix.

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
