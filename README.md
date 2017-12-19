# renew
Golang project self updater

<img src="https://i.imgur.com/Ll0gTjt.png" width="100"/>

```
func programStateChange(s renew.StatusCode) {
	switch s {
	case renew.RUNNING:
		fmt.Println("State has changed to running")
	case renew.FETCHING:
		fmt.Println("State has changed to fetched...")
	case renew.UPDATEFETCHED:
		fmt.Println("State has changed to update fetched")
	}
}

func mainStarted() {
	fmt.Println("Started renew")
	time.Sleep(time.Second * 20)
	fmt.Println("Ended renew")
}

func main() {
	renew.Run(&renew.Configuration{
		Process:      mainStarted,
		StateMonitor: programStateChange,
		Fetcher: &fetcher.GithubFetcher{
			Interval:         time.Second * 5,
			GithubRepository: "https://github.com/AlexsJones/renew.git",
		},
	})
}
```
