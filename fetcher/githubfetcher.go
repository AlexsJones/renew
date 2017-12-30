package fetcher

import (
	"fmt"
	"log"
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
	log.Printf("Application base path: %s\n", applicationBasePath)
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
	//Horrible usage of command here due to gitv4 having ssh passthrough issues
	out, err := exec.Command("git", "pull", "origin", branch).Output()
	if err != nil {
		return false, err
	}
	log.Println(string(out))
	updatedHash, err := fetchHash()
	if err != nil {
		return false, err
	}
	err = os.Chdir(dir)
	if err != nil {
		return false, err
	}
	if strings.Compare(initialhash, updatedHash) == 0 {
		log.Printf("%s %s\n", initialhash, updatedHash)
		return false, nil
	}
	log.Printf("%s %s\n", initialhash, updatedHash)
	return true, err
}

//ShouldRun ...
func (g *GithubFetcher) ShouldRun() bool {

	nextRunTime := g.LastRun.Add(g.Interval)

	if time.Now().After(nextRunTime) {
		now := time.Now()
		log.Printf("Running now and updating next run to %s\n", time.Now().Add(g.Interval).String())
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
