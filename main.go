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
)

var token = flag.String("t", "", "Token for bot")

func main() {
	flag.Parse()
	posturl := "https://slack.com/api/rtm.start"
	form := url.Values{}
	form.Set("token", *token)

	req, reqErr := http.NewRequest("POST", posturl, bytes.NewBufferString(form.Encode()))
	if reqErr != nil {
		log.Println(reqErr)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}

	resp, errDo := client.Do(req)
	if errDo != nil {
		log.Println(errDo)
	}
	defer resp.Body.Close()
	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		log.Println(errRead)
	}

	result := make(map[string]interface{})
	json.Unmarshal(body, &result)
	fmt.Println(result["url"])
}
