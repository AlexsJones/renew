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

	log.Printf("watching %q\n", c.ApplicationDirectory)
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	done := make(chan struct{})
	go func() {
		for {
			select {
			case e := <-w.Events:
				log.Printf("watcher received: %+v", e)
				err = syscall.Exec(c.ApplicationBinaryPath, os.Args, os.Environ())
				if err != nil {
					log.Fatal(err)
				}
			case err = <-w.Errors:
				log.Printf("watcher error: %+v", err)
			case <-done:
				log.Print("watcher shutting down")
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

		for {
			c.Fetcher.Init()
			if c.Fetcher.ShouldRun() {

				err := c.Fetcher.Perform()
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}

	}()

	//Run the sub process
	c.Process()

	close(watcher)
	os.Exit(0)
}
