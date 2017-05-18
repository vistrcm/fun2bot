package pts

import (
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/anaskhan96/soup/fetch"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"
)

const ptsURL = "http://linorgoralik.com/pts" // NOTE: no leading slash

// initialize http client
var cookieJar, _ = cookiejar.New(nil)
var httpClient = &http.Client{
	Timeout: time.Minute * 1,
	Jar:     cookieJar,
}

// GetHTML mimic to  soup.GET but with own http client
func get(url string) (string, error) {
	defer fetch.CatchPanic("Get()")
	resp, err := httpClient.Get(url)
	if err != nil {
		panic("Couldn't perform GET request to " + url)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Unable to read the response body")
	}
	s := string(bytes)
	return s, nil
}

// getLastStrip returns file name for last strip, like "pts247.html"
// parent tag <a> of "all.jpg" contains link to the last strip. Can help to calculate amount of strips.
// so let's get a parrent of this img and value of the href attribute
func getLastStrip(img soup.Root) (string, error) {

	for _, attr := range img.Pointer.Parent.Attr {
		if attr.Key == "href" {
			return attr.Val, nil
		}
	}
	return "", errors.New("Can't find last strip.")
}

// GetRandomImageURL return url of random PTS strip
func GetRandomImageURL() (string, error) {
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	max, err := getMaxStrip()
	if err != nil {
		return "wowow", err
	}
	generated := r.Intn(max + 1)
	return getStripImageURL(generated)
}

// getStrip returns url of specific strip
func getStripImageURL(id int) (string, error) {
	stripURL := fmt.Sprintf("%v/pts%d.html", ptsURL, id)

	resp, err := get(stripURL)
	if err != nil {
		return "can't get url", err
	}

	doc := soup.HTMLParse(resp)
	img := doc.Find("img")
	imageFile := img.Attrs()["src"]
	return fmt.Sprintf("%v/%v", ptsURL, imageFile), nil
}

func getMaxStrip() (int, error) {
	resp, err := get(ptsURL)
	if err != nil {
		return -1, err
	}

	doc := soup.HTMLParse(resp)
	img := doc.Find("img", "src", "all.jpg")

	lastFile, err := getLastStrip(img)
	if err != nil {
		return -2, err
	}

	// lastFile in format "pts247.html" need only number 247.
	strips := strings.TrimPrefix(strings.TrimSuffix(lastFile, ".html"), "pts")

	nStrips, err := strconv.Atoi(strips)
	if err != nil {
		return -3, err
	}

	return nStrips, nil

}
