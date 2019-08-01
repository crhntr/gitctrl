package main

import (
	"flag"
	"fmt"
)

func main() {

	var fn func()
	repo := flag.String("repo", ".", "path to repo")
	flag.Parse()

	switch flag.Arg(0) {
	case "statuses":
		fn = func() { statuses(*repo) }
	case "remote-origin-must-ssh":
		fn = func() { remoteOriginMustSSH(*repo) }
	case "multi-branch-view-file", "mbvf":
		fn = func() { multibranchedit(repo) }
	default:
		fn = func() { fmt.Println("command not known") }
	}

	fn()
}
