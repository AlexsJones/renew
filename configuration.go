package renew

import (
	"time"

	"github.com/AlexsJones/renew/fetcher"
)

//Configuration ...
type Configuration struct {
	StartTime             time.Time
	ApplicationDirectory  string
	ApplicationBinaryPath string
	ApplicationArguments  []string
	Process               func()
	StateChange           chan StatusCode
	Fetcher               fetcher.IFetcher
}
