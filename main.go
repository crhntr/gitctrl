package main

import "flag"

func main() {
	flag.Parse()
	command := flag.Arg(0)

	wd := "./"
	if len(flag.Args()) > 1 {
		wd = flag.Arg(1)
	}

	var args []string
	if a := flag.Args(); len(a) > 1 {
		args = a[1:]
	}

	switch command {
	case "statuses":
		statuses(wd)
	case "remote-origin-must-ssh":
		remoteOriginMustSSH(wd)
	case "multi-branch-edit", "mbe":
		multibranchedit(wd, args)
	}
}
