MumbleDJ
========
A Mumble bot that plays music fetched from YouTube videos.

**IMPORTANT NOTE:** If you were using the Lua version of MumbleDJ previously, you will need to follow the installation guide once more to install new dependencies.

## USAGE
`$ mumbledj -server=localhost -port=64738 -username=MumbleDJ -password="" -channel=root`  
All parameters are optional, the example above shows the default values for each field.

## COMMANDS
These are all of the chat commands currently supported by MumbleDJ. All command names and command prefixes may be changed in `mumbledj.gcfg`. All fields surrounded by `<>` indicate fields that *must* be supplied to the bot for the command to execute. All fields surrounded by `<>?` are optional fields.

####`!add <youtube_video_url OR youtube_playlist_url>`
Adds a YouTube video's audio to the song queue. If no songs are currently in the queue, the audio will begin playing immediately. YouTube playlists may also be added using this command. Please note, 
however, that if a YouTube playlist contains over 25 videos only the first 25 videos will be placed in the song queue.

####`!skip`
Submits a vote to skip the current song. Once the skip ratio target (specified in `mumbledj.gcfg`) is met, the song will be skipped and the next will start playing. Each user may only submit one skip per song.

####`!skipplaylist`
Submits a vote to skip the current playlist. Once the skip ratio target (specified in `mumbledj.gcfg`) is met, the playlist will be skipped and the next song/playlist will start playing. Each user may only submit one skip per playlist.

####`!forceskip`
An admin command that forces a song skip.

####`!forceskipplaylist`
An admin command that forces a playlist skip.

####`!volume <desired_volume>?`
Either outputs the current volume or changes the current volume. If `desired_volume` is not provided, the current volume will be displayed in chat. Otherwise, the volume for the bot will be changed to `desired_volume` if it is within the allowed volume range.

####`!move <channel>`
Moves MumbleDJ into `channel` if it exists.

####`!reload`
Reloads `mumbledj.gcfg` to retrieve updated configuration settings.

####`!reset`
Resets the song queue.

####`!kill`
Safely cleans the bot environment and disconnects from the server. Please use this command to stop the bot instead of force closing, as the kill command deletes any remaining songs in the `~/.mumbledj/songs` directory.

## INSTALLATION
Installation for v2 of MumbleDJ is much easier than it was before, due to the reduced dependency list and a `Makefile` which automates some of the process.

**NOTE:** This bot was designed for use on Linux machines. If you wish to run the bot on another OS, it will require tweaking and is not something I will be able to help with.  
**NOTE #2:** Your Mumble server MUST be using the Opus audio codec, not CELT. Audio will not play if your server uses CELT.

**SETUP GUIDE**  
**1)** Install and correctly configure [`Go`](https://golang.org/) (1.3 or higher). Specifically, make sure to follow [this guide](https://golang.org/doc/code.html) and set the `GOPATH` environment variable properly.

**2)** Install [`ffmpeg`](https://www.ffmpeg.org/) and [`mercurial`](http://mercurial.selenic.com/) if they are not already installed on your system.

**3)** Install [`youtube-dl`](https://github.com/rg3/youtube-dl#installation).

**4)** Clone the `MumbleDJ` repository or [download the latest release](https://github.com/matthieugrieger/mumbledj/releases).

**5)** `cd` into the `MumbleDJ` repository directory and execute the following commands: 
```
$ make install_deps
$ make
$ make install
```

**5)** Edit `~/.mumbledj/config/mumbledj.gcfg` to your liking. This file will be overwritten if the config file structure is changed in a commit, but a backup is always stored at `~/.mumbledj/config/mumbledj_backup.gcfg`.

**6)** Execute the command shown at the top of this `README` document with your credentials, and the bot should be up and running!

## AUTHOR
[Matthieu Grieger](http://matthieugrieger.com)

## LICENSE
```
The MIT License (MIT)

Copyright (c) 2014 Matthieu Grieger

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
* [kennygrant](https://github.com/kennygrant) for [sanitize](https://github.com/kennygrant/sanitize).
* [Jason Moiron](https://github.com/jmoiron) for [jsonq](https://github.com/jmoiron/jsonq).
