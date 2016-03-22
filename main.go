package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gsora/tfanc/tpmodule"
)

// The fan controller itself
var cFan tpmodule.Fan

func main() {
	// check if user running this is root
	if os.Getuid() != 0 {
		fmt.Println("This program have to be executed with root rights.")
		os.Exit(1)
	}

	// then check if the kernel module is loaded with fan_control=1, else exit
	if err := tpmodule.IsModuleLoaded(); err != nil {
		log.Fatal(err)
	}

	// calm down, golint.
	cFan.Update()
	fmt.Println(cFan.Status)
}
