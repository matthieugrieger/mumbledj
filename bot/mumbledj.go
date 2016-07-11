/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/mumbledj.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleffmpeg"
	"github.com/layeh/gumble/gumbleutil"
	"github.com/matthieugrieger/mumbledj/interfaces"
	"github.com/spf13/viper"
)

// MumbleDJ is a struct that keeps track of all aspects of the bot's state.
type MumbleDJ struct {
	AvailableServices []interfaces.Service
	Client            *gumble.Client
	GumbleConfig      *gumble.Config
	TLSConfig         *tls.Config
	AudioStream       *gumbleffmpeg.Stream
	Queue             interfaces.Queue
	Cache             *Cache
	Skips             interfaces.SkipTracker
	Commands          []interfaces.Command
	Version           string
	Volume            float32
	YouTubeDL         *YouTubeDL
	KeepAlive         chan bool
}

// DJ is a struct that keeps track of all aspects of MumbleDJ's environment.
var DJ *MumbleDJ

// NewMumbleDJ initializes and returns a MumbleDJ type.
func NewMumbleDJ() *MumbleDJ {
	SetDefaultConfig()

	return &MumbleDJ{
		AvailableServices: make([]interfaces.Service, 0),
		TLSConfig:         new(tls.Config),
		Queue:             NewQueue(),
		Cache:             NewCache(),
		Skips:             NewSkipTracker(),
		Commands:          make([]interfaces.Command, 0),
		YouTubeDL:         new(YouTubeDL),
		KeepAlive:         make(chan bool),
	}
}

// OnConnect event. First moves MumbleDJ into the default channel if one exists.
// The configuration is loaded and the audio stream is initialized.
func (dj *MumbleDJ) OnConnect(e *gumble.ConnectEvent) {
	dj.AudioStream = nil
	logrus.WithFields(logrus.Fields{
		"volume": fmt.Sprintf("%.2f", viper.GetFloat64("volume.default")),
	}).Infoln("Setting default volume...")
	dj.Volume = float32(viper.GetFloat64("volume.default"))

	if viper.GetBool("cache.enabled") {
		logrus.Infoln("Caching enabled.")
		dj.Cache.UpdateStatistics()
		go dj.Cache.CleanPeriodically()
	} else {
		logrus.Infoln("Caching disabled.")
	}
}

// OnDisconnect event. Terminates MumbleDJ process or retries connection if
// automatic connection retries are enabled.
func (dj *MumbleDJ) OnDisconnect(e *gumble.DisconnectEvent) {
	dj.Queue.Reset()
	if viper.GetBool("connection.retry_enabled") &&
		(e.Type == gumble.DisconnectError || e.Type == gumble.DisconnectKicked) {
		logrus.WithFields(logrus.Fields{
			"interval_secs": fmt.Sprintf("%d", viper.GetInt("connection.retry_interval")),
			"attempts":      fmt.Sprintf("%d", viper.GetInt("connection.retry_attempts")),
		}).Warnln("Disconnected from server. Retrying connection...")

		success := false
		for retries := 0; retries < viper.GetInt("connection.retry_attempts"); retries++ {
			logrus.Infoln("Retrying connection...")
			if client, err := gumble.DialWithDialer(new(net.Dialer), viper.GetString("connection.address")+":"+viper.GetString("connection.port"), dj.GumbleConfig, dj.TLSConfig); err == nil {
				dj.Client = client
				logrus.Infoln("Successfully reconnected to the server!")
				success = true
				break
			}
			time.Sleep(time.Duration(viper.GetInt("connection.retry_interval")) * time.Second)
		}
		if !success {
			dj.KeepAlive <- true
			logrus.Fatalln("Could not reconnect to server. Exiting...")
		}
	} else {
		dj.KeepAlive <- true
		logrus.Fatalln("Disconnected from server. No reconnect attempts will be made.")
	}
}

