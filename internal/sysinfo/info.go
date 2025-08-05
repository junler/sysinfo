package sysinfo

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemInfo struct {
	OS            string      `json:"os"`
	Hostname      string      `json:"hostname"`
	Uptime        string      `json:"uptime"`
	CPU           CPUInfo     `json:"cpu"`
	Memory        MemoryInfo  `json:"memory"`
	Disk          []DiskInfo  `json:"disk"`
	Network       NetworkInfo `json:"network"`
	LoadAverage   []float64   `json:"load_average"`
	ProcessCount  uint64      `json:"process_count"`
	Architecture  string      `json:"architecture"`
	KernelVersion string      `json:"kernel_version"`
	LastBoot      string      `json:"last_boot"`
}

type CPUInfo struct {
	ModelName    string    `json:"model_name"`
	Cores        int32     `json:"cores"`
	LogicalCores int32     `json:"logical_cores"`
	Usage        []float64 `json:"usage"`
	Frequency    float64   `json:"frequency"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Free        uint64  `json:"free"`
	Cached      uint64  `json:"cached"`
	Buffers     uint64  `json:"buffers"`
}

type DiskInfo struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

type NetworkInfo struct {
	Interfaces []NetworkInterface `json:"interfaces"`
	BytesSent  uint64             `json:"bytes_sent"`
	BytesRecv  uint64             `json:"bytes_recv"`
}

type NetworkInterface struct {
	Name      string   `json:"name"`
	Addresses []string `json:"addresses"`
	MTU       int      `json:"mtu"`
	Flags     []string `json:"flags"`
}

func GetSystemInfo() (*SystemInfo, error) {
	// Host Info
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	// CPU Info
	cpuInfos, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	cpuUsage, err := cpu.Percent(time.Second, true)
	if err != nil {
		cpuUsage = []float64{}
	}

	// Memory Info
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// Disk Info
	diskPartitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var diskInfos []DiskInfo
	for _, partition := range diskPartitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		diskInfos = append(diskInfos, DiskInfo{
			Device:      partition.Device,
			Mountpoint:  partition.Mountpoint,
			Fstype:      partition.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	// Network Info
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var networkInterfaces []NetworkInterface
	for _, iface := range netInterfaces {
		var addresses []string
		for _, addr := range iface.Addrs {
			addresses = append(addresses, addr.Addr)
		}
		networkInterfaces = append(networkInterfaces, NetworkInterface{
			Name:      iface.Name,
			Addresses: addresses,
			MTU:       iface.MTU,
			Flags:     iface.Flags,
		})
	}

	netIO, err := net.IOCounters(false)
	var bytesSent, bytesRecv uint64
	if err == nil && len(netIO) > 0 {
		bytesSent = netIO[0].BytesSent
		bytesRecv = netIO[0].BytesRecv
	}

	// Load Average (may not be available on all systems)
	var loadAvgSlice []float64
	// LoadAverage is not available on all platforms, so we'll skip it for now
	loadAvgSlice = []float64{0.0, 0.0, 0.0}

	var cpuInfo CPUInfo
	if len(cpuInfos) > 0 {
		cpuInfo = CPUInfo{
			ModelName:    cpuInfos[0].ModelName,
			Cores:        cpuInfos[0].Cores,
			LogicalCores: int32(len(cpuInfos)),
			Usage:        cpuUsage,
			Frequency:    cpuInfos[0].Mhz,
		}
	}

	return &SystemInfo{
		OS:       fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion),
		Hostname: hostInfo.Hostname,
		Uptime:   fmt.Sprintf("%.0f hours", float64(hostInfo.Uptime)/3600),
		CPU:      cpuInfo,
		Memory: MemoryInfo{
			Total:       memInfo.Total,
			Available:   memInfo.Available,
			Used:        memInfo.Used,
			UsedPercent: memInfo.UsedPercent,
			Free:        memInfo.Free,
			Cached:      memInfo.Cached,
			Buffers:     memInfo.Buffers,
		},
		Disk: diskInfos,
		Network: NetworkInfo{
			Interfaces: networkInterfaces,
			BytesSent:  bytesSent,
			BytesRecv:  bytesRecv,
		},
		LoadAverage:   loadAvgSlice,
		ProcessCount:  hostInfo.Procs,
		Architecture:  hostInfo.KernelArch,
		KernelVersion: hostInfo.KernelVersion,
		LastBoot:      time.Unix(int64(hostInfo.BootTime), 0).Format("2006-01-02 15:04:05"),
	}, nil
}
