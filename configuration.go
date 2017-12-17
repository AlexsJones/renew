package renew

import (
  "github.com/AlexsJones/renew/fetcher"
  "time"
)
//Configuration ...
type Configuration struct {
  StartTime time.Time
  Process func()
  StateMonitor func(State)
  Fetcher fetcher.IFetcher
}
