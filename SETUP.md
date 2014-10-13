MumbleDJ Setup Guide
====================
This setup guide is written for installation on an Ubuntu-based system. If I ever get the time I may also write a setup guide for Fedora and Arch Linux as I have gotten them to work on those distros as well.

**NOTE:** This installation guide is written from memory so if there is something wrong, please let me know!

## Installing Dependencies
This is the bulk of the setup process. Most of the dependencies can be installed with the following `apt-get` commands:

```
$ sudo add-apt-repository ppa:mc3man/trusty-media
$ sudo apt-get update
$ sudo apt-get install protobuf-c-compiler libprotobuf-c0-dev lua5.2 liblua5.2-dev libvorbis-dev libssl-dev libev-dev python-pip ffmpeg
```

Then just install the python module `pafy` with `pip`:

```
$ sudo pip install pafy
```

The rest of the dependencies will have to be compiled from source. First we will compile and install `protobuf` and `protobuf-c`.

### protobuf & protobuf-c

```
$ wget https://protobuf.googlecode.com/svn/rc/protobuf-2.6.0.tar.gz
$ tar xzf protobuf-2.6.0.tar.gz
$ cd protobuf-2.6.0
$ ./configure && make && sudo make install && cd ..
$ wget https://github.com/protobuf-c/protobuf-c/releases/download/v1.0.2/protobuf-c-1.0.2.tar.gz
$ tar xzf protobuf-c-1.0.2.tar.gz
$ cd protobuf-c-1.0.2
$ ./configure && make && sudo make install && cd ..
```

You have now (hopefully) successfully installed `protobuf` and `protobuf-c`. Now we will install `jshon`, which is required for `piepan`.

### jshon and jansson

`jshon` depends on `jansson`, so we will compile and install that beforehand.

```
$ wget http://www.digip.org/jansson/releases/jansson-2.7.tar.gz
$ tar xzf jansson-2.7.tar.gz
$ rm jansson-2.7.tar.gz
$ cd jansson-2.7.tar.gz
$ ./configure && make && sudo make install && cd ..
$ wget http://kmkeen.com/jshon/jshon.tar.gz
$ tar xzf jshon.tar.gz
$ cd jshon-*
$ make && sudo cp jshon /usr/local/bin && cd ..
```

Cool, we now have the necessary dependencies installed to compile `piepan`. We will now install `piepan`.

### piepan

```
$ sudo ldconfig
$ git clone https://github.com/layeh/piepan.git
$ cd piepan
$ make
```

Then `cp` the `piepan` executable into your bot's directory or `/usr/local/bin`.

**NOTE:** If your system cannot find `lua` and complains of `pkg-config` stuff, try the following "fix":

```
$ cd piepan
$ make clean
$ nano Makefile
```

Replace the first like of the `Makefile` with the following:

```
CFLAGS = `pkg-config --libs --cflags libssl lua5.2 libprotobuf-c opus vorbis vorbisfile` -lev -pthread
```

If you recompile piepan it *should* work now.


## Configuring & Running MumbleDJ
Now we're ready to finally use MumbleDJ! First clone the MumbleDJ project.

```
$ git clone https://github.com/matthieugrieger/mumbledj
```

Within the `mumbledj` directory you will see another directory inside named `mumbledj`. Open that directory and edit `config.lua` to your liking.

When you're ready to use the bot, just use the following command (this assumes that the `piepan` executable is in the `mumbledj/mumbledj` directory):

```
./piepan -u NAME_OF_BOT -s SERVER_IP mumbledj.lua
```


####All done!


