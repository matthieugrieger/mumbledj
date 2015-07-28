package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"

	"github.com/layeh/gumble/gumble"
)

var client_token = new(map[string]string)
var token_client = new(map[string]string)
var external_ip = ""

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	var uname = token_client[r.URL.Path[1:]]
	if uname == nil {
		fmt.Fprintf(w, "I don't know you")
	} else {
		fmt.Fprintf(w, "Hi there, I love %s!", uname)
	}
}

func Webserver() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9563", nil)
	rand.Seed(time.Now().UnixNano())
}

func GetWebAddress(user *gumble.User) {
	if client_token[user.Name] != nil {
		token_client[client_token[user.Name]] = nil
	}
	client_token[user.Name] = randSeq(10)
	token_client[client_token[user.Name]] = user.Name
	dj.SendPrivateMessage(user, fmt.Sprintf(WEB_ADDRESS, getIP(), client_token[user.Name], getIP(), client_token[user.Name]))
}

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

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
