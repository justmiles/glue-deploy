package main

import (
	cmd "github.com/justmiles/glue-deploy/cmd"
)

// version of github.com/justmiles/glue-deploy. Overwritten during build
var version = "development"

func main() {
	cmd.Execute(version)
}
