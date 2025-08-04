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
		fmt.Printf("OS: %s\nCPU: %s\nMemory: %s / %s (%s)\n",
			info.OS, info.CPU, info.MemoryUsed, info.MemoryTotal, info.MemoryPct)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
