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
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jmoiron/jsonq"
)

// PerformGetRequest does all the grunt work for a YouTube HTTPS GET request.
func PerformGetRequest(url string) (*jsonq.JsonQuery, error) {
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
			}
			return nil, errors.New("Invalid YouTube ID supplied.")
		}
	} else {
		return nil, errors.New("An error occurred while receiving HTTP GET response.")
	}

	jsonData := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(jsonString))
	decoder.Decode(&jsonData)
	jq := jsonq.NewQuery(jsonData)

	return jq, nil
}
