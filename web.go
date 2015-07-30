package main

import (
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
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

func NewWebServer(port int) *WebServer {
	return &WebServer{
		port:         port,
		client_token: make(map[*gumble.User]string),
		token_client: make(map[string]*gumble.User),
	}
}

func (web *WebServer) makeWeb() {
	http.HandleFunc("/", web.homepage)
	http.HandleFunc("/add", web.add)
	http.HandleFunc("/volume", web.volume)
	http.HandleFunc("/skip", web.skip)
	http.ListenAndServe(":"+strconv.Itoa(web.port), nil)
	rand.Seed(time.Now().UnixNano())
}

func (web *WebServer) homepage(w http.ResponseWriter, r *http.Request) {
	var uname = web.token_client[r.URL.Path[1:]]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, Page{"http://" + getIP() + ":" + strconv.Itoa(web.port) + "/", r.URL.Path[1:]})
	}
}

func (web *WebServer) add(w http.ResponseWriter, r *http.Request) {
	var uname = web.token_client[r.FormValue("token")]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		add(uname, html.UnescapeString(r.FormValue("value")))
	}
}

func (web *WebServer) volume(w http.ResponseWriter, r *http.Request) {
	var uname = web.token_client[r.FormValue("token")]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		var url = html.UnescapeString(r.FormValue("value"))
		add(uname, url)
	}
}

func (web *WebServer) skip(w http.ResponseWriter, r *http.Request) {
	var uname = web.token_client[r.FormValue("token")]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		var url = html.UnescapeString(r.FormValue("value"))
		add(uname, url)
	}
}

func (website *WebServer) GetWebAddress(user *gumble.User) {
	Verbose("Port number: " + strconv.Itoa(web.port))
	if web.client_token[user] != "" {
		web.token_client[web.client_token[user]] = nil
	}
	// dealing with collisions
	var firstLoop = true
	for firstLoop || web.token_client[web.client_token[user]] != nil {
		web.client_token[user] = randSeq(10)
		firstLoop = false
	}
	web.token_client[web.client_token[user]] = user
	dj.SendPrivateMessage(user, fmt.Sprintf(WEB_ADDRESS, getIP(), web.client_token[user], getIP(), web.client_token[user]))
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
