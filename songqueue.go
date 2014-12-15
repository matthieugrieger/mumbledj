/*
 * MumbleDJ
 * By Matthieu Grieger
 * songqueue.go
 * Copyright (c) 2014 Matthieu Grieger (MIT License)
 */

 package main

 type SongQueue struct {
 	songs []Song
 }

 func NewSongQueue() *SongQueue {
 	return &SongQueue{}
 }