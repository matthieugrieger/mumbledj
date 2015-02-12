/*
 * MumbleDJ
 * By Matthieu Grieger
 * main.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/layeh/gopus"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/layeh/gumble/gumbleutil"
	"os"
	"os/user"
)

// MumbleDJ type declaration
type mumbledj struct {
	config         gumble.Config
	client         *gumble.Client
	keepAlive      chan bool
	defaultChannel string
	conf           DjConfig
	queue          *SongQueue
	audioStream    *gumble_ffmpeg.Stream
	homeDir        string
	playlistSkips  map[string][]string
}

// OnConnect event. First moves MumbleDJ into the default channel specified
// via commandline args, and moves to root channel if the channel does not exist. The current
// user's homedir path is stored, configuration is loaded, and the audio stream is set up.
func (dj *mumbledj) OnConnect(e *gumble.ConnectEvent) {
	if dj.client.Channels.Find(dj.defaultChannel) != nil {
		dj.client.Self.Move(dj.client.Channels.Find(dj.defaultChannel))
	} else {
		fmt.Println("Channel doesn't exist or one was not provided, staying in root channel...")
	}

	if currentUser, err := user.Current(); err == nil {
		dj.homeDir = currentUser.HomeDir
	}

	if err := loadConfiguration(); err == nil {
		fmt.Println("Configuration successfully loaded!")
	} else {
		panic(err)
	}

	if audioStream, err := gumble_ffmpeg.New(dj.client); err == nil {
		dj.audioStream = audioStream
		dj.audioStream.Volume = dj.conf.Volume.DefaultVolume
	} else {
		panic(err)
	}

	dj.client.AudioEncoder.SetApplication(gopus.Audio)

	dj.client.Self.SetComment(dj.conf.General.DefaultComment)
}

// OnDisconnect event. Terminates MumbleDJ thread.
func (dj *mumbledj) OnDisconnect(e *gumble.DisconnectEvent) {
	dj.keepAlive <- true
}

// OnTextMessage event. Checks for command prefix, and calls parseCommand if it exists. Ignores
// the incoming message otherwise.
func (dj *mumbledj) OnTextMessage(e *gumble.TextMessageEvent) {
	plainMessage := gumbleutil.PlainText(&e.TextMessage)
	if len(plainMessage) != 0 {
		if plainMessage[0] == dj.conf.General.CommandPrefix[0] && plainMessage != dj.conf.General.CommandPrefix {
			parseCommand(e.Sender, e.Sender.Name, plainMessage[1:])
		}
	}
}

// OnUserChange event. Checks UserChange type, and adjusts items such as skiplists to reflect
// the current status of the users on the server.
func (dj *mumbledj) OnUserChange(e *gumble.UserChangeEvent) {
	if e.Type.Has(gumble.UserChangeDisconnected) {
		if dj.audioStream.IsPlaying() {
			if dj.queue.CurrentSong().playlist != nil {
				dj.queue.CurrentSong().playlist.RemoveSkip(e.User.Name)
			}
			dj.queue.CurrentSong().RemoveSkip(e.User.Name)
		}
	}
}

// Checks if username has the permissions to execute a command. Permissions are specified in
// mumbledj.gcfg.
func (dj *mumbledj) HasPermission(username string, command bool) bool {
	if dj.conf.Permissions.AdminsEnabled && command {
		for _, adminName := range dj.conf.Permissions.Admins {
			if username == adminName {
				return true
			}
		}
		return false
	} else {
		return true
	}
}

// Sends a private message to a user. Essentially just checks if a user is still in the server
// before sending them the message.
func (dj *mumbledj) SendPrivateMessage(user *gumble.User, message string) {
	if targetUser := dj.client.Self.Channel.Users.Find(user.Name); targetUser != nil {
		targetUser.Send(message)
	}
}

// dj variable declaration. This is done outside of main() to allow global use.
var dj = mumbledj{
	keepAlive:     make(chan bool),
	queue:         NewSongQueue(),
	playlistSkips: make(map[string][]string),
}

// Main function, but only really performs startup tasks. Grabs and parses commandline
// args, sets up the gumble client and its listeners, and then connects to the server.
func main() {
	var address, port, username, password, channel, pemCert, pemKey string

	flag.StringVar(&address, "server", "localhost", "address for Mumble server")
	flag.StringVar(&port, "port", "64738", "port for Mumble server")
	flag.StringVar(&username, "username", "MumbleDJ", "username of MumbleDJ on server")
	flag.StringVar(&password, "password", "", "password for Mumble server (if needed)")
	flag.StringVar(&channel, "channel", "root", "default channel for MumbleDJ")
	flag.StringVar(&pemCert, "cert", "", "path to user PEM certificate for MumbleDJ")
	flag.StringVar(&pemKey, "key", "", "path to user PEM key for MumbleDJ")
	flag.Parse()

	dj.client = gumble.NewClient(&dj.config)
	dj.config = gumble.Config{
		Username: username,
		Password: password,
		Address:  address + ":" + port,
	}
	if pemCert != "" {
		if pemKey == "" {
			pemKey = pemCert
		}
		if certificate, err := tls.LoadX509KeyPair(pemCert, pemKey); err != nil {
			panic(err)
		} else {
			dj.config.TLSConfig.Certificates = append(dj.config.TLSConfig.Certificates, certificate)
		}
	}

	dj.defaultChannel = channel

	dj.client.Attach(gumbleutil.Listener{
		Connect:     dj.OnConnect,
		Disconnect:  dj.OnDisconnect,
		TextMessage: dj.OnTextMessage,
		UserChange:  dj.OnUserChange,
	})
	dj.client.Attach(gumbleutil.AutoBitrate)

	// IMPORTANT NOTE: This will be changed later once released. Not really safe at the
	// moment.
	dj.config.TLSConfig.InsecureSkipVerify = true
	if err := dj.client.Connect(); err != nil {
		fmt.Printf("Could not connect to Mumble server at %s:%s.\n", address, port)
		os.Exit(1)
	}

	<-dj.keepAlive
}
