package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	git "gopkg.in/src-d/go-git.v4"
)

type DirStatus struct {
	Path   string
	Status string
}

func (ds DirStatus) String() string {
	return fmt.Sprintf("=========\n%s\n%s\n\n", ds.Path, ds.Status)
}

func statuses(wd string) {
	wg := sync.WaitGroup{}
	stats := make(chan DirStatus)

	printGitDirStats := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasSuffix(path, ".git") || strings.HasSuffix(path, "node_modules") {
			return filepath.SkipDir
		}

		stat, err := os.Stat(filepath.Join(path, ".git"))
		if err != nil || !stat.IsDir() {
			return nil
		}

		go func() {
			wg.Add(1)
			defer wg.Done()

			repo, err := git.PlainOpen(path)
			if err != nil {
				if err == git.ErrRepositoryNotExists {
					return
				}
				log.Println(err)
				return
			}
			wtree, err := repo.Worktree()
			if err != nil {
				log.Println(err)
				return
			}
			status, err := wtree.Status()
			if err != nil {
				log.Println(err)
				return
			}
			if !status.IsClean() {
				stats <- DirStatus{path, status.String()}
			}
		}()

		return filepath.SkipDir
	}

	go func() {
		for stat := range stats {
			stat.Path = strings.TrimPrefix(stat.Path, wd)
			fmt.Println(stat)
		}
	}()
	err := filepath.Walk(wd, printGitDirStats)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
