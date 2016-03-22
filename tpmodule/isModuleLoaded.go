package tpmodule

import (
	"errors"
	"io/ioutil"
)

// IsModuleLoaded checks if the kernel module is loaded with fan_control=1, else exit
func IsModuleLoaded() error {
	basePath := "/sys/module/"
	modList, _ := ioutil.ReadDir(basePath)

	found := false

	for _, module := range modList {
		if module.Name() == "thinkpad_acpi" {
			found = true
			break
		}
	}

	if found == false {
		return errors.New("thinkpad_acpi error: module not loaded")
	}

	basePath = basePath + "thinkpad_acpi/parameters/"
	fanControlFile, _ := ioutil.ReadFile(basePath + "fan_control")
	fanControlFlag := string(fanControlFile)

	if fanControlFlag == "N" {
		return errors.New("thinkpad_acpi error: module loaded, but not initialized with \"fan_control=1\" parameter")
	}

	return nil
}
