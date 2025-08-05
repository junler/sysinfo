package main

import (
	"fmt"
	"os"

	"github.com/junler/sysinfo/cmd"
)

func main() {
	// If no arguments, show usage
	if len(os.Args) == 1 {
		fmt.Println("System Information Tool")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  sysinfo info     - Show system information")
		fmt.Println("  sysinfo ports    - Show open ports")
		fmt.Println("  sysinfo serve    - Start web server")
		fmt.Println("  sysinfo version  - Show version information")
		fmt.Println("  sysinfo --help   - Show help")
		return
	}

	cmd.Execute()
}
