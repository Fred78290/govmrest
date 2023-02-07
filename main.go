package main

import (
	"os"

	_ "github.com/Fred78290/govmrest/vm"
	"github.com/vmware/govmomi/govc/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
