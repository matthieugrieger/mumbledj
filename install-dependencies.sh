#!/bin/sh
set -e

# removing old ffmpeg
sudo rm -rf /usr/bin/ffmpeg
sudo rm -rf /usr/bin/X11/ffmpeg
sudo rm -rf /usr/share/man/man1/ffmpeg.1.gz

# check to see if ffmpeg is installed
if [ ! -f "$HOME/bin/ffmpeg" ]; then
    echo 'Installing ffmpeg'
    wget http://johnvansickle.com/ffmpeg/releases/ffmpeg-release-64bit-static.tar.xz -O /tmp/ffmpeg.tar.xz
    tar -xvf /tmp/ffmpeg.tar.xz --strip 1 --no-anchored ffmpeg ffprobe
    chmod a+rx ffmpeg ffprobe
    mv ff* ~/bin
else
  echo 'Using cached version of ffmpeg.';
fi

# check to see if youtube-dl is installed
if [ ! -f "$HOME/bin/youtube-dl" ]; then
    echo 'Installing youtube-dl'
    curl https://yt-dl.org/latest/youtube-dl -o ~/bin/youtube-dl
    chmod a+rx ~/bin/youtube-dl
else
  echo 'Using cached version of youtube-dl.';
fi

# check to see if opus is installed
if [ ! -d "$HOME/opus/lib" ]; then
    echo 'Installing opus'
    wget http://downloads.xiph.org/releases/opus/opus-1.0.3.tar.gz
    tar xzvf opus-1.0.3.tar.gz
    cd opus-1.0.3 && ./configure --prefix=$HOME/opus && make && make install
else
  echo 'Using cached version of opus.';
fi