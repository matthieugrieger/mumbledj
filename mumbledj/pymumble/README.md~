PYMUMBLE python library
=======================



Description
-----------
This library act as a mumble client, connecting to a murmur server, exchanging states and audio.
It as been developped as a hobby project, for a specific purpose of implementing an automatic audio
recorder, but it has grown beyond that point.
There are probably a lot a space for improvements, and a lot of features are missing, but I hope
the basic framework is solid and it should be relatively easy to add functionnalites

I'm actually not a professionnal developper (not anymore anyway, and a Cobol implementation of Mumble would be a bit obsolete...),
so probably I took some weird choices or bad designs...  but I though it could be usefull for someone else.

For a client application example, you can check https://github.com/Robert904/mumblerecbot

Status
------
- Compatible with mumble 1.2.4 and normally 1.2.3 and 1.2.2
- Support OPUS and CELT (both 0.7 and 0.11) codecs.  Speex is not supported
- Receive and send audio, get users and channels status
- Set properties for users (mute, comments, etc.) and go to a specific channel
- Callback mechanism to react on server events
- Manage the blobs (images, long comments, etc.)
- Handle text messages

### What is missing:
- UDP media.  Currently it works only in TCP tunneling mode (the standard fallback of mumble when UDP is not working)
- basically server management (user creation and registration, ACLs, groups, bans, etc.)
- Positionning is not managed, but it should be easy to add
- Audio targets (whisper, etc.) is not managed in outgoing audio, and has very basic support in incoming
- ping statistics
- Probably a lot of other small features
- polishing ?

Architecture
------------
The library is based on the Mumble object, which is basically a thread.  When started, it will try
to connect to the server and start exchange the connections messages with the server.
This thread is in a loop that take care of the pings, send the commands to the server,
check for incoming messages including audio and check for audio to be sent.
The rate of that loop is controlled by how long it will wait for an incoming message before going further.

It rely on several other modules and objects, but they should probably never be instanciated by an application.

The OPUS and CELT support is achieved by wrapping the C library in a cython module.
these implementations are focussed on what mumble need, but pyopus should be quite clean and extended to a full blown
python binding...

Requirements/installation
-------------------------
It seems to work fine on Python 2.6 and 2.7.
I have used it on both Windows and Linux

Cython is needed, at least 0.14, and you need a worinkg compiler environment (I use MINGW for windows)

in the pyopus and pycelt directories, there is a basic Makefile that should compile the library and create the
loadable compiled Python module at the correct place.  You have to edit these Makefiles to select your environment (Linux or Mingw)
If your cython installation use PIC compiled libraries, you will also hav to uncomment the "CONFIGURE_OPTS = "--with-pic" line.  You'll know if you get compilation error like "relocation R_X86_64_32S against `.rodata' can not be used when making a shared object; recompile with -fPIC"...  no idea how to check that automatically...

Issues
------
It seems one thread keep running when the other crash...  I have to look into that

License
-------
Copyright Robert Hendrickx <rober@percu.be> - 2014

pymumble is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.


Included opus and celt libraries sources have their own licensing
