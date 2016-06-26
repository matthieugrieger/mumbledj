
<h1 align="center">MumbleDJ</h1>
<p align="center"><b>A Mumble bot that plays audio fetched from various media websites.</b></p>
<p align="center"><a href="https://travis-ci.org/matthieugrieger/mumbledj"><img src="https://travis-ci.org/matthieugrieger/mumbledj.svg?branch=master"/></a> <a href="https://raw.githubusercontent.com/matthieugrieger/mumbledj/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg"/></a> <a href="https://github.com/matthieugrieger/mumbledj/releases"><img src="https://img.shields.io/github/release/matthieugrieger/mumbledj.svg"/></a> <a href="https://goreportcard.com/report/github.com/matthieugrieger/mumbledj"><img src="https://goreportcard.com/badge/github.com/matthieugrieger/mumbledj"/></a> <a href="https://codecov.io/gh/matthieugrieger/mumbledj"><img src="https://img.shields.io/codecov/c/github/matthieugrieger/mumbledj.svg"/></a> <a href="https://gitter.im/matthieugrieger/mumbledj"><img src="https://img.shields.io/gitter/room/matthieugrieger/mumbledj.svg" /></a></p>

## Table of Contents

* [Features](#features)
* [Installation](#installation)
  * [Requirements](#requirements)
    * [YouTube API Key](#youtube-api-key)
    * [SoundCloud API Key](#soundcloud-api-key)
  * [Via `go get`](#via-go-get-recommended)
  * [Pre-compiled Binaries](#pre-compiled-binaries-easiest)
  * [From Source](#from-source)
* [Usage](#usage)
* [Commands](#commands)
* [Contributing](#contributing)
* [Author](#author)
* [License](#license)
* [Thanks](#thanks)

## Features
* Plays audio from many media websites, including YouTube, SoundCloud, and Mixcloud.
* Supports playlists and individual videos/tracks.
* Displays metadata in the text chat whenever a new track starts playing.
* Incredibly customizable. Nearly everything is able to be tweaked via configuration files (by default located at `$HOME/.config/mumbledj/config.yaml`).
* A large array of [commands](#commands) that perform a wide variety of functions.
* Built-in vote-skipping.
* Built-in caching system (disabled by default).
* Built-in play/pause/volume control.

## Installation
**IMPORTANT NOTE:** MumbleDJ is only tested and developed for Linux systems. Support will not be given for non-Linux systems if problems are encountered.

### Requirements
**All MumbleDJ installations must also have the following installed:**
* [`youtube-dl`](https://rg3.github.io/youtube-dl/download.html)
* [`ffmpeg`](https://ffmpeg.org) OR [`avconv`](https://libav.org)
* [`aria2`](https://aria2.github.io/) if you plan on using services that throttle download speeds (like Mixcloud)

**If installing via `go install` or from source, the following must be installed:**
* [Go 1.5+](https://golang.org)
  * __NOTE__: Extra installation steps are required for a working Go installation. Once Go is installed, type `go help gopath` for more information.
  * If the repositories for your distro contain a version of Go older than 1.5, try using [`gvm`](https://github.com/moovweb/gvm) to install Go 1.5 or newer.

#### YouTube API Key
A YouTube API key must be present in your configuration file in order to use the YouTube service within the bot. Below is a guide for retrieving an API key:

**1)** Navigate to the [Google Developers Console](https://console.developers.google.com) and sign in with your Google account, or create one if you haven't already.

**2)** Click the "Create Project" button and give your project a name. It doesn't matter what you set your project name to. Once you have a name click the "Create" button. You should be redirected to your new project once it's ready.

**3)** Click on "APIs & auth" on the sidebar, and then click APIs. Under the "YouTube APIs" header, click "YouTube Data API". Click on the "Enable API" button.

**4)** Click on the "Credentials" option underneath "APIs & auth" on the sidebar. Underneath "Public API access" click on "Create New Key". Choose the "Server key" option.

**5)** Add the IP address of the machine MumbleDJ will run on in the box that appears (this is optional, but improves security). Click "Create".

**6)** You should now see that an API key has been generated. Copy/paste this API key into the configuration file located at `$HOME/.config/mumbledj/mumbledj.yaml`. 

#### SoundCloud API Key
A SoundCloud client ID must be present in your configuration file in order to use the SoundCloud service within the bot. Below is a guide for retrieving a client ID:

**1)** Login/sign up for a SoundCloud account on https://soundcloud.com.

**2)** Create a new app: https://soundcloud.com/you/apps/new.

**3)** You should now see that a client ID has been generated. Copy/paste this ID (NOT the client secret) into the configuration file located at `$HOME/.config/mumbledj/mumbledj.yaml`.


### Via `go get` (recommended)
After verifying that the [requirements](#requirements) are installed, simply issue the following command:
```
go get -u github.com/matthieugrieger/mumbledj
```

This should place a binary in `$GOPATH/bin` that can be used to start the bot.

**NOTE:** If using Go 1.5, you MUST execute the following for `go get` to work:
```
export GO15VENDOREXPERIMENT=1
```

### Pre-compiled Binaries (easiest)
Pre-compiled binaries are provided for convenience. Overall, I do not recommend using these unless you cannot get `go install` to work properly. Binaries compiled on your own machine are likely more efficient as these binaries are cross-compiled from a 64-bit Linux system.

After verifying that the [requirements](#requirements) are installed, simply visit the [releases page](https://github.com/matthieugrieger/mumbledj/releases) and download the appropriate binary for your platform.

### From Source
First, clone the MumbleDJ repository to your machine:
```
git clone https://github.com/matthieugrieger/mumbledj.git
```

Install the required software as described in the [requirements section](#requirements), and execute the following:
```
make
```

This will place a compiled `mumbledj` binary in the cloned directory if successful. If you would like to make the binary more accessible by adding it to `/usr/local/bin`, simply execute the following:
```
sudo make install
```

## Usage
MumbleDJ is a compiled program that is executed via a terminal.

Here is an example helptext that gives you a feel for the various commandline arguments you can give MumbleDJ:

```
NAME:
   MumbleDJ - A Mumble bot that plays audio from various media sites.

USAGE:
   mumbledj [global options] command [command options] [arguments...]
   
VERSION:
   3.0.0-alpha
   
COMMANDS:
GLOBAL OPTIONS:
   --config value, -c value        location of MumbleDJ configuration file (default: "$HOME/.config/mumbledj/mumbledj.yaml")
   --server value, -s value        address of Mumble server to connect to (default: "127.0.0.1")
   --port value, -o value          port of Mumble server to connect to (default: "64738")
   --username value, -u value      username for the bot (default: "MumbleDJ")
   --password value, -p value      password for the Mumble server
   --channel value, -n value       channel the bot enters after connecting to the Mumble server
   --cert value, -e value          path to PEM certificate
   --key value, -k value           path to PEM key
   --accesstokens value, -a value  list of access tokens separated by spaces
   --insecure, -i                  if present, the bot will not check Mumble certs for consistency
   --debug, -d                     if present, all debug messages will be shown
   --help, -h                      show help
   --version, -v                   print the version

```

__NOTE__: You can also override all settings found within `config.yaml` directly from the commandline. Here's an example:

```
mumbledj --admins.names="SuperUser,Matt" --volume.default="0.5" --volume.lowest="0.2" --queue.automatic_shuffle_on="true"
```

Keep in mind that values that contain commas (such as `"SuperUser,Matt"`) will be interpreted as string slices, or arrays if you are not familiar with Go. If you want your value to be interpreted as a normal string, it is best to avoid commas for now.

## Commands

### add
* __Description__: Adds a track or playlist from a media site to the queue.
* __Default Aliases__: add, a
* __Arguments__: (Required) URL(s) to a track or playlist from a supported media site.
* __Admin-only by default__: No
* __Example__: `!add https://www.youtube.com/watch?v=KQY9zrjPBjo`

### addnext
* __Description__: Adds a track or playlist from a media site as the next item in the queue.
* __Default Aliases__: addnext, an
* __Arguments__: (Required) URL(s) to a track or playlist from a supported media site.
* __Admin-only by default__: Yes
* __Example__: `!addnext https://www.youtube.com/watch?v=KQY9zrjPBjo`

### cachesize
* __Description__: Outputs the file size of the cache in MiB if caching is enabled.
* __Default Aliases__: cachesize, cs
* __Arguments__: None
* __Admin-only by default__: Yes
* __Example__: `!cachesize`

### currenttrack
* __Description__: Outputs information about the current track in the queue if one exists.
* __Default Aliases__: currenttrack, currentsong, current
* __Arguments__: None
* __Admin-only by default__: No
* __Example__: `!currenttrack`

### forceskip
* __Description__: Immediately skips the current track.
* __Default Aliases__: forceskip, fs
* __Arguments__: None
* __Admin-only by default__: Yes
* __Example__: `!forceskip`

### forceskipplaylist
* __Description__: Immediately skips the current playlist.
* __Default Aliases__: forceskipplaylist, fsp
* __Arguments__: None
* __Admin-only by default__: Yes
* __Example__: `!forceskipplaylist`

### help
* __Description__: Outputs a list of available commands and their descriptions.
* __Default Aliases__: help, h
* __Arguments__: None
* __Admin-only by default__: No
* __Example__: `!help`

### joinme
* __Description__: Moves MumbleDJ into your current channel if not playing audio to someone else.
* __Default Aliases__: joinme, join
* __Arguments__: None
* __Admin-only by default__: Yes
* __Example__: `!joinme`

### kill
* __Description__: Stops the bot and cleans its cache directory.
* __Default Aliases__: kill, k
* __Arguments__: None
* __Admin-only by default__: Yes
* __Example__: `!kill`

### listtracks
* __Description__: Outputs a list of the tracks currently in the queue.
* __Default Aliases__: listtracks, listsongs, list, l
* __Arguments__: (Optional) Number of tracks to list
* __Admin-only by default__: No
* __Example__: `!listtracks 10`

### move
* __Description__: Moves the bot into the Mumble channel provided via argument.
* __Default Aliases__: move, m
* __Arguments__: (Required) Mumble channel to move the bot into
* __Admin-only by default__: Yes
* __Example__: `!move Music`

### nexttrack
* __Description__: Outputs information about the next track in the queue if one exists.
* __Default Aliases__: nexttrack, nextsong, next
* __Arguments__: None
* __Admin-only by default__: No
* __Example__: `!nexttrack`

### numcached
* __Description__: Outputs the number of tracks cached on disk if caching is enabled.
* __Default Aliases__: numcached, nc
* __Arguments__: None
* __Admin-only by default__: Yes
* __Example__: `!numcached`

### numtracks
* __Description__: Outputs the number of tracks currently in the queue.
* __Default Aliases__: numtracks, numsongs, nt
* __Arguments__: None
* __Admin-only by default__: No
* __Example__: `!numtracks`

### pause
* __Description__: Pauses audio playback.
* __Default Aliases__: pause
* __Arguments__: None
* __Admin-only by default__: No
* __Example__: `!pause`

### reload
* __Description__: Reloads the configuration file.
* __Default Aliases__: reload, r
* __Arguments__: None
* __Admin-only by default__: Yes
* __Example__: `!reload`

### reset
* __Description__: Resets the queue by removing all queue items.
* __Default Aliases__: reset, re
* __Arguments__: None
* __Admin-only by default__: Yes
* __Example__: `!reset`

### resume
* __Description__: Resumes audio playback.
* __Default Aliases__: resume
* __Arguments__: None
* __Admin-only by default__: No
* __Example__: `!pause`

### setcomment
* __Description__: Sets the comment displayed next to MumbleDJ's username in Mumble. If the argument is left empty, the current comment is removed.
* __Default Aliases__: setcomment, comment, sc
* __Arguments__: (Optional) New comment
* __Admin-only by default__: Yes
* __Example__: `!setcomment Hello! I'm a bot. Beep boop.`

### shuffle
* __Description__: Randomizes the tracks currently in the queue.
* __Default Aliases__: shuffle, shuf, sh
* __Arguments__: None
* __Admin-only by default__: Yes
* __Example__: `!shuffle`

### skip
* __Description__: Places a vote to skip the current track.
* __Default Aliases__: skip, s
* __Arguments__: None
* __Admin-only by default__: No
* __Example__: `!skip`

### skipplaylist
* __Description__: Places a vote to skip the current playlist.
* __Default Aliases__: skipplaylist, sp
* __Arguments__: None
* __Admin-only by default__: No
* __Example__: `!skipplaylist`

### toggleshuffle
* __Description__: Toggles permanent track shuffling on/off.
* __Default Aliases__: toggleshuffle, toggleshuf, togshuf, tsh
* __Arguments__: None
* __Admin-only by default__: Yes
* __Example__: `!toggleshuffle`

### version
* __Description__: Outputs the current version of MumbleDJ.
* __Default Aliases__: version, v
* __Arguments__: None
* __Admin-only by default__: No
* __Example__: `!version`

### volume
* __Description__: Changes the volume if an argument is provided, outputs the current volume otherwise.
* __Default Aliases__: volume, vol
* __Arguments__: (Optional) New volume
* __Admin-only by default__: No
* __Example__: `!volume 0.5`

## Contributing

Contributions to MumbleDJ are always welcome! Please see the [contribution guidelines](https://github.com/matthieugrieger/mumbledj/blob/master/CONTRIBUTING.md) for instructions and suggestions!

## Author
[Matthieu Grieger](https://github.com/matthieugrieger)

## License
```
The MIT License (MIT)

Copyright (c) 2016 Matthieu Grieger

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

## Thanks
* [All those who contribute to Mumble](https://github.com/mumble-voip/mumble/graphs/contributors)
* [Tim Cooper](https://github.com/bontibon) for [gumble, gumbleffmpeg, and gumbleutil](https://github.com/layeh/gumble)
* [Jeremy Saenz](https://github.com/codegangsta) for [cli](https://github.com/urfave/cli)
* [Anton Holmquist](https://github.com/antonholmquist) for [jason](https://github.com/antonholmquist/jason)
* [Stretchr, Inc.](https://github.com/stretchr) for [testify](https://github.com/stretchr/testify)
* [ChannelMeter](https://github.com/ChannelMeter) for [iso8601duration](https://github.com/ChannelMeter/iso8601duration)
* [Steve Francia](https://github.com/spf13) for [viper](https://github.com/spf13/viper)
* [Simon Eskildsen](https://github.com/Sirupsen) for [logrus](https://github.com/Sirupsen/logrus)
* [Mitchell Hashimoto](https://github.com/mitchellh) for [gox](https://github.com/mitchellh/gox)
* [Jim Teeuwen](https://github.com/jteeuwen) for [go-bindata](https://github.com/jteeuwen/go-bindata)
