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

		fmt.Printf("%-10s %-10s %-10s\n", "PORT", "PROTOCOL", "STATUS")
		for _, p := range ports {
			fmt.Printf("%-10s %-10s %-10s\n", p.Port, p.Protocol, p.Process)
		}
	},
}

func init() {
	rootCmd.AddCommand(portsCmd)
}
