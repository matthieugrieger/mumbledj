MumbleDJ
========
A Mumble bot that plays music fetched from YouTube videos. There are ways to play music with a bot on Mumble already, but I wasn't really satisfied with them. Many of them require the Windows client to be opened along with other applications which is not ideal. My goal with this project is to make a Linux-friendly Mumble bot that can run on a webserver instead of a personal computer.

## Planned Features
These are features that I would like to complete before I consider the project "finished."
* Play YouTube audio through Mumble to be heard by others. (**Done**)
* Commands for adding/queueing YouTube audio by linking to the URL. (**Done**)
* Current song updates in the chat. (**Done**)
* Command to raise/lower volume. (**Done**)
* Commands to allow users to move the bot from one channel to another. (**Done**)
* Admin commands to completely shut down the bot and perform other admin actions. (**Done**)
* Automatically turn on/off music when a certain number of users are inside the channel.

## Hopeful Features
These are features that would be cool to have, but might not ever make it into the bot.
* Ability to control music using voice commands instead of text commands.
* Ability to enable text-to-speech for the bot to announce the name of the next track over voice chat.
* A search command where the bot searches for a song and then plays it. Basically the add/queue command but without the URL.
* Statistics tracking.

## Setup
Since the setup process is a bit extensive, the setup guide can be found in [SETUP.md](https://github.com/matthieugrieger/mumbledj/blob/master/SETUP.md).

## Dependencies
* [OpenSSL](http://www.openssl.org/)
* [Lua 5.2](http://www.lua.org/)
* [libev](http://libev.schmorp.de/)
* [protobuf-c](https://github.com/protobuf-c/protobuf-c)
* [Ogg Vorbis](https://xiph.org/vorbis/)
* [Opus](http://www.opus-codec.org/)
* [Python 2.6 or above](https://www.python.org/)
* [pafy](https://github.com/np1/pafy/)
* [piepan](https://github.com/layeh/piepan)
* [Jansson](http://www.digip.org/jansson/)
* [jshon](http://kmkeen.com/jshon/)
* [ffmpeg](https://www.ffmpeg.org/)

## Author
[Matthieu Grieger](http://matthieugrieger.com)

## License
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

## Thanks
* All those who contribute to [Mumble](https://github.com/mumble-voip/mumble).  
* [Tim Cooper](https://github.com/bontibon) for [piepan](https://github.com/layeh/piepan).
* [Pierre Chapuis](https://github.com/catwell) for [deque](https://github.com/catwell/cw-lua/tree/master/deque).
