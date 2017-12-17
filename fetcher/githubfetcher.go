package fetcher

import "time"

//GithubFetcher for retrieving and updating remote golang projects
type GithubFetcher struct {
	Interval time.Duration
}
