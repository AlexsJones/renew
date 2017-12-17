package renew

import "github.com/AlexsJones/renew/fetcher"
//Configuration ...
type Configuration struct {
  Process func()
  StateMonitor func(State)
  Fetcher fetcher.IFetcher
}
