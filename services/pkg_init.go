/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/pkg_init.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package services

import (
	"reik.pl/mumbledj/bot"
	"reik.pl/mumbledj/interfaces"
)

// DJ is an injected MumbleDJ struct.
var DJ *bot.MumbleDJ

// Services is a slice of enabled MumbleDJ services.
var Services []interfaces.Service

func init() {
	Services = []interfaces.Service{
		NewMixcloudService(),
		NewSoundCloudService(),
		NewYouTubeService(),
	}
}
