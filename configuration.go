package renew

import (
  "github.com/AlexsJones/renew/fetcher"
  "time"
)
//Configuration ...
type Configuration struct {
  StartTime time.Time
  ApplicationDirectory string
  ApplicationBinaryPath string
  Process func()
  StateMonitor func(StatusCode)
  Fetcher fetcher.IFetcher
}
