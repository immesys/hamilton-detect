package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/immesys/spawnpoint/objects"
	"github.com/immesys/spawnpoint/spawnclient"
	bw "gopkg.in/immesys/bw2bind.v5"
)

const TIMEOUT = 5 * time.Second

var bwc *bw.BW2Client
var sc *spawnclient.SpawnClient

func main() {
	bwc = bw.ConnectOrExit("")
	bwc.SetEntityFromEnvironOrExit()
	sc, _ = spawnclient.NewFromBwClient(bwc)
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
	config := &objects.SvcConfig{
		ServiceName: "lifx-controller",
		Entity:      "lifx.ent",
		Image:       "immesys/eopdemo",
		MemAlloc:    "512M",
		CPUShares:   512,
	}
	_, err := sc.DeployService("", config, "ucberkeley/eop/spawnpoint/showroom", "lifx-controller")
	if err != nil {
		fmt.Printf("Failed to deploy service: %v\n", err)
	}

}
func hamiltonLeft() {
	sc.StopService("ucberkeley/eop/spawnpoint/showroom", "lifx-controller")
}
