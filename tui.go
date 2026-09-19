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

type model struct {
	width      int
	height     int
	lastUpdate time.Time

	processTable table.Model
	tableStyle   table.Styles
	baseStyle    lipgloss.Style
	viewStyle    lipgloss.Style

	CpuUsage cpu.TimesStat
	MemUsage mem.VirtualMemoryStat
}
type TickMsg time.Time

func convertBytes(b uint64) (string, string) {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	convert := float64(b)
	i := 0

	for convert >= 1024 && i < len(units)-1 {
		convert /= 1024
		i++
	}

	return fmt.Sprintf("%.1f", convert), units[i]
}

func (m model) View() string {
	column := m.baseStyle.Width(m.width).Padding(1, 0, 0, 0).Render
	content := m.baseStyle.
		Width(m.width).
		Height(m.height).
		Render(
			lipgloss.JoinVertical(lipgloss.Left,
				column(m.ViewHeader()),
				column(m.ViewProcess()),
			),
		)

	return content
}

var Color = struct {
	Border    lipgloss.Color
	Green     lipgloss.Color
	Secondary lipgloss.Color
}{
	Border:    lipgloss.Color("#5A5A5A"),
	Green:     lipgloss.Color("#04B575"),
	Secondary: lipgloss.Color("#3C3C3C"),
}

func (m model) ViewHeader() string {
	list := m.baseStyle.
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(Color.Border).
		Height(4).
		Padding(0, 1)

	listHeader := m.baseStyle.Bold(true).Render

	listItem := func(key string, value string, suffix ...string) string {
		finalSuffix := ""
		if len(suffix) > 0 {
			finalSuffix = suffix[0]
		}

		listItemValue := m.baseStyle.Align(lipgloss.Right).Render(fmt.Sprintf("%s%s", value, finalSuffix))
		listItemKey := func(key string) string {
			return m.baseStyle.Render(key + ":")
		}

		return fmt.Sprintf("%s %s", listItemKey(key), listItemValue)
	}

	return m.viewStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			fmt.Sprintf("Last update: %d milliseconds ago\n", time.Since(m.lastUpdate).Milliseconds()),
			lipgloss.JoinHorizontal(lipgloss.Top,
				// Progress Bars
				list.Render(
					lipgloss.JoinVertical(lipgloss.Left,
						listHeader("% Usage"),
						listItem("CPU", fmt.Sprintf("%s %.1f", ProgressBar(100-m.CpuUsage.Idle, m.baseStyle), 100-m.CpuUsage.Idle), "%"),
						listItem("MEM", fmt.Sprintf("%s %.1f", ProgressBar(m.MemUsage.UsedPercent, m.baseStyle), m.MemUsage.UsedPercent), "%"),
					),
				),

				// CPU Details
				list.Border(lipgloss.NormalBorder(), false).Render(
					lipgloss.JoinVertical(lipgloss.Left,
						listHeader("CPU"),
						listItem("user", fmt.Sprintf("%.1f", m.CpuUsage.User), "%"),
						listItem("sys", fmt.Sprintf("%.1f", m.CpuUsage.System), "%"),
						listItem("idle", fmt.Sprintf("%.1f", m.CpuUsage.Idle), "%"),
					),
				),
				list.Border(lipgloss.NormalBorder(), false).Render(
					lipgloss.JoinVertical(lipgloss.Left,
						listHeader(""),
						listItem("nice", fmt.Sprintf("%.1f", m.CpuUsage.Nice), "%"),
						listItem("iowait", fmt.Sprintf("%.1f", m.CpuUsage.Iowait), "%"),
						listItem("irq", fmt.Sprintf("%.1f", m.CpuUsage.Irq), "%"),
					),
				),
				list.Render(
					lipgloss.JoinVertical(lipgloss.Left,
						listHeader(""),
						listItem("softirq", fmt.Sprintf("%.1f", m.CpuUsage.Softirq), "%"),
						listItem("steal", fmt.Sprintf("%.1f", m.CpuUsage.Steal), "%"),
						listItem("guest", fmt.Sprintf("%.1f", m.CpuUsage.Guest), "%"),
					),
				),

				// MEM Details
				list.Border(lipgloss.NormalBorder(), false).Render(
					lipgloss.JoinVertical(lipgloss.Left,
						listHeader("MEM"),
						func() string {
							val, unit := convertBytes(m.MemUsage.Total)
							return listItem("total", val, unit)
						}(),
						func() string {
							val, unit := convertBytes(m.MemUsage.Used)
							return listItem("used", val, unit)
						}(),
						func() string {
							val, unit := convertBytes(m.MemUsage.Available)
							return listItem("free", val, unit)
						}(),
					),
				),
				list.Render(
					lipgloss.JoinVertical(lipgloss.Left,
						listHeader(""),
						func() string {
							val, unit := convertBytes(m.MemUsage.Active)
							return listItem("active", val, unit)
						}(),
						func() string {
							val, unit := convertBytes(m.MemUsage.Buffers)
							return listItem("buffers", val, unit)
						}(),
						func() string {
							val, unit := convertBytes(m.MemUsage.Cached)
							return listItem("cached", val, unit)
						}(),
					),
				),
			),
		),
	)
}

func (m model) ViewProcess() string {
	return m.viewStyle.Render(m.processTable.View())
}

func ProgressBar(percentage float64, baseStyle lipgloss.Style) string {
	totalBars := 20
	fillBars := int(percentage / 100 * float64(totalBars))

	if fillBars > totalBars {
		fillBars = totalBars
	} else if fillBars < 0 {
		fillBars = 0
	}

	filled := baseStyle.
		Foreground(Color.Green).
		Render(strings.Repeat("|", fillBars))

	empty := baseStyle.
		Foreground(Color.Secondary).
		Render(strings.Repeat("|", totalBars-fillBars))

	return baseStyle.Render(fmt.Sprintf("[%s%s]", filled, empty))
}
