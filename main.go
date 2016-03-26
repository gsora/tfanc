package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/VividCortex/godaemon"
	"github.com/gsora/tfanc/config"
	"github.com/gsora/tfanc/tpmodule"
)

// The fan controller itself
var cFan tpmodule.Fan

func main() {
	// program parameters
	benchmark := flag.Bool("benchmark", false, "run all the possible fan levels, 10 seconds at time")
	foreground := flag.Bool("foreground", false, "don't fork in background, useful for debug purposes")
	flag.Parse()

	conf := securityChecks()

	if *foreground == false && *benchmark == false {
		godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	} else if *foreground == false && *benchmark == true {
		fmt.Println("Cannot fork to background while running in benchmark mode.")
	}

	// if -benchmark passed, run it
	if *benchmark == true {
		tpmodule.Benchmark(cFan)
		return
	}

	// calm down, golint.
	cFan.Update()
	fmt.Println(cFan.Status)
	fmt.Println(conf.Targets)
}

func securityChecks() config.Configuration {
	// check if user running this is root
	if os.Getuid() != 0 {
		fmt.Println("This program have to be executed with root rights.")
		os.Exit(1)
	}

	// then check if the kernel module is loaded with fan_control=1, else exit
	if err := tpmodule.IsModuleLoaded(); err != nil {
		log.Fatal(err)
	}

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	return conf
}
