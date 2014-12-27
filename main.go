/*
 * MumbleDJ
 * By Matthieu Grieger
 * main.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

import (
	"flag"
	"fmt"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/layeh/gumble/gumbleutil"
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
	currentSong    *Song
	audioStream    *gumble_ffmpeg.Stream
	homeDir        string
}

func (dj *mumbledj) OnConnect(e *gumble.ConnectEvent) {
	if dj.client.Channels().Find(dj.defaultChannel) != nil {
		dj.client.Self().Move(dj.client.Channels().Find(dj.defaultChannel))
	} else {
		fmt.Println("Channel doesn't exist, staying in root channel...")
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
		dj.audioStream.Done = dj.OnSongFinished
		dj.audioStream.SetVolume(dj.conf.Volume.DefaultVolume)
	} else {
		panic(err)
	}
}

func (dj *mumbledj) OnDisconnect(e *gumble.DisconnectEvent) {
	dj.keepAlive <- true
}

func (dj *mumbledj) OnTextMessage(e *gumble.TextMessageEvent) {
	if e.Message[0] == '!' {
		parseCommand(e.Sender, e.Sender.Name(), e.Message[1:])
	}
}

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

func (dj *mumbledj) OnSongFinished() {
	if err := dj.currentSong.Delete(); err == nil {
		if dj.queue.Len() != 0 {
			dj.currentSong = dj.queue.NextSong()
			if dj.currentSong != nil {
				if err := dj.currentSong.Download(); err == nil {
					dj.currentSong.Play()
				} else {
					panic(err)
				}
			}
		}
	} else {
		panic(err)
	}
}

var dj = mumbledj{
	keepAlive: make(chan bool),
	queue:     NewSongQueue(),
}

func main() {
	var address, port, username, password, channel string

	flag.StringVar(&address, "server", "localhost", "address for Mumble server")
	flag.StringVar(&port, "port", "64738", "port for Mumble server")
	flag.StringVar(&username, "username", "MumbleDJ", "username of MumbleDJ on server")
	flag.StringVar(&password, "password", "", "password for Mumble server (if needed)")
	flag.StringVar(&channel, "channel", "root", "default channel for MumbleDJ")
	flag.Parse()

	dj.client = gumble.NewClient(&dj.config)
	dj.config = gumble.Config{
		Username: username,
		Password: password,
		Address:  address + ":" + port,
	}
	dj.defaultChannel = channel

	dj.client.Attach(gumbleutil.Listener{
		Connect:     dj.OnConnect,
		Disconnect:  dj.OnDisconnect,
		TextMessage: dj.OnTextMessage,
	})

	// IMPORTANT NOTE: This will be changed later once released. Not really safe at the
	// moment.
	dj.config.TLSConfig.InsecureSkipVerify = true
	if err := dj.client.Connect(); err != nil {
		panic(err)
	}

	<-dj.keepAlive
}
