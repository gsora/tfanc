package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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
	configFile := flag.String("conf", "", "configuration file - optional, will load /etc/tfanc.conf by default")
	flag.Parse()

	// Make some security checks and (hopefully) return a config.Configuration
	conf := securityChecks(*configFile)

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

	cFan := tpmodule.NewFan()

	loop := time.NewTicker(time.Second * 1)

	for range loop.C {
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

	fan.CurrentTarget, _ = conf.WhatLevelDoIBelong(cpuTemp)
	fan.SetLevel(fan.CurrentTarget.Level)
}

func securityChecks(configFile string) config.Configuration {
	// check if user running this is root
	if os.Getuid() != 0 {
		log.Fatal("This program have to be executed with root rights.")
	}

	// then check if the kernel module is loaded with fan_control=1, else exit
	if err := tpmodule.IsModuleLoaded(); err != nil {
		log.Fatal(err)
	}

	var conf config.Configuration
	var err error

	// Load the config file
	if configFile == "" {
		conf, err = config.LoadConfig(config.ConfigFilePath)
	} else {
		conf, err = config.LoadConfig(configFile)
	}

	if err != nil {
		log.Fatal(err)
	}

	return conf
}