// OnTextMessage event. Checks for command prefix and passes it to the Commander
// if it exists. Ignores the incoming message otherwise.
func (dj *MumbleDJ) OnTextMessage(e *gumble.TextMessageEvent) {
	plainMessage := gumbleutil.PlainText(&e.TextMessage)
	if len(plainMessage) != 0 {
		if plainMessage[0] == viper.GetString("commands.prefix")[0] &&
			plainMessage != viper.GetString("commands.prefix") {
			go func() {
				message, isPrivateMessage, err := dj.FindAndExecuteCommand(e.Sender, plainMessage[1:])
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"user":    e.Sender.Name,
						"message": err.Error(),
					}).Warnln("Sending an error message...")
					dj.SendPrivateMessage(e.Sender, fmt.Sprintf("<b>Error:</b> %s", err.Error()))
				} else {
					if isPrivateMessage {
						logrus.WithFields(logrus.Fields{
							"user":    e.Sender.Name,
							"message": message,
						}).Infoln("Sending a private message...")
						dj.SendPrivateMessage(e.Sender, message)
					} else {
						logrus.WithFields(logrus.Fields{
							"channel": dj.Client.Self.Channel.Name,
							"message": message,
						}).Infoln("Sending a message to channel...")
						dj.Client.Self.Channel.Send(message, false)
					}
				}
			}()
		}
	}
}

// OnUserChange event. Checks UserChange type and adjusts skip trackers to
// reflect the current status of the users on the server.
func (dj *MumbleDJ) OnUserChange(e *gumble.UserChangeEvent) {
	if e.Type.Has(gumble.UserChangeDisconnected) || e.Type.Has(gumble.UserChangeChannel) {
		logrus.WithFields(logrus.Fields{
			"user": e.User.Name,
		}).Infoln("A user has disconnected or changed channels, updating skip trackers...")
		dj.Skips.RemoveTrackSkip(e.User)
		dj.Skips.RemovePlaylistSkip(e.User)
	}
}

// SendPrivateMessage sends a private message to the specified user. This method
// verifies that the targeted user is still present in the server before attempting
// to send the message.
func (dj *MumbleDJ) SendPrivateMessage(user *gumble.User, message string) {
	dj.Client.Do(func() {
		if targetUser := dj.Client.Self.Channel.Users.Find(user.Name); targetUser != nil {
			targetUser.Send(message)
		}
	})
}

// IsAdmin checks whether a particular Mumble user is a MumbleDJ admin.
// Returns true if the user is an admin, and false otherwise.
func (dj *MumbleDJ) IsAdmin(user *gumble.User) bool {
	for _, admin := range viper.GetStringSlice("admins.names") {
		if user.Name == admin {
			return true
		}
	}
	return false
}

