package renew

import (
	"fmt"
	"os"
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

	go func() {

	}()

	c.Process()

	os.Exit(0)
}
