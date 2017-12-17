package renew

import (
	"fmt"
	"os"
	"time"
)

//Run ...
func Run(c *Configuration) {
	if c == nil {
		fmt.Println("No configuration")
		os.Exit(1)
	}
	if c.Process == nil {
		fmt.Println("No process function has been defined")
		os.Exit(1)
	}
	if c.Fetcher == nil {
		fmt.Println("No fetch process configured")
		os.Exit(1)
	}
	c.StartTime = time.Now()

	go func() {
		c.Fetcher.Init()
		for {
			if c.Fetcher.ShouldRun() {
				c.Fetcher.Perform()
			}
			time.Sleep(time.Second)
		}
	}()

	c.Process()

	os.Exit(0)
}
