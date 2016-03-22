package tpmodule

import (
	"fmt"
	"time"
)

// Benchmark cycles through all the supported fan levels, reporting its speed
func Benchmark(f Fan) {
	fmt.Printf("Running benchmark...\n")
	for i := 0; i < 8; i++ {
		f.Update()
		fmt.Println("Setting level ", i, " was ", f.Level)
		f.SetLevel(i)
		time.Sleep(10 * time.Second)
		fmt.Println("Speed at level ", i, ": ", f.Speed)
		fmt.Printf("\n")
	}

	fmt.Println("Setting level \"auto\" for safety.")
	f.SetAutoLevel()
}
