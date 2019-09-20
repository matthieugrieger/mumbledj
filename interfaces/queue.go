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
	PrependTrack(Track) error
	InsertTrack(int, Track) error
	GetTrack(int) Track
	GetTrackNoWait(index int) Track
	PeekNextTrack() (Track, error)
	RemoveTrack(int) Track
	RemoveTrackIf(func(int, Track) bool) int
	Traverse(func(int, Track))
	ShuffleTracks()
	RandomNextTrack(bool)
}
