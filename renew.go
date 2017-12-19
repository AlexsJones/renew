package renew

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/kardianos/osext"
)

func restart(c *Configuration) {
	log.Println("Restarting now...")

	currentPid := syscall.Getppid()
	log.Printf("Current pid %d\n", currentPid)
	pid, err := syscall.ForkExec(c.ApplicationBinaryPath, c.ApplicationArguments[1:], &syscall.ProcAttr{
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys:   &syscall.SysProcAttr{},
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = process.Release()
	if err != nil {
		log.Fatal(err.Error())
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Started with pid %d\n", pid)
	currentPid = syscall.Getpid()

	process, err = os.FindProcess(currentPid)
	if err != nil {
		log.Fatal(err.Error())
	}
	process.Signal(syscall.SIGTERM)
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
			case syscall.SIGHUP:
				fmt.Println("hungup")
			case syscall.SIGINT:
				fmt.Println("Warikomi")
			case syscall.SIGTERM:
				fmt.Println("force stop")
				exitChan <- 0
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

	pid := os.Getpid()
	log.Printf("Started with process id %d", pid)
	go signalHandler()

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
