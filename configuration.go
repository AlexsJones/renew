package renew

import (
	"time"

	"github.com/AlexsJones/renew/fetcher"
)

//Configuration ...
type Configuration struct {
	StartTime             time.Time
	ApplicationGoPath     string
	applicationBinaryPath string
	ApplicationArguments  []string
	Process               func()
	StateChange           chan StatusCode
	Fetcher               fetcher.IFetcher
}
