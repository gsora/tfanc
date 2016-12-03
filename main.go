package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/VividCortex/godaemon"
	"github.com/gsora/tfanc/config"
	"github.com/gsora/tfanc/tpmodule"
)

// The fan controller itself
var cFan tpmodule.Fan

// LastCPUTemp holds the last CPU temp read
var LastCPUTemp int

func main() {
	// program parameters
	benchmark := flag.Bool("benchmark", false, "run all the possible fan levels, 10 seconds at time")
	foreground := flag.Bool("foreground", false, "don't fork in background, useful for debug purposes")
	flag.Parse()

	// Make some security checks and (hopefully) return a config.Configuration
	conf := securityChecks()

	// If the user didn't want to run in benchmark mode, or in foreground mode,
	// fork into background.
	//
	// Otherwise, if he wanted to benchmark, run it.
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
	cFan := tpmodule.NewFan()
	fmt.Println(cFan.Levels)
	fmt.Println(cFan.Status)
	fmt.Println(conf.Targets)
	sort.Sort(config.ByRange(conf.Targets))
	fmt.Println(conf.Targets)

	loop := time.NewTicker(time.Second * 1)

	for _ = range loop.C {
		fanCycle(cFan, conf)
	}
}

/*

fanCycle is the main core of tfanc, it's the algorithm behind level switching.

It will be run inside a time.Ticker loop every second to guarantee an adequate
level of cooling for every application.

*/
func fanCycle(fan tpmodule.Fan, conf config.Configuration) {
	cpuTemp := tpmodule.GetCPUTemp()

	if LastCPUTemp > cpuTemp {
		if !(fan.CurrentTarget.MinTemp < cpuTemp && fan.CurrentTarget.MinTemp > cpuTemp) {
			fan.CurrentTarget = conf.NextTarget(fan.CurrentTarget)
			fan.SetLevel(fan.CurrentTarget.Level)
		}
	} else {
		fan.CurrentTarget = conf.PrevTarget(fan.CurrentTarget)
		fan.SetLevel(fan.CurrentTarget.Level)
	}

	LastCPUTemp = cpuTemp
}

func securityChecks() config.Configuration {
	// check if user running this is root
	if os.Getuid() != 0 {
		log.Fatal("This program have to be executed with root rights.")
	}

	// then check if the kernel module is loaded with fan_control=1, else exit
	if err := tpmodule.IsModuleLoaded(); err != nil {
		log.Fatal(err)
	}

	// Load the config file
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	return conf
}
