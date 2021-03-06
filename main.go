package main

import (
	"flag"
	"github.com/vistrcm/fun2bot/pts"
	"github.com/vistrcm/fun2bot/randomizer"
	"github.com/vistrcm/fun2bot/sender"
	"github.com/vistrcm/fun2bot/xkcd"
	"log"
	"math/rand"
	"time"
)

type fetchResult struct {
	fetched string
	err     error
}

type botResponse struct {
	response string
	err      error
}

// phraseToSend throw a dice. Decide which channel to use and retrieve phrase to send.
// values to decide hardcoded. for now about 30% goes to pts, other to xkcd
func phraseToSend() (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if r.Intn(100) >= 30 {
		return xkcd.RandomStrip()
	} else {
		return pts.GetRandomImageURL()
	}
}

func main() {
	// parse command line args
	botUrl := flag.String("botUrl", "http://host", "specify bot url")
	room := flag.String("room", "myroom", "specify room to send message to")
	flag.Parse()

	var next time.Time
	//var err error
	var fetchDone chan fetchResult // if non-nil, Fetch is running
	var postDone chan botResponse

	for {
		var fetchDelay time.Duration // initially 0 (no delay)
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}

		var startFetch <-chan time.Time

		if fetchDone == nil {
			startFetch = time.After(fetchDelay)
		}

		select {
		case <-startFetch:
			log.Printf("staring fetching")
			fetchDone = make(chan fetchResult, 1)
			go func() {
				fetched, err := phraseToSend()
				fetchDone <- fetchResult{fetched, err}
			}()

		case result := <-fetchDone:
			fetchDone = nil
			if result.err != nil {
				log.Fatalf("Something goes wrong with %v\n", result)
				break
			}

			postDone = make(chan botResponse, 1)

			go func() {
				log.Printf("strip url to send: %s\n", result.fetched)
				botResp, err := sender.SendToURL(*botUrl, *room, result.fetched)
				if err != nil {
					panic(err)
				}

				postDone <- botResponse{botResp, err}

			}()

			next = randomizer.GetNextTime()
			log.Printf("next fetch time: %v", next)

		case postResult := <-postDone:
			log.Printf("Bot response: %v\n", postResult.response)
		}
	}
}
