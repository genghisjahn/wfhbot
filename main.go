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
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}
