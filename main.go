package main

import (
	"github.com/camptocamp/prometheus-orchestrators-sd/cmd"
)

var exitCode int

// Following variables are filled in by the build script
var posdVersion = "<<< filled in by build >>>"

func main() {
	cmd.Execute()
}
