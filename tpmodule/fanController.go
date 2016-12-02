package tpmodule

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var fanPath = "/proc/acpi/ibm/fan"

// Fan struct represents the content of /proc/acpi/ibm/fan at a given time. One need to call its "update" method to populate or update.
type Fan struct {
	Status string
	Speed  int
	Level  string
}

// Update struct with data from the module
func (f *Fan) Update() {
	fPointer, _ := os.Open(fanPath)
	defer fPointer.Close()
	scanner := bufio.NewScanner(fPointer)
	scanner.Split(bufio.ScanLines)
	data := make([]string, 3)
	for i := 0; i < 3; i++ {
		scanner.Scan()
		s := strings.Split(scanner.Text(), ":")
		k := strings.TrimSpace(s[1])
		data[i] = k
	}

	f.Status = data[0]
	sp, _ := strconv.Atoi(data[1])
	f.Speed = sp
	f.Level = data[2]
}

// ToggleFan toggles the fan between enabled and disabled state
func (f *Fan) ToggleFan() {
	f.Update()
	switch f.Status {
	case "enabled":
		dest := []byte("disable")
		err := ioutil.WriteFile(fanPath, dest, 0644)
		if err != nil {
			panic(err)
		}
	case "disabled":
		dest := []byte("enable")
		err := ioutil.WriteFile(fanPath, dest, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// SetLevel sets the fan speed level to the "level" argument
func (f *Fan) SetLevel(level int) {
	l := "level " + strconv.Itoa(level)
	dest := []byte(l)
	err := ioutil.WriteFile(fanPath, dest, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// SetAutoLevel sets the fan speed to the "auto" level, meaning that the firmware will select the best fan level, ignoring what the configuration file says
func (f *Fan) SetAutoLevel() {
	dest := []byte("level auto")
	err := ioutil.WriteFile(fanPath, dest, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// SetFullSpeedLevel sets the fan speed to the maximum possible, defined by the firmware
func (f *Fan) SetFullSpeedLevel() {
	dest := []byte("level full-speed")
	err := ioutil.WriteFile(fanPath, dest, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
