# renew
Golang project self updater

<img src="https://i.imgur.com/Ll0gTjt.png" width="100"/>

```
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
		StateChange: stateChange,
		Fetcher: &fetcher.GithubFetcher{
			Interval:         time.Second * 5,
			GithubRepository: "https://github.com/AlexsJones/renew.git",
		},
	})
}
```
