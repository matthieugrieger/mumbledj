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
	"os"
	"os/user"
	"reflect"
	"strings"
	"time"

	"github.com/layeh/gopus"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/layeh/gumble/gumbleutil"
)

// mumbledj is a struct that keeps track of all aspects of the bot's current
// state.
type mumbledj struct {
	config         gumble.Config
	client         *gumble.Client
	keepAlive      chan bool
	defaultChannel []string
	conf           DjConfig
	queue          *SongQueue
	audioStream    *gumble_ffmpeg.Stream
	homeDir        string
	playlistSkips  map[string][]string
	cache          *SongCache
}

// OnConnect event. First moves MumbleDJ into the default channel specified
// via commandline args, and moves to root channel if the channel does not exist. The current
// user's homedir path is stored, configuration is loaded, and the audio stream is set up.
func (dj *mumbledj) OnConnect(e *gumble.ConnectEvent) {
	if dj.client.Channels.Find(dj.defaultChannel...) != nil {
		dj.client.Self.Move(dj.client.Channels.Find(dj.defaultChannel...))
	} else {
		fmt.Println("Channel doesn't exist or one was not provided, staying in root channel...")
	}

	dj.audioStream = gumble_ffmpeg.New(dj.client)
	dj.audioStream.Volume = dj.conf.Volume.DefaultVolume

	if dj.conf.General.PlayerCommand == "ffmpeg" || dj.conf.General.PlayerCommand == "avconv" {
		dj.audioStream.Command = dj.conf.General.PlayerCommand
	} else {
		fmt.Println("Invalid PlayerCommand configuration value. Only \"ffmpeg\" and \"avconv\" are supported. Defaulting to ffmpeg...")
	}

	dj.client.AudioEncoder.SetApplication(gopus.Audio)

	dj.client.Self.SetComment(dj.conf.General.DefaultComment)

	if dj.conf.Cache.Enabled {
		dj.cache.Update()
		go dj.cache.ClearExpired()
	}
}

// OnDisconnect event. Terminates MumbleDJ thread.
func (dj *mumbledj) OnDisconnect(e *gumble.DisconnectEvent) {
	if e.Type == gumble.DisconnectError || e.Type == gumble.DisconnectKicked {
		fmt.Println("Disconnected from server... Will retry connection in 30 second intervals for 15 minutes.")
		reconnectSuccess := false
		for retries := 0; retries <= 30; retries++ {
			fmt.Println("Retrying connection...")
			if err := dj.client.Connect(); err == nil {
				fmt.Println("Successfully reconnected to the server!")
				reconnectSuccess = true
				break
			}
			time.Sleep(30 * time.Second)
		}
		if !reconnectSuccess {
			fmt.Println("Could not reconnect to server. Exiting...")
			dj.keepAlive <- true
			os.Exit(1)
		}
	} else {
		dj.keepAlive <- true
	}
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
			if !isNil(dj.queue.CurrentSong().Playlist()) {
				dj.queue.CurrentSong().Playlist().RemoveSkip(e.User.Name)
			}
			dj.queue.CurrentSong().RemoveSkip(e.User.Name)
		}
	}
}

// HasPermission checks if username has the permissions to execute a command. Permissions are specified in
// mumbledj.gcfg.
func (dj *mumbledj) HasPermission(username string, command bool) bool {
	if dj.conf.Permissions.AdminsEnabled && command {
		for _, adminName := range dj.conf.Permissions.Admins {
			if username == adminName {
				return true
			}
		}
		return false
	}
	return true
}

// SendPrivateMessage sends a private message to a user. Essentially just checks if a user is still in the server
// before sending them the message.
func (dj *mumbledj) SendPrivateMessage(user *gumble.User, message string) {
	if targetUser := dj.client.Self.Channel.Users.Find(user.Name); targetUser != nil {
		targetUser.Send(message)
	}
}

// CheckAPIKeys enables the services with API keys in the environment varaibles
func CheckAPIKeys() {
	anyDisabled := false

	// Checks YouTube API key
	if dj.conf.ServiceKeys.Youtube == "" {
		anyDisabled = true
		fmt.Printf("The youtube service has been disabled as you do not have a YouTube API key defined in your config file!\n")
	} else {
		services = append(services, YouTube{})
	}

	// Checks Soundcloud API key
	if dj.conf.ServiceKeys.SoundCloud == "" {
		anyDisabled = true
		fmt.Printf("The soundcloud service has been disabled as you do not have a Soundcloud API key defined in your config file!\n")
	} else {
		services = append(services, SoundCloud{})
	}

	// Checks to see if any service was disabled
	if anyDisabled {
		fmt.Printf("Please see the following link for info on how to enable missing services: https://github.com/matthieugrieger/mumbledj\n")
	}

	// Exits application if no services are enabled
	if services == nil {
		fmt.Printf("No services are enabled, and thus closing\n")
		os.Exit(1)
	}
}

// isNil checks to see if an object is nil
func isNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}

// dj variable declaration. This is done outside of main() to allow global use.
var dj = mumbledj{
	keepAlive:     make(chan bool),
	queue:         NewSongQueue(),
	playlistSkips: make(map[string][]string),
	cache:         NewSongCache(),
}

// main primarily performs startup tasks. Grabs and parses commandline
// args, sets up the gumble client and its listeners, and then connects to the server.
func main() {

	if currentUser, err := user.Current(); err == nil {
		dj.homeDir = currentUser.HomeDir
	}

	if err := loadConfiguration(); err == nil {
		fmt.Println("Configuration successfully loaded!")
	} else {
		panic(err)
	}

	var address, port, username, password, channel, pemCert, pemKey, accesstokens string
	var insecure bool

	flag.StringVar(&address, "server", "localhost", "address for Mumble server")
	flag.StringVar(&port, "port", "64738", "port for Mumble server")
	flag.StringVar(&username, "username", "MumbleDJ", "username of MumbleDJ on server")
	flag.StringVar(&password, "password", "", "password for Mumble server (if needed)")
	flag.StringVar(&channel, "channel", "root", "default channel for MumbleDJ")
	flag.StringVar(&pemCert, "cert", "", "path to user PEM certificate for MumbleDJ")
	flag.StringVar(&pemKey, "key", "", "path to user PEM key for MumbleDJ")
	flag.StringVar(&accesstokens, "accesstokens", "", "list of access tokens for channel auth")
	flag.BoolVar(&insecure, "insecure", false, "skip certificate checking")
	flag.Parse()

	dj.config = gumble.Config{
		Username: username,
		Password: password,
		Address:  address + ":" + port,
		Tokens:   strings.Split(accesstokens, " "),
	}
	dj.client = gumble.NewClient(&dj.config)

	dj.config.TLSConfig.InsecureSkipVerify = true
	if !insecure {
		gumbleutil.CertificateLockFile(dj.client, fmt.Sprintf("%s/.mumbledj/cert.lock", dj.homeDir))
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

	dj.defaultChannel = strings.Split(channel, "/")

	CheckAPIKeys()

	dj.client.Attach(gumbleutil.Listener{
		Connect:     dj.OnConnect,
		Disconnect:  dj.OnDisconnect,
		TextMessage: dj.OnTextMessage,
		UserChange:  dj.OnUserChange,
	})
	dj.client.Attach(gumbleutil.AutoBitrate)

	if err := dj.client.Connect(); err != nil {
		fmt.Printf("Could not connect to Mumble server at %s:%s.\n", address, port)
		os.Exit(1)
	}

	<-dj.keepAlive
}
