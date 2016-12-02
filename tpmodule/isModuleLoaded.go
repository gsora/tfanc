package tpmodule

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"
)

// IsModuleLoaded checks if the kernel module is loaded with fan_control=1, else exit
func IsModuleLoaded() error {
	basePath := "/sys/module/thinkpad_acpi"
	_, err := ioutil.ReadDir(basePath)

	if err != nil {
		return errors.New("thinkpad_acpi module not loaded.")
	}

	fanControlFile, err := ioutil.ReadFile(basePath + "/parameters/fan_control")

	if err != nil {
		log.Fatal(err)

	}
	fanControlFlag := strings.TrimSpace(string(fanControlFile))

	if fanControlFlag != "Y" {
		return errors.New("thinkpad_acpi error: module loaded, but not initialized with \"fan_control=1\" parameter")
	}

	return nil
}
