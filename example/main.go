package main

import (
	"fmt"
	"time"

	renew "github.com/AlexsJones/renew"
	"github.com/AlexsJones/renew/fetcher"
)

func main() {

	stateChange := make(chan renew.StatusCode)

	go func() {
		for {
			select {
			case evt := <-stateChange:
				switch evt {
				case renew.RUNNING:
					fmt.Println("State has changed to running")
				case renew.FETCHING:
					fmt.Println("State has changed to fetched...")
				case renew.NOUPDATEFETCHED:
					fmt.Println("No update to fetch")
				case renew.UPDATEFETCHED:
					fmt.Println("State has changed to update fetched")
				case renew.RESTARTING:
					fmt.Println("-----restarting-----")
				}
			}

			time.Sleep(time.Second)
		}
	}()
	renew.Run(&renew.Configuration{
		Process: func() {
			fmt.Println("Started renew")
			time.Sleep(time.Second * 20)
			fmt.Println("Ended renew")
		},
		ApplicationGoPath:    "github.com/AlexsJones/renew",
		ApplicationArguments: []string{},

		StateChange: stateChange,
		Fetcher: &fetcher.GithubFetcher{
			Interval:          time.Second * 5,
			DefaultOriginName: "origin",
			GithubRepository:  "https://github.com/AlexsJones/renew.git",
		},
	})
}
