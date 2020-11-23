package main

import "fmt"

var (
	buildstamp  = "not set"
	buildnumber = "not set"
	githash     = "not set"
)

func showVersion() {
	fmt.Println("Git hash: ", githash)
	fmt.Println("Build time: ", buildstamp)
	fmt.Println("Build number: ", buildnumber)
}
