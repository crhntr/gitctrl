package main

import (
  "log"
  "strings"
  "regexp"
  "fmt"
  git "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing"
)

func multibranchedit(wd string, args []string) {
    var branchRegex, filename string

    if len(args) > 0 {
      filename = args[0]
    }
    if len(args) > 1 {
      branchRegex = args[1]
    }

    repo, err := git.PlainOpen(wd)
    if err != nil {
      log.Fatal(err)
    }
    branchItr, err := repo.Branches()
    if err != nil {
      log.Fatal(err)
    }
    var branches []*plumbing.Reference
    err = branchItr.ForEach(func(branch *plumbing.Reference) error {
      name := strings.TrimPrefix(branch.Name().Short(), "refs/heads/")
      if match, err := regexp.MatchString(branchRegex, name); err != nil {
        log.Fatal(err)
      } else if match {
          branches = append(branches, branch)
      }
      return nil
    })
    if err != nil {
      log.Fatal(err)
    }

    for _, branch := range branches {
      fmt.Printf("%+v", branch.Name())
    }

    fmt.Println(filename)
}
