package main

import (
	"math"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

func GetCPUStats() (cpu.TimesStat, error) {
	times, err := cpu.Times(false)

	if err != nil {
		return cpu.TimesStat{}, err
	}

	if len(times) == 0 {
		return cpu.TimesStat{}, nil
	}

	curTimes := times[0]

	total := curTimes.User +
		curTimes.System +
		curTimes.Idle +
		curTimes.Nice +
		curTimes.Iowait +
		curTimes.Irq +
		curTimes.Softirq +
		curTimes.Steal +
		curTimes.Guest +
		curTimes.GuestNice

	if total == 0 {
		return cpu.TimesStat{}, nil
	}

	curTimes.User = math.Round((curTimes.User/total)*100*100) / 100
	curTimes.System = math.Round((curTimes.System/total)*100*100) / 100
	curTimes.Idle = math.Round((curTimes.Idle/total)*100*100) / 100
	curTimes.Nice = math.Round((curTimes.Nice/total)*100*100) / 100
	curTimes.Iowait = math.Round((curTimes.Iowait/total)*100*100) / 100
	curTimes.Irq = math.Round((curTimes.Irq/total)*100*100) / 100
	curTimes.Softirq = math.Round((curTimes.Softirq/total)*100*100) / 100
	curTimes.Steal = math.Round((curTimes.Steal/total)*100*100) / 100
	curTimes.Guest = math.Round((curTimes.Guest/total)*100*100) / 100
	curTimes.GuestNice = math.Round((curTimes.GuestNice/total)*100*100) / 100

	return curTimes, nil
}

type MyCustomMemStat struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Free        uint64  `json:"free"`
}

func getMEMStats() (MyCustomMemStat, error) {
	v, err := mem.VirtualMemory()

	if err != nil {
		return MyCustomMemStat{}, err
	}

	return MyCustomMemStat{
		Total:       v.Total,
		Used:        v.Used,
		Free:        v.Free,
		UsedPercent: v.UsedPercent,
		Available:   v.Available,
	}, nil
}

type ProcessInfo struct {
	PID         int32
	Name        string
	Username    string
	Memory      uint64
	CPUPercent  float64 // CPU usage percentage
	RunningTime string
}

func GetProcesses(n int) ([]ProcessInfo, error) {
	processes, err = process.Processes()

	if err != nil {
		return
	}

	var processInfos []ProcessInfo

	for _, p := range processes {
		pid := p.Pid

		name, err := p.Name()
		if err != nil {
			name = "Unknown"
		}

		createTime, err := p.CreateTime()
		if err != nil {
			createTime = 0
		}

		startingTime := time.Unix(createTime/1000, 0)
		runningTime := time.Since(startingTime).Truncate(time.Second)

		username, err := p.Username()
		if err != nil {
			username = "Unknown"
		}

		memoryInfo, err := p.MemoryInfo()
		if err != nil {
			processInfos = append(processInfos, ProcessInfo{
				PID:         pid,
				Name:        name,
				RunningTime: runningTime.String(),
				Username:    username,
				Memory:      0,
				CPUPercent:  0,
			})
			continue
		}

		memory := memoryInfo.RSS

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			cpuPercent = 0
		}

		processInfos = append(processInfos, ProcessInfo{
			PID:         pid,
			Name:        name,
			RunningTime: runningTime.String(),
			Username:    username,
			Memory:      memory,
			CPUPercent:  cpuPercent,
		})

	}

	sort.Slice(processInfos, func(i, j int) bool {
		return processInfos[i].CPUPercent > processInfos[j].CPUPercent
	})

	return processInfos, nil
}
