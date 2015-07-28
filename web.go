package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/layeh/gumble/gumble"
)

var client_token = make(map[*gumble.User]string)
var token_client = make(map[string]*gumble.User)
var external_ip = ""

func Webserver() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/add", addSong)
	http.ListenAndServe(":9563", nil)
	rand.Seed(time.Now().UnixNano())
}

func homepage(w http.ResponseWriter, r *http.Request) {
	var uname = token_client[r.URL.Path[1:]]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		fmt.Fprintf(w, "Hang in there %s, I haven't made the website yet!", uname.Name)
	}
}

func addSong(w http.ResponseWriter, r *http.Request) {
	var uname = token_client[r.FormValue("token")]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		var url = html.UnescapeString(r.FormValue("url"))
		fmt.Fprintf(w, url)
	}
}

func GetWebAddress(user *gumble.User) {
	if client_token[user] != "" {
		token_client[client_token[user]] = nil
	}
	// dealing with collisions
	var firstLoop = true
	for firstLoop || token_client[client_token[user]] != nil {
		client_token[user] = randSeq(10)
		firstLoop = false
	}
	token_client[client_token[user]] = user
	dj.SendPrivateMessage(user, fmt.Sprintf(WEB_ADDRESS, getIP(), client_token[user], getIP(), client_token[user]))
}

// Gets the external ip address for the server
func getIP() string {
	if external_ip != "" {
		return external_ip
	} else {
		if response, err := http.Get("http://myexternalip.com/raw"); err == nil {
			defer response.Body.Close()
			if response.StatusCode == 200 {
				if body, err := ioutil.ReadAll(response.Body); err == nil {
					external_ip = strings.TrimSpace(string(body))
				}
			}
		}
		return external_ip
	}
}

// Generates a pseudorandom string of characters
func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
