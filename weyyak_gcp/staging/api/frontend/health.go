package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
)

func health() string {
	percent, _ := cpu.Percent(time.Second, true)
	load, _ := load.Avg()

	return fmt.Sprintf("CPU: %v\n Load: %v\n", percent, load)
}
