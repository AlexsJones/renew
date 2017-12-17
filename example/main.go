package main

import (
	"fmt"
	"time"

	renew "github.com/AlexsJones/renew"
	"github.com/AlexsJones/renew/fetcher"
)

func programStateChange(s renew.State) {
	fmt.Printf("State has changed to %s\n", s.Description)
}

func mainStarted() {
	fmt.Println("Started renew")

	fmt.Println("Ended renew")
}

func main() {
	renew.Run(&renew.Configuration{
		Process:      mainStarted,
		StateMonitor: programStateChange,
		Fetcher: &fetcher.GithubFetcher{
			Interval: time.Second,
		},
	})
}
