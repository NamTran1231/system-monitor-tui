package main

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func printStats() {
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

	ProcessData, _ := GetProcesses(10)

	data, err := json.MarshalIndent(ProcessData, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
	if err != nil {
		fmt.Println(err)
		return
	}

}

func main() {
	style := lipgloss.NewStyle()

	fmt.Println("Test 0%:")
	fmt.Println(ProgressBar(0, style))

	fmt.Println("Test 32.8%:")
	fmt.Println(ProgressBar(32.8, style))
}
