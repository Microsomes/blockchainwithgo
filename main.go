package main

import (
	"os"
)

const bootstrapIRC = "http://www.dal.net:9090/"

func main() {

	defer os.Exit(0)

	cli := CommandLine{}

	cli.run()

}
