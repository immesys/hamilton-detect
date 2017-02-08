package main

import (
	"sync"
	"time"

	bw "gopkg.in/immesys/bw2bind.v5"
)

const TIMEOUT = 5 * time.Second

func main() {
	bwc := bw.ConnectOrExit("")
	bwc.SetEntityFromEnvironOrExit()
	subchan := bwc.SubscribeOrExit(&bw.SubscribeParams{
		URI:       "amplab/brz/hbr003/s.hamilton/00126d0700000061/i.l7g/signal/raw",
		AutoChain: true,
	})
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			checkLeft()
		}
	}()
	for _ = range subchan {
		hamiltonHeard()
	}
}

var isIn bool
var stateMu sync.Mutex
var lastHeard time.Time

func hamiltonHeard() {
	stateMu.Lock()
	lastHeard = time.Now()
	if !isIn {
		isIn = true
		stateMu.Unlock()
		go hamiltonEntered()
	}
	stateMu.Unlock()
}
func checkLeft() {
	stateMu.Lock()
	if time.Now().Sub(lastHeard) > TIMEOUT {
		isIn = false
		stateMu.Unlock()
		go hamiltonLeft()
	}
	stateMu.Unlock()
}
func hamiltonEntered() {
	//JACK put your action here
}
func hamiltonLeft() {
	//JACK put your action here
}
