package cmd

import (
	"fmt"

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

		if len(info.Disk) > 0 {
			fmt.Println("\n=== Disk Information ===")
			for _, disk := range info.Disk {
				fmt.Printf("%s (%s) - %s\n", disk.Device, disk.Mountpoint, disk.Fstype)
				fmt.Printf("  Total: %.2f GB, Used: %.2f GB (%.1f%%), Free: %.2f GB\n",
					float64(disk.Total)/1e9, float64(disk.Used)/1e9, disk.UsedPercent, float64(disk.Free)/1e9)
			}
		}

		fmt.Println("\n=== Network Information ===")
		fmt.Printf("Bytes Sent: %.2f MB\n", float64(info.Network.BytesSent)/1e6)
		fmt.Printf("Bytes Received: %.2f MB\n", float64(info.Network.BytesRecv)/1e6)
		fmt.Printf("Network Interfaces: %d\n", len(info.Network.Interfaces))

		if len(info.LoadAverage) > 0 {
			fmt.Println("\n=== Load Average ===")
			fmt.Printf("1 min: %.2f, 5 min: %.2f, 15 min: %.2f\n",
				info.LoadAverage[0], info.LoadAverage[1], info.LoadAverage[2])
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
