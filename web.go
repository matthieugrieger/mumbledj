package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var client_token = new(map[string]string)
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
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func Webserver() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9563", nil)
}

func GetWebAddress(user *gumble.User) {
	dj.SendPrivateMessage(user, fmt.Sprintf(WEB_ADDRESS, getIP(), user.Name, getIP(), user.Name))
}

func getIP() string {
	if external_ip != "" {
		return external_ip
	} else {
		if response, err := http.Get("http://myexternalip.com/raw"); err == nil {
			defer response.Body.Close()
			if response.StatusCode == 200 {
				if body, err := ioutil.ReadAll(response.Body); err == nil {
					external_ip = string(body)
				}
			}
		}
		
		return external_ip
	}
}
