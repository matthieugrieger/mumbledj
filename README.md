MumbleDJ
========
A Mumble bot that plays music fetched from YouTube videos.

## USAGE
`$ mumbledj -server=localhost -port=64738 -username=MumbleDJ -password="" -channel=root -cert="" -key=""`
All parameters are optional, the example above shows the default values for each field.

## COMMANDS
These are all of the chat commands currently supported by MumbleDJ. All command names and command prefixes may be changed in `~/.mumbledj/config/mumbledj.gcfg`.

Command | Description | Arguments | Admin | Example
--------|-------------|-----------|-------|--------
**add** | Adds a YouTube video's audio to the song queue. If no songs are currently in the queue, the audio will begin playing immediately. YouTube playlists may also be added using this command. Please note, however, that if a YouTube playlist contains over 25 videos only the first 25 videos will be placed in the song queue. | youtube_video_url OR youtube_playlist_url | No | `!add https://www.youtube.com/watch?v=5xfEr2Oxdys`
**skip**| Submits a vote to skip the current song. Once the skip ratio target (specified in `mumbledj.gcfg`) is met, the song will be skipped and the next will start playing. Each user may only submit one skip per song. | None | No | `!skip`
**skipplaylist** | Submits a vote to skip the current playlist. Once the skip ratio target (specified in mumbledj.gcfg) is met, the playlist will be skipped and the next song/playlist will start playing. Each user may only submit one skip per playlist. | None | No | `!skipplaylist`
**forceskip** | An admin command that forces a song skip. | None | Yes | `!forceskip`
**forceskipplaylist** | An admin command that forces a playlist skip. | None | Yes | `!forceskipplaylist`
**help** | Displays this list of commands in Mumble chat. | None | No | `!help`
**volume** | Either outputs the current volume or changes the current volume. If desired volume is not provided, the current volume will be displayed in chat. Otherwise, the volume for the bot will be changed to desired volume if it is within the allowed volume range. | None OR desired volume | No | `!volume 0.5`, `!volume`
**move** | Moves MumbleDJ into channel if it exists. | Channel | Yes | `!move Music`
**reload** | Reloads `mumbledj.gcfg` to retrieve updated configuration settings. | None | Yes | `!reload`
**reset** | Stops all audio and resets the song queue. | None | Yes | `!reset`
**numsongs** | Outputs the number of songs in the queue in chat. Individual songs and songs within playlists are both counted. | None | No | `!numsongs`
**nextsong** | Outputs the title and name of the submitter of the next song in the queue if it exists. | None | No | `!nextsong`
**currentsong** | Outputs the title and name of the submitter of the song currently playing. | None | No | `!currentsong`
**setcomment** | Sets the comment for the bot. If no argument is given, the current comment will be removed. | None OR new_comment | Yes | `!setcomment Hello! I am a bot. Type !help for the available commands.`
**kill** | Safely cleans the bot environment and disconnects from the server. Please use this command to stop the bot instead of force closing, as the kill command deletes any remaining songs in the `~/.mumbledj/songs` directory. | None | Yes | `!kill`




## INSTALLATION
Installation for v2 of MumbleDJ is much easier than it was before, due to the reduced dependency list and a `Makefile` which automates some of the process.

**NOTE:** This bot was designed for use on Linux machines. If you wish to run the bot on another OS, it will require tweaking and is not something I will be able to help with.  
**NOTE #2:** Your Mumble server MUST be using the Opus audio codec, not CELT. Audio will not play if your server uses CELT.

###SETUP GUIDE  
**1)** Install and correctly configure [`Go`](https://golang.org/) (1.3 or higher). Specifically, make sure to follow [this guide](https://golang.org/doc/code.html) and set the `GOPATH` environment variable properly.

**2)** Install [`ffmpeg`](https://www.ffmpeg.org/) and [`mercurial`](http://mercurial.selenic.com/) if they are not already installed on your system. Also be sure that you have 
[`opus`](http://www.opus-codec.org/) and its development headers installed on your system, as well as `openal` (check your distributions repo for the package name).

**3)** Install [`youtube-dl`](https://github.com/rg3/youtube-dl#installation). It is recommended to install `youtube-dl` through the method described on the linked GitHub page, rather than installing through a distribution repository. This ensures that you get the most up-to-date version of `youtube-dl`.

**4)** If you wish to install MumbleDJ without any further root privileges, make sure that `~/bin` exists and is added to your `$PATH`. If this step is not done, the `Makefile` will place the MumbleDJ binary in `/usr/local/bin` instead, which requires root privileges.

**5)** Clone the `MumbleDJ` repository or [download the latest release](https://github.com/matthieugrieger/mumbledj/releases).

**6)** `cd` into the `MumbleDJ` repository directory and execute the following commands: 
```
$ make
$ make install
```

**7)** Edit `~/.mumbledj/config/mumbledj.gcfg` to your liking. This file will be overwritten if the config file structure is changed in a commit, but a backup is always stored at 
`~/.mumbledj/config/mumbledj_backup.gcfg`.

**8)** Execute the command shown at the top of this `README` document with your credentials, and the bot should be up and running!

**Recommended, but not required:** Set `opusthreshold=0` in `/etc/mumble-server.ini` or `/etc/murmur.ini`. This will force the server to always use the Opus audio codec, which is the only codec that MumbleDJ supports.

###UPDATE GUIDE
**1)** `git pull` or [download the latest release](https://github.com/matthieugrieger/mumbledj/releases).

**2)** Issue the following commands within your updated MumbleDJ directory:
```
$ make clean
$ make
$ make install
```
**NOTE**: It is *very* important that you use `make` instead of `make build` when updating MumbleDJ as the first option will grab the latest updates from MumbleDJ's dependencies.

**3)** Check to make sure your configuration in `~/.mumbledj/config/mumbledj.gcfg` is the same as before. If it is back to the default, a backup should have been created at `~/.mumbledj/config/mumbledj_backup.gcfg` so you can copy the values back over.

## AUTHOR
[Matthieu Grieger](http://matthieugrieger.com)

## LICENSE
```
The MIT License (MIT)

Copyright (c) 2014, 2015 Matthieu Grieger

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```

## THANKS
* All those who contribute to [Mumble](https://github.com/mumble-voip/mumble).
* [Tim Cooper](https://github.com/bontibon) for [gumble](https://github.com/layeh/gumble).
* [Ricardo Garcia](https://github.com/rg3) for [youtube-dl](https://github.com/rg3/youtube-dl).
* [ScalingData](https://github.com/scalingdata) for [gcfg](https://github.com/scalingdata/gcfg).
* [Jason Moiron](https://github.com/jmoiron) for [jsonq](https://github.com/jmoiron/jsonq).
* [Nitrous.IO](https://github.com/nitrous-io) for [goop](https://github.com/nitrous-io/goop).
