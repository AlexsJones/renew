package renew

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/kardianos/osext"
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
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		p := path.Dir(filename)
		c.ApplicationDirectory = p
	}
	osex, err := osext.Executable()
	if err != nil {
		fmt.Println("An error occured with binary location search")
		os.Exit(1)
	}
	c.ApplicationBinaryPath = osex
	c.StartTime = time.Now()

	go func() {
		c.Fetcher.Init()
		c.StateMonitor(RUNNING)
		for {
			if c.Fetcher.ShouldRun() {
				c.StateMonitor(FETCHING)
				if err := c.Fetcher.Perform(); err != nil {
					c.StateMonitor(FAILURE)
				} else {
					c.StateMonitor(UPDATEFETCHED)
				}
			}
			time.Sleep(time.Second)
		}
	}()

	c.Process()

	os.Exit(0)
}
