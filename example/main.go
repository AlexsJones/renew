package main

import (
	"fmt"
	"time"

	renew "github.com/AlexsJones/renew"
	"github.com/AlexsJones/renew/fetcher"
)

func programStateChange(s renew.State) {
	switch s.StatusCode {
	case renew.RUNNING:
		fmt.Println("State has changed to running")
	case renew.FETCHING:
		fmt.Println("State has changed to fetched...")
	case renew.UPDATEFETCHED:
		fmt.Println("State has changed to update fetched")
	}
}

func mainStarted() {
	fmt.Println("Started renew")
	time.Sleep(time.Second * 30)
	fmt.Println("Ended renew")
}

func main() {
	renew.Run(&renew.Configuration{
		Process:      mainStarted,
		StateMonitor: programStateChange,
		Fetcher: &fetcher.GithubFetcher{
			Interval: time.Second * 5,
		},
	})
}
