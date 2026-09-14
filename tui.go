package main

import (
    "fmt"
    "strings"
    "time"

    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/lipgloss"
    "github.com/shirou/gopsutil/v4/cpu"
    "github.com/shirou/gopsutil/v4/mem"
)

type model struct{
	width  int
	height int

	processTable table.Model
	tableStyle   table.Styles
	baseStyle    lipgloss.Style
	viewStyle    lipgloss.Style

	CpuUsage cpu.TimesStat
	MemUsage mem.VirtualMemoryStat
}
type TickMsg time.Time


func (m Model) View() string{

}


func ProgressBar(percentage float64, baseStyle lipgloss.Style) string{
	totalBars := 20
	fillBar := int(percentage / 100 * float64(totalBars))

	filled := baseStyle.
			Foreground(lipgloss.Color("#04B575")).
			Render(strings.Repeat("|", fillBar))

	empty := baseStyle.
			Foreground(lipgloss.Color("#3C3C3C")).
			Render(strings.Repeat("|", totalBars - fillBar))

	
	return baseStyle.Render(fmt.Sprintf("[%s%s]", filled, empty))
}