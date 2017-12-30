package renew

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/AlexsJones/renew/plumbing"
	"github.com/fsnotify/fsnotify"
	"github.com/kardianos/osext"
)

func watch(c *Configuration) (chan struct{}, error) {

	log.Printf("watching %q\n", c.ApplicationGoPath)
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
				if c.StateChange != nil {
					c.StateChange <- RESTARTING
				}
				err = syscall.Exec(c.applicationBinaryPath, os.Args, os.Environ())
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
	err = w.Add(c.ApplicationGoPath)
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
	if c.ApplicationGoPath == "" {
		fmt.Println("No ApplicationGoPath configured")
		os.Exit(1)
	}
	//Modify path to absolute
	c.ApplicationGoPath = path.Join(build.Default.GOPATH, "src", c.ApplicationGoPath)

	//If statechange is an open channel then defer the close to the program exit
	if c.StateChange != nil {
		defer close(c.StateChange)
	}
	//This runs against the binary path, you'll need to `go build` your application
	osex, err := osext.Executable()
	if err != nil {
		fmt.Println("An error occured with binary location search")
		os.Exit(1)
	}
	c.applicationBinaryPath = osex
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
				if c.StateChange != nil {
					c.StateChange <- FETCHING
				}
				//Perform the fetch
				updated, err := c.Fetcher.Perform(c.ApplicationGoPath)
				if err != nil {
					if c.StateChange != nil {
						c.StateChange <- FAILURE
					}
					fmt.Println(err.Error())
				}
				if updated {
					if c.StateChange != nil {
						c.StateChange <- UPDATEFETCHED
					}
					plumbing.RebuildAndInstall(c.ApplicationGoPath)
				} else {
					if c.StateChange != nil {
						c.StateChange <- NOUPDATEFETCHED
					}
				}
			}
		}

	}()

	//Run the sub process
	c.Process()

	close(watcher)
	os.Exit(0)
}
