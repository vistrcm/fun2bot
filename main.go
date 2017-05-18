package main

import (
	"flag"
	"github.com/vistrcm/fun2bot/pts"
	"github.com/vistrcm/fun2bot/randomizer"
	"github.com/vistrcm/fun2bot/sender"
	"log"
	"time"
)

type fetchResult struct {
	fetched string
	err     error
}

func main() {
	// parse command line args
	botUrl := flag.String("botUrl", "http://host", "specify bot url")
	room := flag.String("room", "myroom", "specify room to send message to")
	flag.Parse()

	var next time.Time
	//var err error
	var fetchDone chan fetchResult // if non-nil, Fetch is running

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
				fetched, err := pts.GetRandomImageURL()
				fetchDone <- fetchResult{fetched, err}
			}()

		case result := <-fetchDone:
			fetchDone = nil
			if result.err != nil {
				log.Fatalf("Something goes wrong with %v\n", result)
				break
			}

			log.Printf("strip url to send: %s\n", result.fetched)
			botResp, err := sender.SendToURL(*botUrl, *room, result.fetched)
			if err != nil {
				panic(err)
			}

			log.Printf("Bot response: %v\n", botResp)

			next = randomizer.GetNextTime()
			log.Printf("next fetch time: %v", next)
		}
	}
}
