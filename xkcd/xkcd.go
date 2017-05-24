package xkcd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
	"math/rand"
	"errors"
)

const (
	xkcdURL = "https://xkcd.com"
	retryNum = 3
)

// initialize http client
var cookieJar, _ = cookiejar.New(nil)
var httpClient = &http.Client{
	Timeout: time.Minute * 1,
	Jar:     cookieJar,
}

// apiResponse represents API response from XKCD json api
type apiResponse struct {
	Num        int    `json:'num'`
	Title      string `json:'title'`
	SafeTitle  string `json:'safe_title'`
	Year       string `json:'year'`
	Month      string `json:'month'`
	Day        string `json:'day'`
	Img        string `json:'img'`
	Alt        string `json:'alt'`
	Transcript string `json:'transcript'`
	News       string `json:'news'`
	Link       string `json:'link'`
}

//send request to url. Handle errors somehow.
func request(url string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error happened %v, req: %v", err, req)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("error on sending request: %v, resp: %v", err, resp)
	}

	// check for status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		log.Fatalf("login failed. Response status: %s\n", resp.Status)
	}

	return resp
}

// lastStrip finds number of last strip
func lastStrip() (int, error) {
	lastStripURL := fmt.Sprintf("%v/info.0.json", xkcdURL)
	resp := request(lastStripURL)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error on reading: %v. body: %v", err, body)
		return -1, errors.New("error on reading xkcd response")
	}

	var apiResp = new(apiResponse)
	err = json.Unmarshal(body, apiResp)
	if err != nil {
		log.Fatalf("something happened during unmarshall: %v. f: %v", err, apiResp)
		return -2, errors.New("error parsing xkcd response")
	}
	return apiResp.Num, nil

}

// getStrip get specific strip
func getStrip(id int) (apiResponse, error) {
	stripURL := fmt.Sprintf("%v/%d/info.0.json", xkcdURL, id)
	log.Printf("getting stripe: %s", stripURL)

	// prepare var
	var apiResp = new(apiResponse)

	resp := request(stripURL)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error on reading: %v. body: %v", err, body)
		return *apiResp, errors.New("error on reading xkcd response")
	}


	err = json.Unmarshal(body, apiResp)
	if err != nil {
		log.Fatalf("something happened during unmarshall: %v. f: %v", err, apiResp)
		return *apiResp, errors.New("error parsing xkcd response")
	}
	return *apiResp, nil

}

// getRandomStrip return random XKCD strip.
func getRandomStrip() (result apiResponse, err error) {
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	max, _ := lastStrip()

	log.Printf("max stip num: %d\n", max)

	// GetStripImageURL can return error, so need to re-try with another num
	for i := 0; i <= retryNum; i++ {
		generated := r.Intn(max + 1)
		log.Printf("generated strip num: %d\n", generated)
		strip, err := getStrip(generated)
		if err == nil { // return if no error
			log.Printf("Strip: %d\n", strip.Num)
			return strip, err
		}
	}

	return *new(apiResponse), errors.New("Can not get image url")
}

// RandomStrip prepares string with image link and alt text
func RandomStrip() (string, error) {
	strip, err := getRandomStrip()
	if err != nil {
		log.Fatalf("error %s happened on getting random strip")
		return "", err
	}

	data := fmt.Sprintf("%s\n%q", strip.Img, strip.Alt)
	return data, nil

}