package main

import "flag"

func main() {
	flag.Parse()
	command := flag.Arg(0)

	wd := "./"
	if len(flag.Args()) > 1 {
		wd = flag.Arg(1)
	}

	switch command {
	case "statuses":
		statuses(wd)
	case "remote-origin-must-ssh":
		remoteOriginMustSSH(wd)
	}
}
