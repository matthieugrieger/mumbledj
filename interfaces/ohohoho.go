/*
 * MumbleDJ
 * By Matthieu Grieger
 * interfaces/command.go
 * Copyright (c) 2019 Reikion (MIT License)
 */

package interfaces

// Ohohoho is an interface that all commands must implement.
type Ohohoho interface {
	IsInterrupting() bool
	Stop() error
	EmptyStop() error
	PlaySample(sampleName string, howMany int) error
}
