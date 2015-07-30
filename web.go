package main

import (
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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
	Site  string
	Token string
	User  string
}

var external_ip = ""

func NewWebServer(port int) *WebServer {
	rand.Seed(time.Now().UnixNano())
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
}

func (web *WebServer) homepage(w http.ResponseWriter, r *http.Request) {
	var uname = web.token_client[r.URL.Path[1:]]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		cwd, _ := os.Getwd()
		t, err := template.ParseFiles(filepath.Join(cwd, "./.mumbledj/web/index.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, Page{getIP() + ":" + strconv.Itoa(web.port), r.URL.Path[1:], uname.Name})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (web *WebServer) add(w http.ResponseWriter, r *http.Request) {
	var uname = web.token_client[r.FormValue("token")]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		add(uname, html.UnescapeString(r.FormValue("value")))
		fmt.Fprintf(w, "Success")
	}
}

func (web *WebServer) volume(w http.ResponseWriter, r *http.Request) {
	var uname = web.token_client[r.FormValue("token")]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		var vol = html.UnescapeString(r.FormValue("value"))
		volume(uname, vol)
		fmt.Fprintf(w, "Success")
	}
}

func (web *WebServer) skip(w http.ResponseWriter, r *http.Request) {
	var uname = web.token_client[r.FormValue("token")]
	if uname == nil {
		fmt.Fprintf(w, "Invalid Token")
	} else {
		value = html.UnescapeString(r.FormValue("value"))
		playlist, err := strconv.ParseBool(playlist)
		if err == nil {
			skip(uname, false, playlist)
			fmt.Fprintf(w, "Success")
		} else {
			fmt.Fprintf(w, "Invalid Value")
		}
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
