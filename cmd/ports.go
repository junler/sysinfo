package cmd

import (
	"fmt"

	"github.com/junler/sysinfo/internal/sysinfo"
	"github.com/spf13/cobra"
)

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Show open ports and their status",
	Run: func(cmd *cobra.Command, args []string) {
		ports, err := sysinfo.GetOpenPorts()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("=== Open Ports ===")
		fmt.Printf("%-8s %-10s %-10s %-20s %-8s %-15s\n", "PORT", "PROTOCOL", "STATUS", "PROCESS", "PID", "ADDRESS")
		fmt.Println("------------------------------------------------------------------------")
		for _, p := range ports {
			pidStr := "N/A"
			if p.PID != 0 {
				pidStr = fmt.Sprintf("%d", p.PID)
			}
			fmt.Printf("%-8s %-10s %-10s %-20s %-8s %-15s\n",
				p.Port, p.Protocol, p.Status, p.Process, pidStr, p.Address)
		}

		if len(ports) == 0 {
			fmt.Println("No open ports found.")
		} else {
			fmt.Printf("\nTotal: %d open ports\n", len(ports))
		}
	},
}

func init() {
	rootCmd.AddCommand(portsCmd)
}
