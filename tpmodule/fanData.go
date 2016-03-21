package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Fan struct represents the content of /proc/acpi/ibm/fan at a given time. One need to call its "update" method to populate or update.
type Fan struct {
	Status string
	Speed  string
	Level  string
}

func (f *Fan) update() {

	fPointer, _ := os.Open("/proc/acpi/ibm/fan")
	defer fPointer.Close()
	scanner := bufio.NewScanner(fPointer)
	scanner.Split(bufio.ScanLines)
	for i := 0; i < 3; i++ {
		scanner.Scan()
		s := strings.Split(scanner.Text(), ":")
		s = strings.TrimSpace(s[1])
		// TODO: map values to struct
	}

}

func main() {
	var a Fan
	a.update()
}
