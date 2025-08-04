package sysinfo

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	OS          string
	CPU         string
	MemoryTotal string
	MemoryUsed  string
	MemoryPct   string
}

func GetSystemInfo() (*SystemInfo, error) {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Info()
	h, _ := host.Info()

	return &SystemInfo{
		OS:          fmt.Sprintf("%v %v (%v)", h.Platform, h.PlatformVersion, h.OS),
		CPU:         fmt.Sprintf("%v - %v cores", c[0].ModelName, c[0].Cores),
		MemoryTotal: fmt.Sprintf("%.2f GB", float64(v.Total)/1e9),
		MemoryUsed:  fmt.Sprintf("%.2f GB", float64(v.Used)/1e9),
		MemoryPct:   fmt.Sprintf("%.2f %%", v.UsedPercent),
	}, nil
}
