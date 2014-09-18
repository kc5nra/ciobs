package main

import (
	"ciobs/lib"
	"flag"
)

var agent = flag.Bool("agent", false, "act as agent")

func init() {
	flag.BoolVar(agent, "a", false, "act as agent")
}

func main() {
	flag.Parse()

	if !*agent {
		ciobs.StartServer()
	}
}