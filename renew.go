package renew

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kardianos/osext"
)

func watch(c *Configuration) (chan struct{}, error) {

	//log.Printf("watching %q\n", c.ApplicationDirectory)
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	done := make(chan struct{})
	go func() {
		for {
			select {
			case _ = <-w.Events:
				//log.Printf("watcher received: %+v", e)
				err = syscall.Exec(c.ApplicationBinaryPath, os.Args, os.Environ())
				if err != nil {
					log.Fatal(err)
				}
			case err = <-w.Errors:
				log.Printf("watcher error: %+v", err)
			case <-done:

				return
			}
		}
	}()
	err = w.Add(c.ApplicationDirectory)
	if err != nil {
		return nil, err
	}
	return done, nil
}

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
	//If statechange is an open channel then defer the close to the program exit
	if c.StateChange != nil {
		defer close(c.StateChange)
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
	c.ApplicationArguments = os.Args
	c.StartTime = time.Now()

	//Run the watch process
	watcher, err := watch(c)
	if err != nil {
		log.Fatal(err)
	}

	//Run the fetch cycle
	go func() {
		c.Fetcher.Init()
		if c.StateChange != nil {
			c.StateChange <- RUNNING
		}
		for {
			c.StateChange <- RUNNING
			if c.Fetcher.ShouldRun() {
				//log.Println("Performing fetch")
				if c.StateChange != nil {
					c.StateChange <- FETCHING
				}
				//Perform the fetch
				err := c.Fetcher.Perform()
				if err != nil {
					if c.StateChange != nil {
						c.StateChange <- FAILURE
					}
					fmt.Println(err.Error())
				}
				if c.StateChange != nil {
					c.StateChange <- UPDATEFETCHED
				}
			}
		}

	}()

	//Run the sub process
	c.Process()

	close(watcher)
	os.Exit(0)
}
