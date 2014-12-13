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
	//"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/layeh/gumble/gumbleutil"
)

// MumbleDJ type declaration
type mumbledj struct {
	config gumble.Config
	client *gumble.Client
	keepAlive chan bool
	defaultChannel string
	conf djConfig
}

func (dj *mumbledj) OnConnect(e *gumble.ConnectEvent) {
	dj.client.Self().Move(dj.client.Channels().Find(dj.defaultChannel))
	
	var err error
	dj.conf, err = loadConfiguration()
	if err == nil {
		fmt.Println("Configuration successfully loaded!")
	} else {
		panic(err)
	}
}

func (dj *mumbledj) OnDisconnect(e *gumble.DisconnectEvent) {
	dj.keepAlive <- true
}

func (dj *mumbledj) OnTextMessage(e *gumble.TextMessageEvent) {
	if e.Message[0] == '!' {
		parseCommand(e.Sender.Name(), e.Message[1:])
	}
}

var dj = mumbledj {
	keepAlive: make(chan bool),
}

func main() {
	var address, port, username, password, channel string
	flag.StringVar(&address, "server", "localhost", "address for Mumble server")
	flag.StringVar(&port, "port", "64738", "port for Mumble server")
	flag.StringVar(&username, "username", "MumbleDJ", "username of MumbleDJ on server")
	flag.StringVar(&password, "password", "", "password for Mumble server (if needed)")
	flag.StringVar(&channel, "channel", "", "default channel for MumbleDJ")
	flag.Parse()
	
	dj.client = gumble.NewClient(&dj.config)
	dj.config = gumble.Config{
		Username: username,
		Password: password,
		Address: address + ":" + port,
	}
	dj.defaultChannel = channel
	
	dj.client.Attach(gumbleutil.Listener{
		Connect: dj.OnConnect,
		Disconnect: dj.OnDisconnect,
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

