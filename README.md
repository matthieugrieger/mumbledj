MumbleDJ
========
A Mumble bot that plays music fetched from YouTube videos. There are ways to play music with a bot on Mumble already, but I wasn't really satisfied with them. Many of them require the Windows client to be opened along with other applications which is not ideal. My goal with this project is to make a Linux-friendly Mumble bot that can run on a webserver instead of a personal computer.

**NOTE:** I am still dealing with some instability issues with regards to the audio playback. If MumbleDJ does not start playing your song within around 30 seconds, just restart the bot. I'm working on a fix.

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
