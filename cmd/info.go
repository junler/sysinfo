package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/junler/sysinfo/internal/sysinfo"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show system information",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := sysinfo.GetSystemInfo()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("=== System Information ===")
		fmt.Printf("OS: %s\n", info.OS)
		fmt.Printf("Hostname: %s\n", info.Hostname)
		fmt.Printf("Architecture: %s\n", info.Architecture)
		fmt.Printf("Kernel Version: %s\n", info.KernelVersion)
		fmt.Printf("Uptime: %s\n", info.Uptime)
		fmt.Printf("Last Boot: %s\n", info.LastBoot)
		fmt.Printf("Process Count: %d\n", info.ProcessCount)

		fmt.Println("\n=== CPU Information ===")
		fmt.Printf("Model: %s\n", info.CPU.ModelName)
		fmt.Printf("Physical Cores: %d\n", info.CPU.Cores)
		fmt.Printf("Logical Cores: %d\n", info.CPU.LogicalCores)
		fmt.Printf("Frequency: %.2f MHz\n", info.CPU.Frequency)
		if len(info.CPU.Usage) > 0 {
			fmt.Printf("CPU Usage per core: ")
			for i, usage := range info.CPU.Usage {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%.1f%%", usage)
			}
			fmt.Println()
		}

		fmt.Println("\n=== Memory Information ===")
		fmt.Printf("Total: %.2f GB\n", float64(info.Memory.Total)/1e9)
		fmt.Printf("Used: %.2f GB (%.1f%%)\n", float64(info.Memory.Used)/1e9, info.Memory.UsedPercent)
		fmt.Printf("Free: %.2f GB\n", float64(info.Memory.Free)/1e9)
		fmt.Printf("Available: %.2f GB\n", float64(info.Memory.Available)/1e9)
		fmt.Printf("Cached: %.2f GB\n", float64(info.Memory.Cached)/1e9)
		fmt.Printf("Buffers: %.2f GB\n", float64(info.Memory.Buffers)/1e9)
		fmt.Printf("Shared: %.2f GB\n", float64(info.Memory.Shared)/1e9)
		fmt.Printf("Active: %.2f GB\n", float64(info.Memory.Active)/1e9)
		fmt.Printf("Inactive: %.2f GB\n", float64(info.Memory.Inactive)/1e9)

		fmt.Println("\n=== Swap Information ===")
		if info.Swap.Total > 0 {
			fmt.Printf("Total: %.2f GB\n", float64(info.Swap.Total)/1e9)
			fmt.Printf("Used: %.2f GB (%.1f%%)\n", float64(info.Swap.Used)/1e9, info.Swap.UsedPercent)
			fmt.Printf("Free: %.2f GB\n", float64(info.Swap.Free)/1e9)
		} else {
			fmt.Println("No swap configured")
		}

		if len(info.Disk) > 0 {
			fmt.Println("\n=== Disk Information ===")
			for _, disk := range info.Disk {
				fmt.Printf("%s (%s) - %s\n", disk.Device, disk.Mountpoint, disk.Fstype)
				fmt.Printf("  Total: %.2f GB, Used: %.2f GB (%.1f%%), Free: %.2f GB\n",
					float64(disk.Total)/1e9, float64(disk.Used)/1e9, disk.UsedPercent, float64(disk.Free)/1e9)
				if disk.InodesTotal > 0 {
					fmt.Printf("  Inodes: Total: %d, Used: %d, Free: %d\n",
						disk.InodesTotal, disk.InodesUsed, disk.InodesFree)
				}
			}
		}

		fmt.Println("\n=== Network Information ===")
		fmt.Printf("Bytes Sent: %.2f MB\n", float64(info.Network.BytesSent)/1e6)
		fmt.Printf("Bytes Received: %.2f MB\n", float64(info.Network.BytesRecv)/1e6)
		fmt.Printf("Packets Sent: %d\n", info.Network.PacketsSent)
		fmt.Printf("Packets Received: %d\n", info.Network.PacketsRecv)
		fmt.Printf("Errors In: %d, Errors Out: %d\n", info.Network.ErrorsIn, info.Network.ErrorsOut)
		fmt.Printf("Drops In: %d, Drops Out: %d\n", info.Network.DropsIn, info.Network.DropsOut)
		fmt.Printf("Network Interfaces: %d\n", len(info.Network.Interfaces))

		fmt.Println("\n=== Load Average ===")
		fmt.Printf("1 min: %.2f, 5 min: %.2f, 15 min: %.2f\n",
			info.LoadAverage.Load1, info.LoadAverage.Load5, info.LoadAverage.Load15)

		fmt.Println("\n=== I/O Statistics ===")
		fmt.Printf("Disk Read: %.2f MB (%d operations, %d ms)\n",
			float64(info.IOStats.DiskReadBytes)/1e6, info.IOStats.DiskReadCount, info.IOStats.DiskReadTime)
		fmt.Printf("Disk Write: %.2f MB (%d operations, %d ms)\n",
			float64(info.IOStats.DiskWriteBytes)/1e6, info.IOStats.DiskWriteCount, info.IOStats.DiskWriteTime)

		if len(info.TopProcesses) > 0 {
			fmt.Println("\n=== Top Processes ===")
			fmt.Printf("%-8s %-20s %-10s %-8s %-8s %-12s\n", "PID", "Name", "Status", "CPU%", "MEM%", "User")
			fmt.Println(strings.Repeat("-", 80))
			for i, proc := range info.TopProcesses {
				if i >= 10 { // Limit to top 10 for display
					break
				}
				fmt.Printf("%-8d %-20s %-10s %-8.1f %-8.1f %-12s\n",
					proc.PID, truncateString(proc.Name, 20), proc.Status,
					proc.CPUPercent, proc.MemPercent, truncateString(proc.Username, 12))
			}
		}

		if len(info.Users) > 0 {
			fmt.Println("\n=== Logged In Users ===")
			for _, user := range info.Users {
				fmt.Printf("User: %s, Terminal: %s, Host: %s, Started: %s\n",
					user.User, user.Terminal, user.Host,
					time.Unix(user.Started, 0).Format("2006-01-02 15:04:05"))
			}
		}

		if len(info.SystemServices) > 0 && len(info.SystemServices) <= 20 {
			fmt.Println("\n=== System Services (Sample) ===")
			for i, service := range info.SystemServices {
				if i >= 10 { // Limit display
					break
				}
				fmt.Printf("%-20s %-10s PID: %d\n", truncateString(service.Name, 20), service.Status, service.PID)
			}
		}

		if info.Temperature.CPUTemp > 0 {
			fmt.Println("\n=== Temperature Information ===")
			fmt.Printf("CPU Temperature: %.1f°C\n", info.Temperature.CPUTemp)
		}
	},
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
