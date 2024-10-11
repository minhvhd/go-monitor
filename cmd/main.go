package main

import (
	"fmt"
	"github.com/minhvhd/go-monitor/internal/hardware"
)

func main() {
	system, err := hardware.GetSystem()
	if err != nil {
		fmt.Println(err)
		return
	}

	cpu, err := hardware.GetCPU()
	if err != nil {
		fmt.Println(err)
		return
	}

	disk, err := hardware.GetDisk()

	fmt.Println(system)
	fmt.Println(cpu)
	fmt.Println(disk)
}
