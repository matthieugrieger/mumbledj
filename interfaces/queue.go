/*
 * MumbleDJ
 * By Matthieu Grieger
 * interfaces/queue.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package interfaces

// Queue is the interface which should be interacted with for queue operations.
// Using the Queue interface ensures thread safety.
type Queue interface {
	Length() int
	Reset()
	AppendTrack(Track) error
	InsertTrack(int, Track) error
	CurrentTrack() (Track, error)
	GetTrack(int) Track
	PeekNextTrack() (Track, error)
	Traverse(func(int, Track))
	ShuffleTracks()
	RandomNextTrack(bool)
	Skip()
	SkipPlaylist()
	PlayCurrent() error
	PauseCurrent() error
	ResumeCurrent() error
	StopCurrent() error
}
