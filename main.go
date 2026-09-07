package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	cpuUsage, _ := GetCPUStats()

	cpuData, err := json.MarshalIndent(cpuUsage, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("CPU Percentage:")
	fmt.Println(string(cpuData))

	memUsage, _ := getMEMStats()

	memData, err := json.MarshalIndent(memUsage, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Memory:")
	fmt.Println(string(memData))
}
