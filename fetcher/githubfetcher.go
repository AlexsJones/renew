package fetcher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	git "gopkg.in/src-d/go-git.v4"
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
	GithubRepository string
}

//Perform updte check
func (g *GithubFetcher) Perform(applicationBasePath string) error {

	b, err := exists(applicationBasePath)
	if err != nil {
		return err
	}
	if !b {
		return errors.New(".git not found in directory")
	}
	r, err := git.PlainOpen(applicationBasePath)
	if err != nil {
		return err
	}

	err = r.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		return err
	}
	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(ref.Hash())

	log.Printf("Updated to commit %s\n", commit)
	return nil
}

//ShouldRun ...
func (g *GithubFetcher) ShouldRun() bool {

	nextRunTime := g.LastRun.Add(g.Interval)

	if time.Now().After(nextRunTime) {
		now := time.Now()
		//log.Printf("Running now and updating next run to %s\n", time.Now().Add(g.Interval).String())
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
}
