package main

import (
	"flkcli/cmd"
	_ "flkcli/cmd/login"
	_ "flkcli/cmd/set"
	_ "flkcli/cmd/setup"
)

func main() {
	cmd.Execute()
}
