package fetcher

import (
	"fmt"
	"log"
	"os"
	"time"
)

//GithubFetcher for retrieving and updating remote golang projects
type GithubFetcher struct {
	Interval time.Duration
	LastRun  time.Time
	//Fetcher specific -------
	GithubRepository string
}

//Perform updte check
func (g *GithubFetcher) Perform() error {
	log.Println("Performing update with github fetcher")

	return nil
}

//ShouldRun ...
func (g *GithubFetcher) ShouldRun() bool {

	nextRunTime := g.LastRun.Add(g.Interval)

	if time.Now().After(nextRunTime) {
		now := time.Now()
		log.Printf("Running now and updating next run to %s\n", time.Now().Add(g.Interval).String())
		g.LastRun = now
		return true
	}

	return false
}

//Init ...
func (g *GithubFetcher) Init() {
	g.LastRun = time.Now()
	if g.GithubRepository == "" {
		fmt.Println("No Github repository specified")
		os.Exit(1)
	}
}
