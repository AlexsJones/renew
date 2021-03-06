# renew

[![Go Report Card](https://goreportcard.com/badge/github.com/AlexsJones/renew)](https://goreportcard.com/report/github.com/AlexsJones/renew)

Golang project self updater

<img src="https://i.imgur.com/Ll0gTjt.png" width="100"/>


The purpose here is to let your applications auto update.
I've built an interface for fetching from remotes; currently supporting github.


Below the example shows how to modify your main to add the renew implementation.
The `Process` field is a function pointer to your application code.

```go
func main() {

	stateChange := make(chan renew.StatusCode)

	go func() {
		for {
			select {
			case evt := <-stateChange:
				switch evt {
				case renew.RUNNING:
					fmt.Println("State has changed to running")
				case renew.FETCHING:
					fmt.Println("State has changed to fetched...")
				case renew.NOUPDATEFETCHED:
					fmt.Println("No update to fetch")
				case renew.UPDATEFETCHED:
					fmt.Println("State has changed to update fetched")
				case renew.RESTARTING:
					fmt.Println("-----restarting-----")
				}
			}

			time.Sleep(time.Second)
		}
	}()
	renew.Run(&renew.Configuration{
		Process: func() {
			fmt.Println("Started renew")
			time.Sleep(time.Second * 20)
			fmt.Println("Ended renew")
		},
		ApplicationGoPath:    "github.com/AlexsJones/renew",
		ApplicationArguments: []string{},

		StateChange: stateChange,
		Fetcher: &fetcher.GithubFetcher{
			Interval:         time.Second * 5,
			GithubRepository: "https://github.com/AlexsJones/renew.git",
		},
	})
}

```
