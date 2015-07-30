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

type WebServer struct {
	port         int
	client_token map[*gumble.User]string
	token_client map[string]*gumble.User
}

type Page struct {
	siteUrl string
	token   string
}

var external_ip = ""

func Webserver(port int) *WebServer {
	var webserver = WebServer{port, make(map[*gumble.User]string), make(map[string]*gumble.User)}
	http.HandleFunc("/", webServer.homepage)
	http.HandleFunc("/add", webserver.add)
	http.HandleFunc("/volume", webserver.volume)
	http.HandleFunc("/skip", webserver.skip)
	http.ListenAndServe(":"+port, nil)
	rand.Seed(time.Now().UnixNano())
	return &webserver
}

func (w WebServer) homepage(w http.ResponseWriter, r *http.Request) {
	var uname = token_client[r.URL.Path[1:]]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, Page{"http://" + getIp() + ":" + w.port + "/", r.URL.Path[1:]})
	}
}

func (w WebServer) add(w http.ResponseWriter, r *http.Request) {
	var uname = token_client[r.FormValue("token")]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		add(uname, html.UnescapeString(r.FormValue("value")))
	}
}

func (w WebServer) volume(w http.ResponseWriter, r *http.Request) {
	var uname = token_client[r.FormValue("token")]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		var url = html.UnescapeString(r.FormValue("value"))
		add(uname, url)
	}
}

func (w WebServer) skip(w http.ResponseWriter, r *http.Request) {
	var uname = token_client[r.FormValue("token")]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		var url = html.UnescapeString(r.FormValue("value"))
		add(uname, url)
	}
}

func (w WebServer) GetWebAddress(user *gumble.User) {
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
