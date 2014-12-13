MumbleDJ
========
A Mumble bot that plays music fetched from YouTube videos. I have decided to experiment with rewriting the bot in Go using [gumble](https://github.com/layeh/gumble). I am hoping this will cut down on the dependency list, make it easier to develop in the future, and allow for some extra functionality that wasn't previously possible.

And yes, I know that technically this is v3. The Ruby implementation had problems with high CPU usage and choppy audio which I couldn't seem to figure out.

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
* [Tim Cooper](https://github.com/bontibon) for [gumble](https://github.com/layeh/gumble).
* [Ricardo Garcia](https://github.com/rg3) for [youtube-dl](https://github.com/rg3/youtube-dl).
* [Tom Preston-Werner](https://github.com/mojombo) for [TOML](https://github.com/toml-lang/toml).
* [Andrew Gallant](https://github.com/BurntSushi) for [TOML Golang parser](https://github.com/BurntSushi/toml).
