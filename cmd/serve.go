package cmd

import (
	"fmt"
	"log"

	"github.com/junler/sysinfo/internal/webserver"
	"github.com/spf13/cobra"
)

var port string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Long:  "Start the web server to display system information in a web interface",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting web server on port %s...\n", port)
		fmt.Printf("Open http://localhost:%s in your browser\n", port)

		server := webserver.NewWebServer(port)
		if err := server.Start(); err != nil {
			log.Fatal("Failed to start web server:", err)
		}
	},
}

func init() {
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the web server on")
	rootCmd.AddCommand(serveCmd)
}
