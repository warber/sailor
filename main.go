package main

import (
	"flag"
	"fmt"
	"github.com/warber/sailor/ui"
	"os"
)

func main() {

	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		if args[0] == "auth" {
			//todo
		} else {
			bailout()
		}
	}
	if len(args) == 0 {
		p := ui.NewProgram()

		if err := p.Start(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	} else {
		bailout()
	}
}

func bailout() {
	fmt.Println("Usage: sailor [auth]")
	flag.PrintDefaults()
	os.Exit(1)
}
