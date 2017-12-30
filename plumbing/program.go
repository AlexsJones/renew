package plumbing

import (
	"log"
	"os/exec"
	"strings"
)

//RebuildAndInstall ...
func RebuildAndInstall(path string) error {
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
