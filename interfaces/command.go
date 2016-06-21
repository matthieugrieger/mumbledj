/*
 * MumbleDJ
 * By Matthieu Grieger
 * interfaces/command.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package interfaces

import "github.com/layeh/gumble/gumble"

// Command is an interface that all commands must implement.
type Command interface {
	Aliases() []string
	Description() string
	IsAdminCommand() bool
	Execute(user *gumble.User, args ...string) (string, bool, error)
}
