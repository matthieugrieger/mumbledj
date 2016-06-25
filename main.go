/*
 * MumbleDJ
 * By Matthieu Grieger
 * main.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/matthieugrieger/mumbledj/bot"
	"github.com/matthieugrieger/mumbledj/commands"
	"github.com/matthieugrieger/mumbledj/services"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// DJ is a global variable that holds various details about the bot's state.
var DJ = bot.NewMumbleDJ()

func init() {
	DJ.Commands = commands.Commands
	DJ.AvailableServices = services.Services

	// Injection into sub-packages.
	commands.DJ = DJ
	services.DJ = DJ
	bot.DJ = DJ

	DJ.Version = "v3.0.4"

	logrus.SetLevel(logrus.WarnLevel)
}

func main() {
	app := cli.NewApp()
	app.Name = "MumbleDJ"
	app.Usage = "A Mumble bot that plays audio from various media sites."
	app.Version = DJ.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: os.ExpandEnv("$HOME/.config/mumbledj/config.yaml"),
			Usage: "location of MumbleDJ configuration file",
		},
		cli.StringFlag{
			Name:  "server, s",
			Value: "127.0.0.1",
			Usage: "address of Mumble server to connect to",
		},
		cli.StringFlag{
			Name:  "port, o",
			Value: "64738",
			Usage: "port of Mumble server to connect to",
		},
		cli.StringFlag{
			Name:  "username, u",
			Value: "MumbleDJ",
			Usage: "username for the bot",
		},
		cli.StringFlag{
			Name:  "password, p",
			Value: "",
			Usage: "password for the Mumble server",
		},
		cli.StringFlag{
			Name:  "channel, n",
			Value: "",
			Usage: "channel the bot enters after connecting to the Mumble server",
		},
		cli.StringFlag{
			Name:  "cert, e",
			Value: "",
			Usage: "path to PEM certificate",
		},
		cli.StringFlag{
			Name:  "key, k",
			Value: "",
			Usage: "path to PEM key",
		},
		cli.StringFlag{
			Name:  "accesstokens, a",
			Value: "",
			Usage: "list of access tokens separated by spaces",
		},
		cli.BoolFlag{
			Name:  "insecure, i",
			Usage: "if present, the bot will not check Mumble certs for consistency",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "if present, all debug messages will be shown",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.Bool("debug") {
			logrus.SetLevel(logrus.InfoLevel)
		}

		viper.SetConfigFile(c.String("config"))
		if err := viper.ReadInConfig(); err != nil {
			logrus.WithFields(logrus.Fields{
				"file":  c.String("config"),
				"error": err.Error(),
			}).Warnln("An error occurred while reading the configuration file. Using default configuration...")
			if _, err := os.Stat(c.String("config")); os.IsNotExist(err) {
				createConfigWhenNotExists()
			}
		} else {
			if duplicateErr := bot.CheckForDuplicateAliases(); duplicateErr != nil {
				logrus.WithFields(logrus.Fields{
					"issue": duplicateErr.Error(),
				}).Fatalln("An issue was discoverd in your configuration.")
			}
			createNewConfigIfNeeded()
			viper.WatchConfig()
		}

		if c.GlobalIsSet("server") {
			viper.Set("connection.address", c.String("server"))
		}
		if c.GlobalIsSet("port") {
			viper.Set("connection.port", c.String("port"))
		}
		if c.GlobalIsSet("username") {
			viper.Set("connection.username", c.String("username"))
		}
		if c.GlobalIsSet("password") {
			viper.Set("connection.password", c.String("password"))
		}
		if c.GlobalIsSet("channel") {
			viper.Set("defaults.channel", c.String("channel"))
		}
		if c.GlobalIsSet("cert") {
			viper.Set("connection.cert", c.String("cert"))
		}
		if c.GlobalIsSet("key") {
			viper.Set("connection.key", c.String("key"))
		}
		if c.GlobalIsSet("accesstokens") {
			viper.Set("connection.access_tokens", c.String("accesstokens"))
		}
		if c.GlobalIsSet("insecure") {
			viper.Set("connection.insecure", c.Bool("insecure"))
		}

		if err := DJ.Connect(); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatalln("An error occurred while connecting to the server.")
		}

		if viper.GetString("defaults.channel") != "" {
			defaultChannel := strings.Split(viper.GetString("defaults.channel"), "/")
			DJ.Client.Do(func() {
				DJ.Client.Self.Move(DJ.Client.Channels.Find(defaultChannel...))
			})
		}

		DJ.Client.Do(func() {
			DJ.Client.Self.SetComment(viper.GetString("defaults.comment"))
		})
		<-DJ.KeepAlive

		return nil
	}

	app.Run(os.Args)
}

func createConfigWhenNotExists() {
	configFile, err := Asset("config.yaml")
	if err != nil {
		logrus.Warnln("An error occurred while accessing config binary data. A new config file will not be written.")
	} else {
		filePath := os.ExpandEnv("$HOME/.config/mumbledj/config.yaml")
		os.Mkdir(os.ExpandEnv("$HOME/.config/mumbledj"), 0777)
		writeErr := ioutil.WriteFile(filePath, configFile, 0644)
		if writeErr == nil {
			logrus.WithFields(logrus.Fields{
				"file_path": filePath,
			}).Infoln("A default configuration file has been written.")
		} else {
			logrus.WithFields(logrus.Fields{
				"error": writeErr.Error(),
			}).Warnln("An error occurred while writing a new config file.")
		}
	}
}

func createNewConfigIfNeeded() {
	newConfigPath := os.ExpandEnv("$HOME/.config/mumbledj/config.yaml.new")

	// Check if we should write an updated config file to config.yaml.new.
	if assetInfo, err := AssetInfo("config.yaml"); err == nil {
		asset, _ := Asset("config.yaml")
		if configFile, err := os.Open(os.ExpandEnv("$HOME/.config/mumbledj/config.yaml")); err == nil {
			configInfo, _ := configFile.Stat()
			defer configFile.Close()
			if configNewFile, err := os.Open(newConfigPath); err == nil {
				defer configNewFile.Close()
				configNewInfo, _ := configNewFile.Stat()
				if assetInfo.ModTime().Unix() > configNewInfo.ModTime().Unix() {
					// The config asset is newer than the config.yaml.new file.
					// Write a new config.yaml.new file.
					ioutil.WriteFile(os.ExpandEnv(newConfigPath), asset, 0644)
					logrus.WithFields(logrus.Fields{
						"file_path": newConfigPath,
					}).Infoln("An updated default configuration file has been written.")
				}
			} else if assetInfo.ModTime().Unix() > configInfo.ModTime().Unix() {
				// The config asset is newer than the existing config file.
				// Write a config.yaml.new file.
				ioutil.WriteFile(os.ExpandEnv(newConfigPath), asset, 0644)
				logrus.WithFields(logrus.Fields{
					"file_path": newConfigPath,
				}).Infoln("An updated default configuration file has been written.")
			}
		}
	}
}
