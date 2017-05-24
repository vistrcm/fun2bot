package pts

import (
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"
)

const (
	ptsURL   = "http://linorgoralik.com/pts" // NOTE: no leading slash
	retryNum = 3                             // how many times to retry to get image
)

// initialize http client
var cookieJar, _ = cookiejar.New(nil)
var httpClient = &http.Client{
	Timeout: time.Minute * 1,
	Jar:     cookieJar,
}

// get mimic to  soup.GET but with own http client
func get(url string) (string, error) {
	//defer fetch.CatchPanic("Get()")
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
func getLastStripFile(img soup.Root) (string, error) {

	for _, attr := range img.Pointer.Parent.Attr {
		if attr.Key == "href" {
			return attr.Val, nil
		}
	}
	return "", errors.New("Can't find last strip.")
}

// GetRandomImageURL return url of random PTS strip
func GetRandomImageURL() (result string, err error) {
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	max, err := getMaxStrip()
	if err != nil {
		return "wowow", err
	}
	log.Printf("max stip num: %d\n", max)

	// GetStripImageURL can return error, so need to re-try with another num
	for i := 0; i <= retryNum; i++ {
		generated := r.Intn(max + 1)
		log.Printf("generated strip num: %d\n", generated)
		url, err := GetStripImageURL(generated)
		if err == nil { // return if no error
			log.Printf("Strip image url: %q\n", url)
			return url, err
		}
	}

	return "", errors.New("Can not get image url.")
}

// getStrip returns url of specific strip
func GetStripImageURL(id int) (result string, err error) {
	stripURL := fmt.Sprintf("%v/pts%d.html", ptsURL, id)

	log.Printf("getting image for strip: %s\n", stripURL)
	resp, err := get(stripURL)
	if err != nil {
		return "can't get url", err
	}

	doc := soup.HTMLParse(resp)

	img := doc.Find("img")

	// Find do not return error if error happened. So let's check if Pointer is set after Find().
	if img.Pointer == nil {
		return "", errors.New(fmt.Sprintf("cant' find image for strip %q", stripURL))
	} else {
		imageFile := img.Attrs()["src"]
		return fmt.Sprintf("%v/%v", ptsURL, imageFile), nil
	}

}

// getMaxStrip returns number of last strip
func getMaxStrip() (int, error) {
	resp, err := get(ptsURL)
	if err != nil {
		return -1, err
	}

	doc := soup.HTMLParse(resp)
	img := doc.Find("img", "src", "all.jpg")

	lastFile, err := getLastStripFile(img)
	if err != nil {
		return -2, err
	}
	log.Printf("last stip file name: %s\n", lastFile)

	// lastFile in format "pts247.html" need only number 247.
	strips := strings.TrimPrefix(strings.TrimSuffix(lastFile, ".html"), "pts")

	nStrips, err := strconv.Atoi(strips)
	if err != nil {
		return -3, err
	}

	return nStrips, nil

}
