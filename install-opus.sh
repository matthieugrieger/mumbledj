#!/bin/sh
set -e
# check to see if opus folder is empty
if [ ! -d "$HOME/opus/lib" ]; then
    wget http://downloads.xiph.org/releases/opus/opus-1.0.3.tar.gz
    tar xzvf opus-1.0.3.tar.gz
    cd opus-1.0.3 && ./configure --prefix=$HOME/opus && make && make install
else
  echo 'Using cached directory.';
fi