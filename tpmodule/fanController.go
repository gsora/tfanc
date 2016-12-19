package tpmodule

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gsora/tfanc/config"
)

var fanPath = "/proc/acpi/ibm/fan"

// Fan struct represents the content of /proc/acpi/ibm/fan at a given time.
//One need to call its "update" method to populate or update.
type Fan struct {
	Status        string
	Speed         int
	Level         string
	Levels        FanLevels
	CurrentTarget config.Target
}

// FanLevels specify what are the minimum and maximum fan levels supported by the ThinkPad where
// tfanc is running.
type FanLevels struct {
	Min int
	Max int
}

// NewFan returns a new Fan{}
func NewFan() Fan {
	var f FanLevels
	f.DetectFanLevels()
	var r Fan
	r.Levels = f
	r.CurrentTarget = config.Target{MinTemp: 0, Level: 0}
	r.Update()

	return r
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

// DetectFanLevels reads from the thinkpad_acpi sysfs what are the supported
// fan levels of the ThinkPad who's running tfanc.
func (fl *FanLevels) DetectFanLevels() {
	path := "/proc/acpi/ibm/fan"
	d, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	cmds := strings.Split(strings.TrimSpace(string(d)), " ")
	levels := strings.Split(cmds[4][:len(cmds[4])-1], "-")
	fl.Min, _ = strconv.Atoi(levels[0])
	fl.Max, _ = strconv.Atoi(levels[1])
}
