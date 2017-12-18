package renew

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/kardianos/osext"
)

func restart(c *Configuration) {
	log.Println("Restarting now...")
	args := []string{
		"-renew"}
	cmd := exec.Command(c.ApplicationBinaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatalf("renew: Failed to launch, error: %v", err)
	}
}

func signalHandler() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitChan := make(chan int)
	go func() {
		for {
			s := <-signalChan
			switch s {
			// kill -SIGHUP XXXX
			case syscall.SIGHUP:
				fmt.Println("hungup")

			// kill -SIGINT XXXX or Ctrl+c
			case syscall.SIGINT:
				fmt.Println("Warikomi")

			// kill -SIGTERM XXXX
			case syscall.SIGTERM:
				fmt.Println("force stop")
				exitChan <- 0

			// kill -SIGQUIT XXXX
			case syscall.SIGQUIT:
				fmt.Println("stop and core dump")
				exitChan <- 0

			default:
				fmt.Println("Unknown signal.")
				exitChan <- 1
			}
		}
	}()

	code := <-exitChan
	os.Exit(code)
}

//Run ...
func Run(c *Configuration) {

	//Capture child process flags
	var renewChild bool
	flag.BoolVar(&renewChild, "renew", false, "Process has a parent")
	flag.Parse()
	if renewChild {
		log.Println("Terminating parent process")
		parent := syscall.Getppid()
		log.Printf("renew: Killing parent pid: %v", parent)
		syscall.Kill(parent, syscall.SIGTERM)
	}

	//

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

	go func() {
		c.Fetcher.Init()
		c.StateMonitor(RUNNING)
		for {
			if c.Fetcher.ShouldRun() {
				// c.StateMonitor(FETCHING)
				// if err := c.Fetcher.Perform(); err != nil {
				// 	c.StateMonitor(FAILURE)
				// } else {
				// 	c.StateMonitor(UPDATEFETCHED)
				// 	//apply update
				// }

				restart(c)
			}
			time.Sleep(time.Second)
		}
	}()

	c.Process()

	os.Exit(0)
}
