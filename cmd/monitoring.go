package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/junler/sysinfo/internal/sysinfo"
	"github.com/spf13/cobra"
)

var monitoringCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Show detailed system monitoring metrics",
	Long: `Display comprehensive system monitoring information including:
- Top processes by CPU and memory usage
- Detailed network statistics
- I/O performance metrics
- System temperature (if available)
- Load averages
- Logged in users
- System services status`,
	Run: func(cmd *cobra.Command, args []string) {
		info, err := sysinfo.GetSystemInfo()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// System Overview
		fmt.Println("=== SYSTEM MONITORING DASHBOARD ===")
		fmt.Printf("Host: %s | OS: %s | Uptime: %s\n", info.Hostname, info.OS, info.Uptime)
		fmt.Printf("Architecture: %s | Kernel: %s\n", info.Architecture, info.KernelVersion)
		fmt.Printf("Processes: %d | Last Boot: %s\n", info.ProcessCount, info.LastBoot)

		// Load Average
		fmt.Println("\n=== LOAD AVERAGE ===")
		fmt.Printf("1min: %6.2f | 5min: %6.2f | 15min: %6.2f\n",
			info.LoadAverage.Load1, info.LoadAverage.Load5, info.LoadAverage.Load15)

		// CPU Details
		fmt.Println("\n=== CPU METRICS ===")
		fmt.Printf("Model: %s\n", info.CPU.ModelName)
		fmt.Printf("Cores: %d Physical / %d Logical | Frequency: %.0f MHz\n",
			info.CPU.Cores, info.CPU.LogicalCores, info.CPU.Frequency)
		if len(info.CPU.Usage) > 0 {
			fmt.Printf("Per-Core Usage: ")
			for i, usage := range info.CPU.Usage {
				if i > 0 && i%8 == 0 {
					fmt.Printf("\n                ")
				}
				fmt.Printf("%5.1f%% ", usage)
			}
			fmt.Println()
		}
		if info.CPU.Temperature > 0 {
			fmt.Printf("Temperature: %.1f°C\n", info.CPU.Temperature)
		}

		// Memory and Swap
		fmt.Println("\n=== MEMORY METRICS ===")
		fmt.Printf("RAM:  Total: %7.1f GB | Used: %7.1f GB (%5.1f%%) | Available: %7.1f GB\n",
			float64(info.Memory.Total)/1e9, float64(info.Memory.Used)/1e9,
			info.Memory.UsedPercent, float64(info.Memory.Available)/1e9)
		fmt.Printf("      Cached: %6.1f GB | Buffers: %5.1f GB | Shared: %7.1f GB\n",
			float64(info.Memory.Cached)/1e9, float64(info.Memory.Buffers)/1e9, float64(info.Memory.Shared)/1e9)

		if info.Swap.Total > 0 {
			fmt.Printf("Swap: Total: %7.1f GB | Used: %7.1f GB (%5.1f%%) | Free: %9.1f GB\n",
				float64(info.Swap.Total)/1e9, float64(info.Swap.Used)/1e9,
				info.Swap.UsedPercent, float64(info.Swap.Free)/1e9)
		} else {
			fmt.Println("Swap: Not configured")
		}

		// I/O Statistics
		fmt.Println("\n=== DISK I/O METRICS ===")
		fmt.Printf("Read:  %8.1f MB (%8d ops) | Avg Time: %6.1f ms\n",
			float64(info.IOStats.DiskReadBytes)/1e6, info.IOStats.DiskReadCount,
			getAvgTime(info.IOStats.DiskReadTime, info.IOStats.DiskReadCount))
		fmt.Printf("Write: %8.1f MB (%8d ops) | Avg Time: %6.1f ms\n",
			float64(info.IOStats.DiskWriteBytes)/1e6, info.IOStats.DiskWriteCount,
			getAvgTime(info.IOStats.DiskWriteTime, info.IOStats.DiskWriteCount))

		// Network Statistics
		fmt.Println("\n=== NETWORK METRICS ===")
		fmt.Printf("Traffic:  Sent: %8.1f MB | Received: %8.1f MB\n",
			float64(info.Network.BytesSent)/1e6, float64(info.Network.BytesRecv)/1e6)
		fmt.Printf("Packets:  Sent: %8d    | Received: %8d\n",
			info.Network.PacketsSent, info.Network.PacketsRecv)
		fmt.Printf("Errors:   In: %10d    | Out: %12d\n",
			info.Network.ErrorsIn, info.Network.ErrorsOut)
		fmt.Printf("Drops:    In: %10d    | Out: %12d\n",
			info.Network.DropsIn, info.Network.DropsOut)

		// Top Processes
		if len(info.TopProcesses) > 0 {
			fmt.Println("\n=== TOP PROCESSES BY CPU ===")
			fmt.Printf("%-8s %-20s %-10s %8s %8s %10s %-15s\n",
				"PID", "NAME", "STATUS", "CPU%", "MEM%", "RSS(MB)", "USER")
			fmt.Println(strings.Repeat("-", 90))
			for i, proc := range info.TopProcesses {
				if i >= 15 { // Show top 15
					break
				}
				fmt.Printf("%-8d %-20s %-10s %7.1f%% %7.1f%% %9.1f %-15s\n",
					proc.PID,
					truncateString(proc.Name, 20),
					proc.Status,
					proc.CPUPercent,
					proc.MemPercent,
					float64(proc.MemoryRSS)/1024/1024,
					truncateString(proc.Username, 15))
			}
		}

		// Disk Usage
		if len(info.Disk) > 0 {
			fmt.Println("\n=== DISK USAGE ===")
			fmt.Printf("%-20s %-15s %-10s %10s %10s %10s %8s\n",
				"DEVICE", "MOUNTPOINT", "FSTYPE", "TOTAL", "USED", "FREE", "USE%")
			fmt.Println(strings.Repeat("-", 95))
			for _, disk := range info.Disk {
				fmt.Printf("%-20s %-15s %-10s %9.1fG %9.1fG %9.1fG %7.1f%%\n",
					truncateString(disk.Device, 20),
					truncateString(disk.Mountpoint, 15),
					disk.Fstype,
					float64(disk.Total)/1e9,
					float64(disk.Used)/1e9,
					float64(disk.Free)/1e9,
					disk.UsedPercent)
			}
		}

		// Active Users
		if len(info.Users) > 0 {
			fmt.Println("\n=== ACTIVE USERS ===")
			fmt.Printf("%-15s %-10s %-15s %-20s\n", "USER", "TTY", "HOST", "LOGIN TIME")
			fmt.Println(strings.Repeat("-", 65))
			for _, user := range info.Users {
				fmt.Printf("%-15s %-10s %-15s %-20s\n",
					truncateString(user.User, 15),
					truncateString(user.Terminal, 10),
					truncateString(user.Host, 15),
					time.Unix(user.Started, 0).Format("2006-01-02 15:04"))
			}
		}

		// System Services (sample)
		if len(info.SystemServices) > 0 {
			fmt.Println("\n=== SYSTEM SERVICES (Sample) ===")
			fmt.Printf("%-25s %-12s %8s\n", "SERVICE", "STATUS", "PID")
			fmt.Println(strings.Repeat("-", 50))
			serviceCount := 0
			for _, service := range info.SystemServices {
				if serviceCount >= 15 { // Limit to 15 services
					break
				}
				fmt.Printf("%-25s %-12s %8d\n",
					truncateString(service.Name, 25),
					service.Status,
					service.PID)
				serviceCount++
			}
		}

		// Network Interfaces
		if len(info.Network.Interfaces) > 0 {
			fmt.Println("\n=== NETWORK INTERFACES ===")
			for _, iface := range info.Network.Interfaces {
				fmt.Printf("Interface: %s (MTU: %d)\n", iface.Name, iface.MTU)
				if len(iface.Addresses) > 0 {
					fmt.Printf("  Addresses: %s\n", strings.Join(iface.Addresses, ", "))
				}
				if iface.HardwareAddr != "" {
					fmt.Printf("  Hardware: %s\n", iface.HardwareAddr)
				}
				if len(iface.Flags) > 0 {
					fmt.Printf("  Flags: %s\n", strings.Join(iface.Flags, ", "))
				}
				fmt.Println()
			}
		}
	},
}

// getAvgTime calculates average time per operation
func getAvgTime(totalTime, count uint64) float64 {
	if count == 0 {
		return 0.0
	}
	return float64(totalTime) / float64(count)
}

func init() {
	rootCmd.AddCommand(monitoringCmd)
}
