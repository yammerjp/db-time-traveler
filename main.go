package main

import (
	"fmt"

	"github.com/yammerjp/db-time-traveler/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	cmd.Execute(builtInformations())
}

func builtInformations() string {
	return fmt.Sprintf("version: %s\ncommit: %s\ndate: %s\nbuiltBy: %s\n", version, commit, date, builtBy)
}
