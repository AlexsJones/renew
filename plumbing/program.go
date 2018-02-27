package plumbing

import (
	"log"
	"os/exec"
	"strings"
	"errors"
	"runtime"
)

//RebuildAndInstall ...
func RebuildAndInstall(path string) error {
	if runtime.GOOS == "windows" {
    return errors.New("Does not support windows")
	}
	cmd := exec.Command("go", "build")
	cmd.Path = path
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	t := strings.TrimSpace(string(out))
	log.Println(t)
	cmd = exec.Command("go", "install")
	cmd.Path = path
	if err != nil {
		return err
	}
	t = strings.TrimSpace(string(out))
	log.Println(t)
	return nil
}
