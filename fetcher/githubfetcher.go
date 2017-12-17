package fetcher

import (
	"log"
	"time"
)

//GithubFetcher for retrieving and updating remote golang projects
type GithubFetcher struct {
	Interval time.Duration
	LastRun  time.Time
}

//Perform updte check
func (g *GithubFetcher) Perform() {
	log.Println("Performing update with github fetcher")
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
}
