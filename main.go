package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/websocket"
)

var token = flag.String("t", "", "Token for bot")

func main() {
	flag.Parse()
	posturl := "https://slack.com/api/rtm.start"
	form := url.Values{}
	form.Set("token", *token)

	req, reqErr := http.NewRequest("POST", posturl, bytes.NewBufferString(form.Encode()))
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}

	resp, errDo := client.Do(req)
	if errDo != nil {
		fmt.Println(errDo)
	}
	defer resp.Body.Close()
	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		fmt.Println(errRead)
	}

	result := make(map[string]interface{})
	json.Unmarshal(body, &result)
	wssURL := result["url"].(string)
	dialWebSocket(wssURL)
}

type SlackMessage struct {
	Type      string  `json:"type"`
	Channel   string  `json:"channel"`
	User      string  `json:"user"`
	Text      string  `json:"text"`
	TimeStamp float64 `json:"ts"`
	Team      string  `json:"team"`
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
				sm := SlackMessage{}
				json.Unmarshal(msg[:n], &sm)
				fmt.Println("Msg:", sm.Text)
			}
			fmt.Printf("Type:%v\n", string(r["type"].(string)))
		}
		fmt.Printf("Received: %s.\n", msg[:n])
	}

}
