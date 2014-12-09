/*
 * MumbleDJ
 * By Matthieu Grieger
 * main.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

package main

import (
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/layeh/gumble/gumbleutil"
)

// MumbleDJ type declaration
type MumbleDJ struct {
	serverAddress		string
	serverPort		int
	username		string
	password		string
}
