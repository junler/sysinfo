package sysinfo

import (
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type SystemInfo struct {
	OS             string          `json:"os"`
	Hostname       string          `json:"hostname"`
	Uptime         string          `json:"uptime"`
	CPU            CPUInfo         `json:"cpu"`
	Memory         MemoryInfo      `json:"memory"`
	Swap           SwapInfo        `json:"swap"`
	Disk           []DiskInfo      `json:"disk"`
	Network        NetworkInfo     `json:"network"`
	LoadAverage    LoadAverageInfo `json:"load_average"`
	ProcessCount   uint64          `json:"process_count"`
	Architecture   string          `json:"architecture"`
	KernelVersion  string          `json:"kernel_version"`
	LastBoot       string          `json:"last_boot"`
	TopProcesses   []ProcessInfo   `json:"top_processes"`
	Temperature    TemperatureInfo `json:"temperature"`
	IOStats        IOStatsInfo     `json:"io_stats"`
	Users          []UserInfo      `json:"users"`
	SystemServices []ServiceInfo   `json:"system_services"`
}

type CPUInfo struct {
	ModelName    string    `json:"model_name"`
	Cores        int32     `json:"cores"`
	LogicalCores int32     `json:"logical_cores"`
	Usage        []float64 `json:"usage"`
	Frequency    float64   `json:"frequency"`
	CacheSize    int32     `json:"cache_size"`
	Temperature  float64   `json:"temperature"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Free        uint64  `json:"free"`
	Cached      uint64  `json:"cached"`
	Buffers     uint64  `json:"buffers"`
	Shared      uint64  `json:"shared"`
	Active      uint64  `json:"active"`
	Inactive    uint64  `json:"inactive"`
}

type SwapInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskInfo struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
	InodesTotal uint64  `json:"inodes_total"`
	InodesUsed  uint64  `json:"inodes_used"`
	InodesFree  uint64  `json:"inodes_free"`
}

type NetworkInfo struct {
	Interfaces  []NetworkInterface `json:"interfaces"`
	BytesSent   uint64             `json:"bytes_sent"`
	BytesRecv   uint64             `json:"bytes_recv"`
	PacketsSent uint64             `json:"packets_sent"`
	PacketsRecv uint64             `json:"packets_recv"`
	ErrorsIn    uint64             `json:"errors_in"`
	ErrorsOut   uint64             `json:"errors_out"`
	DropsIn     uint64             `json:"drops_in"`
	DropsOut    uint64             `json:"drops_out"`
}

type NetworkInterface struct {
	Name         string   `json:"name"`
	Addresses    []string `json:"addresses"`
	MTU          int      `json:"mtu"`
	Flags        []string `json:"flags"`
	HardwareAddr string   `json:"hardware_addr"`
}

type LoadAverageInfo struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type ProcessInfo struct {
	PID         int32   `json:"pid"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemPercent  float32 `json:"mem_percent"`
	MemoryRSS   uint64  `json:"memory_rss"`
	MemoryVMS   uint64  `json:"memory_vms"`
	CreateTime  int64   `json:"create_time"`
	NumThreads  int32   `json:"num_threads"`
	Username    string  `json:"username"`
	CommandLine string  `json:"command_line"`
}

type TemperatureInfo struct {
	CPUTemp     float64            `json:"cpu_temp"`
	Sensors     []SensorInfo       `json:"sensors"`
	ThermalZone map[string]float64 `json:"thermal_zone"`
}

type SensorInfo struct {
	Name        string  `json:"name"`
	Temperature float64 `json:"temperature"`
	High        float64 `json:"high"`
	Critical    float64 `json:"critical"`
}

type IOStatsInfo struct {
	DiskReadBytes  uint64 `json:"disk_read_bytes"`
	DiskWriteBytes uint64 `json:"disk_write_bytes"`
	DiskReadCount  uint64 `json:"disk_read_count"`
	DiskWriteCount uint64 `json:"disk_write_count"`
	DiskReadTime   uint64 `json:"disk_read_time"`
	DiskWriteTime  uint64 `json:"disk_write_time"`
}

type UserInfo struct {
	User     string `json:"user"`
	Terminal string `json:"terminal"`
	Host     string `json:"host"`
	Started  int64  `json:"started"`
}

type ServiceInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	PID    int32  `json:"pid"`
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

	// Swap Info
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		swapInfo = &mem.SwapMemoryStat{}
	}

	// Load Average
	loadAvg, err := load.Avg()
	var loadAvgInfo LoadAverageInfo
	if err == nil {
		loadAvgInfo = LoadAverageInfo{
			Load1:  loadAvg.Load1,
			Load5:  loadAvg.Load5,
			Load15: loadAvg.Load15,
		}
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
			InodesTotal: usage.InodesTotal,
			InodesUsed:  usage.InodesUsed,
			InodesFree:  usage.InodesFree,
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
			Name:         iface.Name,
			Addresses:    addresses,
			MTU:          iface.MTU,
			Flags:        iface.Flags,
			HardwareAddr: iface.HardwareAddr,
		})
	}

	netIO, err := net.IOCounters(false)
	var bytesSent, bytesRecv, packetsSent, packetsRecv, errorsIn, errorsOut, dropsIn, dropsOut uint64
	if err == nil && len(netIO) > 0 {
		bytesSent = netIO[0].BytesSent
		bytesRecv = netIO[0].BytesRecv
		packetsSent = netIO[0].PacketsSent
		packetsRecv = netIO[0].PacketsRecv
		errorsIn = netIO[0].Errin
		errorsOut = netIO[0].Errout
		dropsIn = netIO[0].Dropin
		dropsOut = netIO[0].Dropout
	}

	// Get top processes
	topProcesses := getTopProcesses(10)

	// Get temperature info
	tempInfo := getTemperatureInfo()

	// Get IO stats
	ioStats := getIOStats()

	// Get users
	users := getLoggedInUsers()

	// Get system services (basic implementation)
	services := getSystemServices()

	var cpuInfo CPUInfo
	if len(cpuInfos) > 0 {
		cpuInfo = CPUInfo{
			ModelName:    cpuInfos[0].ModelName,
			Cores:        cpuInfos[0].Cores,
			LogicalCores: int32(len(cpuInfos)),
			Usage:        cpuUsage,
			Frequency:    cpuInfos[0].Mhz,
			CacheSize:    cpuInfos[0].CacheSize,
			Temperature:  getCPUTemperature(),
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
			Shared:      memInfo.Shared,
			Active:      memInfo.Active,
			Inactive:    memInfo.Inactive,
		},
		Swap: SwapInfo{
			Total:       swapInfo.Total,
			Used:        swapInfo.Used,
			Free:        swapInfo.Free,
			UsedPercent: swapInfo.UsedPercent,
		},
		Disk: diskInfos,
		Network: NetworkInfo{
			Interfaces:  networkInterfaces,
			BytesSent:   bytesSent,
			BytesRecv:   bytesRecv,
			PacketsSent: packetsSent,
			PacketsRecv: packetsRecv,
			ErrorsIn:    errorsIn,
			ErrorsOut:   errorsOut,
			DropsIn:     dropsIn,
			DropsOut:    dropsOut,
		},
		LoadAverage:    loadAvgInfo,
		ProcessCount:   hostInfo.Procs,
		Architecture:   hostInfo.KernelArch,
		KernelVersion:  hostInfo.KernelVersion,
		LastBoot:       time.Unix(int64(hostInfo.BootTime), 0).Format("2006-01-02 15:04:05"),
		TopProcesses:   topProcesses,
		Temperature:    tempInfo,
		IOStats:        ioStats,
		Users:          users,
		SystemServices: services,
	}, nil
}

