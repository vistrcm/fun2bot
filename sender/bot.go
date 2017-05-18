package sender

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// initialize http client
var cookieJar, _ = cookiejar.New(nil)
var httpClient = &http.Client{
	Timeout: time.Minute * 1,
	Jar:     cookieJar,
}

// send http request to make bot send message to chat
func SendToURL(botURL, room, message string) (string, error) {
	resp, err := httpClient.PostForm(botURL, url.Values{"room": {room}, "message": {message}})

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	s := string(bytes)
	return s, nil

}
