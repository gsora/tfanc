package tpmodule

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// GetCPUTemp returns the CPU temperature as read by thinkpad_acpi
func GetCPUTemp() int {
	dat, _ := ioutil.ReadFile("/proc/acpi/ibm/thermal")
	content := strings.Split(string(dat), " ")
	content = strings.Split(content[0], ":")
	temp := strings.TrimSpace(content[1])
	k, _ := strconv.Atoi(temp)
	return k
}