// Connect starts the process for connecting to a Mumble server.
func (dj *MumbleDJ) Connect() error {
	// Perform startup checks before connecting.
	logrus.Infoln("Performing startup checks...")
	PerformStartupChecks()

	// Create Gumble config.
	dj.GumbleConfig = gumble.NewConfig()
	dj.GumbleConfig.Username = viper.GetString("connection.username")
	dj.GumbleConfig.Password = viper.GetString("connection.password")
	dj.GumbleConfig.Tokens = strings.Split(viper.GetString("connection.access_tokens"), ",")

	// Initialize key pair if needed.
	if viper.GetBool("connection.insecure") {
		dj.TLSConfig.InsecureSkipVerify = true
	} else {
		dj.TLSConfig.ServerName = viper.GetString("connection.address")

		if viper.GetString("connection.cert") != "" {
			if viper.GetString("connection.key") == "" {
				viper.Set("connection.key", viper.GetString("connection.cert"))
			}

			if certificate, err := tls.LoadX509KeyPair(viper.GetString("connection.cert"), viper.GetString("connection.key")); err == nil {
				dj.TLSConfig.Certificates = append(dj.TLSConfig.Certificates, certificate)
			} else {
				return err
			}
		}
	}

	// Add user p12 cert if needed.
	if viper.GetString("connection.user_p12") != "" {
		if _, err := os.Stat(viper.GetString("connection.user_p12")); os.IsNotExist(err) {
			return err
		}

		// Create temporary directory for converted p12 file.
		dir, err := ioutil.TempDir("", "mumbledj")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)

		// Create temporary mumbledj.crt.pem from p12 file.
		command := exec.Command("openssl", "pkcs12", "-password", "pass:", "-in", viper.GetString("connection.user_p12"), "-out", dir+"/mumbledj.crt.pem", "-clcerts", "-nokeys")
		if err := command.Run(); err != nil {
			return err
		}

		// Create temporary mumbledj.key.pem from p12 file.
		command = exec.Command("openssl", "pkcs12", "-password", "pass:", "-in", viper.GetString("connection.user_p12"), "-out", dir+"/mumbledj.key.pem", "-nocerts", "-nodes")
		if err := command.Run(); err != nil {
			return err
		}

		if certificate, err := tls.LoadX509KeyPair(dir+"/mumbledj.crt.pem", dir+"/mumbledj.key.pem"); err == nil {
			dj.TLSConfig.Certificates = append(dj.TLSConfig.Certificates, certificate)
		} else {
			return err
		}
	}

	dj.GumbleConfig.Attach(gumbleutil.Listener{
		Connect:     dj.OnConnect,
		Disconnect:  dj.OnDisconnect,
		TextMessage: dj.OnTextMessage,
		UserChange:  dj.OnUserChange,
	})
	dj.GumbleConfig.Attach(gumbleutil.AutoBitrate)

	var connErr error

	logrus.WithFields(logrus.Fields{
		"address": viper.GetString("connection.address"),
		"port":    viper.GetString("connection.port"),
	}).Infoln("Attempting connection to server...")
	if dj.Client, connErr = gumble.DialWithDialer(new(net.Dialer), viper.GetString("connection.address")+":"+viper.GetString("connection.port"), dj.GumbleConfig, dj.TLSConfig); connErr != nil {
		return connErr
	}

	return nil
}

// FindAndExecuteCommand attempts to find a reference to a command in an
// incoming message. If found, the command is executed and the resulting
// message/error is returned.
func (dj *MumbleDJ) FindAndExecuteCommand(user *gumble.User, message string) (string, bool, error) {
	command, err := dj.findCommand(message)
	if err != nil {
		return "", true, errors.New("No command was found in this message")
	}
	return dj.executeCommand(user, message, command)
}

// GetService loops through the available services and determines if a URL
// matches a particular service. If a match is found, the service object is
// returned.
func (dj *MumbleDJ) GetService(url string) (interfaces.Service, error) {
	for _, service := range dj.AvailableServices {
		if service.CheckURL(url) {
			return service, nil
		}
	}
	return nil, errors.New("The provided URL does not match an enabled service")
}

func (dj *MumbleDJ) findCommand(message string) (interfaces.Command, error) {
	var possibleCommand string
	if strings.Contains(message, " ") {
		possibleCommand = strings.ToLower(message[:strings.Index(message, " ")])
	} else {
		possibleCommand = strings.ToLower(message)
	}
	for _, command := range dj.Commands {
		for _, alias := range command.Aliases() {
			if possibleCommand == alias {
				return command, nil
			}
		}
	}
	return nil, errors.New("No command was found in this message")
}

func (dj *MumbleDJ) executeCommand(user *gumble.User, message string, command interfaces.Command) (string, bool, error) {
	canExecute := false
	if viper.GetBool("admins.enabled") && command.IsAdminCommand() {
		canExecute = dj.IsAdmin(user)
	} else {
		canExecute = true
	}

	if canExecute {
		return command.Execute(user, strings.Split(message, " ")[1:]...)
	}
	return "", true, errors.New("You do not have permission to execute this command")
}
