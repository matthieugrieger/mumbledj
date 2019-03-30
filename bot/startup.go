/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/startup.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// PerformStartupChecks executes the suite of startup checks that are run before the bot
// connects to the server.
func PerformStartupChecks() {
	logrus.WithFields(logrus.Fields{
		"num_services": fmt.Sprintf("%d", len(DJ.AvailableServices)),
	}).Infoln("Checking for availability of services...")

	for i := len(DJ.AvailableServices) - 1; i >= 0; i-- {
		if err := DJ.AvailableServices[i].CheckAPIKey(); err != nil {
			name := DJ.AvailableServices[i].GetReadableName()
			logrus.WithFields(logrus.Fields{
				"service": name,
				"error":   err.Error(),
			}).Warnln("A startup check discovered an issue. The service will be disabled.")

			// Remove service from enabled services.
			DJ.AvailableServices = append(DJ.AvailableServices[:i], DJ.AvailableServices[i+1:]...)
		}
	}

	if len(DJ.AvailableServices) == 0 {
		logrus.Fatalln("The bot cannot continue as no services are enabled.")
	}

	if err := checkYouTubeDLInstallation(); err != nil {
		logrus.Fatalln("youtube-dl is either not installed or is not discoverable in $PATH. youtube-dl is required to download audio.")
	}
	if viper.GetString("defaults.player_command") == "ffmpeg" {
		if err := checkFfmpegInstallation(); err != nil {
			logrus.Fatalln("ffmpeg is either not installed or is not discoverable in $PATH. If you would like to use avconv instead, change the defaults.player_command value in the configuration file.")
		}
	} else if viper.GetString("defaults.player_command") == "avconv" {
		if err := checkAvconvInstallation(); err != nil {
			logrus.Fatalln("avconv is either not installed or is not discoverable in $PATH. If you would like to use ffmpeg instead, change the defaults.player_command value in the configuration file.")
		}
	} else {
		logrus.Fatalln("The player command provided in the configuration file is invalid. Valid choices are: \"ffmpeg\", \"avconv\".")
	}

	if err := checkAria2Installation(); err != nil {
		logrus.Warnln("aria2 is not installed or is not discoverable in $PATH. The bot will still partially work, but some services will not work properly.")
	}

	if err := checkOpenSSLInstallation(); err != nil {
		logrus.Warnln("openssl is not installed or is not discoverable in $PATH. p12 certificate files will not work.")
	}
}

func checkYouTubeDLInstallation() error {
	logrus.Infoln("Checking YouTubeDL installation...")
	command := exec.Command("youtube-dl", "--version")
	if err := command.Run(); err != nil {
		return errors.New("youtube-dl is not properly installed")
	}
	return nil
}

func checkFfmpegInstallation() error {
	logrus.Infoln("Checking ffmpeg installation...")
	command := exec.Command("ffmpeg", "-version")
	if err := command.Run(); err != nil {
		return errors.New("ffmpeg is not properly installed")
	}
	return nil
}

func checkAvconvInstallation() error {
	logrus.Infoln("Checking avconv installation...")
	command := exec.Command("avconv", "-version")
	if err := command.Run(); err != nil {
		return errors.New("avconv is not properly installed")
	}
	return nil
}

func checkAria2Installation() error {
	logrus.Infoln("Checking aria2c installation...")
	command := exec.Command("aria2c", "-v")
	if err := command.Run(); err != nil {
		return errors.New("aria2c is not properly installed")
	}
	return nil
}

func checkOpenSSLInstallation() error {
	logrus.Infoln("Checking openssl installation...")
	command := exec.Command("openssl", "version")
	if err := command.Run(); err != nil {
		return errors.New("openssl is not properly installed")
	}
	return nil
}
