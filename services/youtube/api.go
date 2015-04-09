/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/youtube/api.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package youtube 

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"

    "github.com/jmoiron/jsonq"
)

// Collection of metadata for a YouTube video.
type YouTubeVideo struct {
  id string
  title string
  duration string
  secondsDuration string
  thumbnail string
}

// Collection of metadata for a YouTube playlist.
type YouTubePlaylist struct {
  id string
  title string
  duration string
  secondsDuration string
}

// Retrieves the metadata for a new YouTube video, and creates and returns a
// YouTubeVideo type.
func GetYouTubeVideo(id string) (*YouTubeVideo, error) {
  url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet%2CcontentDetails&id=%s&key=%s",
    id, os.Getenv("YOUTUBE_API_KEY"))
  jsonString := ""

  if response, err := http.Get(url); err == nil {
    defer response.Body.Close()
    if response.StatusCode == 200 {
      if body, err := ioutil.ReadAll(response.Body); err == nil {
        jsonString = string(body)
      }
    } else {
      if response.StatusCode == 403 {
        return nil, errors.New("Invalid API key supplied.")
      } else {
        return nil, errors.New("Invalid YouTube ID supplied.")
      }
    }
  } else {
    return nil, errors.New("An error occurred while receiving HTTP GET response.")
  }

  jsonData := map[string]interface{}{}
  decoder := json.NewDecoder(strings.NewReader(jsonString))
  decoder.Decode(&jsonData)
  jq := jsonq.NewQuery(jsonData)

  title, _ := jq.String("items", "0", "snippet", "title")
  thumbnail, _ := jq.String("items", "0", "snippet", "thumbnails", "high", "url")
  duration, _ := jq.String("items", "0", "contentDetails", "duration")

  minutes := int(duration[2 : strings.Index(duration, "M")])
  seconds := int(duration[strings.Index(duration, "M")+1 : len(duration)-1])
  totalSeconds := (minutes * 60) + seconds
  durationString := fmt.Sprintf("%d:%d", minutes, seconds)

  video := &YoutubeVideo {
    id: id,
    title: title,
    duration: durationString,
    secondsDuration: totalSeconds,
    thumbnail: thumbnail,
  }
  return video, nil
}
