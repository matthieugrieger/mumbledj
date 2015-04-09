/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/youtube/song.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package youtube

type Song struct {
	submitter    string
	title        string
	id           string
	duration     string
	thumbnailUrl string
	skippers     []string
	playlist     *Playlist
	dontSkip     bool
}

func NewSong(user, id string, playlist *Playlist) (*Song, error) {

}