// getTopProcesses returns the top N processes by CPU usage
func getTopProcesses(limit int) []ProcessInfo {
	processes, err := process.Processes()
	if err != nil {
		return []ProcessInfo{}
	}

	var procInfos []ProcessInfo
	for _, p := range processes {
		// Get process info
		name, _ := p.Name()
		status, _ := p.Status()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		memInfo, _ := p.MemoryInfo()
		createTime, _ := p.CreateTime()
		numThreads, _ := p.NumThreads()
		username, _ := p.Username()
		cmdline, _ := p.Cmdline()

		var memRSS, memVMS uint64
		if memInfo != nil {
			memRSS = memInfo.RSS
			memVMS = memInfo.VMS
		}

		// Convert status slice to string
		statusStr := ""
		if len(status) > 0 {
			statusStr = status[0]
		}

		procInfo := ProcessInfo{
			PID:         p.Pid,
			Name:        name,
			Status:      statusStr,
			CPUPercent:  cpuPercent,
			MemPercent:  memPercent,
			MemoryRSS:   memRSS,
			MemoryVMS:   memVMS,
			CreateTime:  createTime,
			NumThreads:  numThreads,
			Username:    username,
			CommandLine: cmdline,
		}
		procInfos = append(procInfos, procInfo)
	}

	// Sort by CPU usage
	sort.Slice(procInfos, func(i, j int) bool {
		return procInfos[i].CPUPercent > procInfos[j].CPUPercent
	})

	// Return top N processes
	if len(procInfos) > limit {
		return procInfos[:limit]
	}
	return procInfos
}

