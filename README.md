MumbleDJ
========
A Mumble bot that plays music fetched from YouTube videos. I have decided to experiment with rewriting the bot in Go using [gumble](https://github.com/layeh/gumble). I am hoping this will cut down on the dependency list, make it easier to develop in the future, and allow for some extra functionality that wasn't previously possible.

And yes, I know that technically this is v3. The Ruby implementation had problems with high CPU usage and choppy audio which I couldn't seem to figure out.

## USAGE
#####`$ mumbledj -server=localhost -port=64738 -username=MumbleDJ -password="" -channel=root`  
All parameters are optional, the example above shows the default values for each field.

## INSTALLATION
Installation for v2 of MumbleDJ is much easier than it was before, due to the reduced dependency list and a `Makefile` which automates some of the process.  

**NOTE:** This bot was designed for use on Linux machines. If you wish to run the bot on another OS, it will require tweaking and is not something I will be able to help with.

**SETUP GUIDE**  
**1)** Install and correctly configure [`Go`](https://golang.org/). Specifically, make sure to follow [this guide](https://golang.org/doc/code.html) and set the `GOPATH` environment variable properly.

**2)** Install [`ffmpeg`](https://www.ffmpeg.org/) if it is not already installed on your system.

**3)** Clone the `MumbleDJ` repository.

**4)** `cd` into the `MumbleDJ` repository directory and execute the following commands: 
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
* [Hicham Bouabdallah](https://github.com/hishboy) for [Golang queue implementation](https://github.com/hishboy/gocommons/blob/master/lang/queue.go).
* [kennygrant](https://github.com/kennygrant) for [sanitize](https://github.com/kennygrant/sanitize).
* [Jason Moiron](https://github.com/jmoiron) for [jsonq](https://github.com/jmoiron/jsonq).
