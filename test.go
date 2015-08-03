package main

import (
	"fmt"
	"github.com/layeh/gumble/gumble"
	"os"
	"time"
)

type TestSettings struct {
	password     string
	ip           string
	port         string
	accesstokens []string
}

var test TestSettings

func Test(password, ip, port string, accesstokens []string) {
	test = TestSettings{
		password:     password,
		ip:           ip,
		port:         port,
		accesstokens: accesstokens,
	}
	test.testYoutubeSong()
}

func (t TestSettings) createClient(uname string) *gumble.Client {
	config := gumble.Config{
		Username: uname,
		Password: t.password,
		Address:  t.ip + ":" + t.port,
		Tokens:   t.accesstokens,
	}
	config.TLSConfig.InsecureSkipVerify = true
	client := gumble.NewClient(&config)
	return client
}

func (t TestSettings) testYoutubeSong() {
	dummyClient := t.createClient("dummy")
	if err := dummyClient.Connect(); err != nil {
		panic(err)
	}

	dj.client.Request(gumble.RequestUserList)
	time.Sleep(time.Second * 5)
	dummyUser := dj.client.Users.Find("dummy")
	if dummyUser == nil {
		fmt.Printf("User does not exist, printing users\n")
		for _, user := range dj.client.Users {
			fmt.Printf(user.Name + "\n")
		}
		fmt.Printf("End of user list\n")
		os.Exit(1)
	}

	// Don't judge, I used the (autogenerated) Top Tracks for United Kingdom playlist
	songs := map[string]string{
		"http://www.youtube.com/watch?v=QcIy9NiNbmo":  "Taylor Swift - Bad Blood ft. Kendrick Lamar",
		"https://www.youtube.com/watch?v=vjW8wmF5VWc": "Silentó - Watch Me (Whip/Nae Nae) (Official)",
		"http://youtu.be/nsDwItoNlLc":                 "Tinie Tempah ft. Jess Glynne - Not Letting Go (Official Video)",
		"https://youtu.be/hXTAn4ELEwM":                "Years & Years - Shine",
		"http://youtube.com/watch?v=RgKAFK5djSk":      "Wiz Khalifa - See You Again ft. Charlie Puth [Official Video] Furious 7 Soundtrack",
		"https://youtube.com/watch?v=qWWSM3wCiKY":     "Calvin Harris & Disciples - How Deep Is Your Love (Audio)",
		"http://www.youtube.com/v/yzTuBuRdAyA":        "The Weeknd - The Hills",
		"https://www.youtube.com/v/cNw8A5pwbVI":       "Pia Mia - Do It Again ft. Chris Brown, Tyga",
	}

	for url, title := range songs {
		err := add(dummyUser, url)
		if err != nil {
			fmt.Printf("For: %s; Expected: %s; Got: %s\n", url, title, err.Error())
		} else if dj.queue.CurrentSong().Title() != title {
			fmt.Printf("For: %s; Expected: %s; Got: %s\n", url, title, dj.queue.CurrentSong().Title())
		}

		time.Sleep(time.Second * 10)
		skip(dummyUser, false, false)
	}

	os.Exit(0)
	dummyClient.Disconnect()
}