// getTemperatureInfo collects temperature information from various sensors
func getTemperatureInfo() TemperatureInfo {
	tempInfo := TemperatureInfo{
		CPUTemp:     getCPUTemperature(),
		Sensors:     []SensorInfo{},
		ThermalZone: make(map[string]float64),
	}

	// Try to get sensor temperatures (this may not work on all systems)
	// gopsutil doesn't have direct temperature support, so we'll return empty for now
	// In a real implementation, you might read from /sys/class/thermal/ on Linux
	// or use platform-specific APIs

	return tempInfo
}

// getCPUTemperature attempts to get CPU temperature
func getCPUTemperature() float64 {
	// This is a placeholder - actual implementation would depend on the platform
	// On Linux, you might read from /sys/class/thermal/thermal_zone*/temp
	// On macOS, you might use IOKit
	// For now, return 0 to indicate unavailable
	return 0.0
}

// getIOStats collects disk I/O statistics
func getIOStats() IOStatsInfo {
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return IOStatsInfo{}
	}

	var totalReadBytes, totalWriteBytes, totalReadCount, totalWriteCount, totalReadTime, totalWriteTime uint64

	for _, counter := range ioCounters {
		totalReadBytes += counter.ReadBytes
		totalWriteBytes += counter.WriteBytes
		totalReadCount += counter.ReadCount
		totalWriteCount += counter.WriteCount
		totalReadTime += counter.ReadTime
		totalWriteTime += counter.WriteTime
	}

	return IOStatsInfo{
		DiskReadBytes:  totalReadBytes,
		DiskWriteBytes: totalWriteBytes,
		DiskReadCount:  totalReadCount,
		DiskWriteCount: totalWriteCount,
		DiskReadTime:   totalReadTime,
		DiskWriteTime:  totalWriteTime,
	}
}

// getLoggedInUsers returns information about currently logged in users
func getLoggedInUsers() []UserInfo {
	users, err := host.Users()
	if err != nil {
		return []UserInfo{}
	}

	var userInfos []UserInfo
	for _, user := range users {
		userInfos = append(userInfos, UserInfo{
			User:     user.User,
			Terminal: user.Terminal,
			Host:     user.Host,
			Started:  int64(user.Started),
		})
	}

	return userInfos
}

// getSystemServices returns basic information about system services
func getSystemServices() []ServiceInfo {
	// This is a basic implementation that shows some common processes
	// A more complete implementation would interact with systemd on Linux,
	// launchd on macOS, or Windows services on Windows

	services := []ServiceInfo{}

	// Get all processes and filter for common system services
	processes, err := process.Processes()
	if err != nil {
		return services
	}

	systemProcessNames := map[string]bool{
		"systemd":        true,
		"kernel":         true,
		"kthreadd":       true,
		"ksoftirqd":      true,
		"rcu_":           true,
		"watchdog":       true,
		"sshd":           true,
		"NetworkManager": true,
		"dbus":           true,
		"cron":           true,
		"rsyslog":        true,
		"apache2":        true,
		"nginx":          true,
		"mysql":          true,
		"postgres":       true,
		"docker":         true,
		"containerd":     true,
	}

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		// Check if this is a system service
		isSystemService := false
		for serviceName := range systemProcessNames {
			if name == serviceName || (len(name) > len(serviceName) && name[:len(serviceName)] == serviceName) {
				isSystemService = true
				break
			}
		}

		if isSystemService {
			status, _ := p.Status()
			statusStr := ""
			if len(status) > 0 {
				statusStr = status[0]
			}
			services = append(services, ServiceInfo{
				Name:   name,
				Status: statusStr,
				PID:    p.Pid,
			})
		}
	}

	return services
}
