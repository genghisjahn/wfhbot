package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/websocket"
)

var token = flag.String("t", "", "Token for bot")
var baseURL = "https://slack.com/api/%v"

func main() {
	flag.Parse()

	f := url.Values{}

	body, cmdErr := doCommand("rtm.start", f)
	if cmdErr != nil {
		fmt.Println("ERROR:", cmdErr)
	}
	result := make(map[string]interface{})
	json.Unmarshal(body, &result)
	wssURL := result["url"].(string)
	dialWebSocket(wssURL)
}

type Message struct {
	Type      string  `json:"type"`
	Channel   string  `json:"channel"`
	User      string  `json:"user"`
	Text      string  `json:"text"`
	TimeStamp float64 `json:"ts"`
	Team      string  `json:"team"`
}

type UserType struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	User    string `json:"user"`
}

type Greeting struct {
	Type string `json:"type"`
}

func doCommand(command string, form url.Values) ([]byte, error) {
	posturl := fmt.Sprintf(baseURL, command)
	if form == nil {
		return nil, errors.New("Form is nil!")
	}
	form.Set("token", *token)
	req, reqErr := http.NewRequest("POST", posturl, bytes.NewBufferString(form.Encode()))
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}

	resp, errDo := client.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	defer resp.Body.Close()
	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}
	return body, nil
}

func dialWebSocket(wssURL string) {
	origin := "http://ms114.slack-msgs.com/"
	ws, err := websocket.Dial(wssURL, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	/*
		if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
			log.Fatal(err)
		}
	*/
	r := make(map[string]interface{})
	for {
		var msg = make([]byte, 512)
		var n int
		if n, err = ws.Read(msg); err != nil {
			log.Println(err)
			log.Println("Lost connection, restarting...")
			dialWebSocket(wssURL)
		}
		json.Unmarshal(msg[:n], &r)
		if r["type"] != nil {
			t := r["type"].(string)
			switch t {
			case "message":
				m := Message{}
				json.Unmarshal(msg[:n], &m)
				fmt.Println("Msg:", m.Text)
			case "hello":
				f := url.Values{}
				f.Set("name", "#general")
				doCommand("channels.join", f)
				f2 := url.Values{}
				f2.Set("channel", "U02LDM150")
				f2.Set("text", "test private message")
				f2.Set("username", "WFHBot")
				doCommand("chat.postMessage", f2)
			}
		}
		fmt.Printf("Received: %s.\n", msg[:n])
	}
}
