package fetcher

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	origin = "origin"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

//GithubFetcher for retrieving and updating remote golang projects
type GithubFetcher struct {
	Interval time.Duration
	LastRun  time.Time
	//Fetcher specific -------
	GithubRepository  string
	DefaultOriginName string
}

func fetchHash() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return "", err
	}
	t := strings.TrimSpace(string(out))
	return t, nil
}

func fetchBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	t := strings.TrimSpace(string(out))
	return t, nil
}

//Perform updte check
func (g *GithubFetcher) Perform(applicationBasePath string) (bool, error) {
	dir, err := os.Getwd()
	if err != nil {
		return false, err
	}
	initialhash, err := fetchHash()
	if err != nil {
		return false, err
	}
	err = os.Chdir(applicationBasePath)
	if err != nil {
		return false, err
	}
	branch, err := fetchBranch()
	if err != nil {
		return false, err
	}
	_, err = exec.Command("git", "pull", "origin", branch).Output()
	if err != nil {
		return false, err
	}
	updatedHash, err := fetchHash()
	if err != nil {
		return false, err
	}
	err = os.Chdir(dir)
	if err != nil {
		return false, err
	}
	if strings.Compare(initialhash, updatedHash) == 0 {
		return false, nil
	}
	return true, err
}

//ShouldRun ...
func (g *GithubFetcher) ShouldRun() bool {

	nextRunTime := g.LastRun.Add(g.Interval)

	if time.Now().After(nextRunTime) {
		now := time.Now()
		g.LastRun = now
		return true
	}

	return false
}

//Init ...
func (g *GithubFetcher) Init() {
	g.LastRun = time.Now()
	if g.GithubRepository == "" {
		fmt.Println("No Github repository specified")
		os.Exit(1)
	}
	if g.DefaultOriginName == "" {
		g.DefaultOriginName = origin
	}
}
