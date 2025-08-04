package cmd

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/junler/sysinfo/internal/sysinfo"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start web server to show system info",
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/", servePage)
		fmt.Println("Serving on http://localhost:8080 ...")
		http.ListenAndServe(":8080", nil)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func servePage(w http.ResponseWriter, r *http.Request) {
	info, _ := sysinfo.GetSystemInfo()

	tmpl := `
		<h1>System Info</h1>
		<p><strong>OS:</strong> {{.OS}}</p>
		<p><strong>CPU:</strong> {{.CPU}}</p>
		<p><strong>Memory:</strong> {{.MemoryUsed}} / {{.MemoryTotal}} ({{.MemoryPct}})</p>
	`

	t := template.Must(template.New("web").Parse(tmpl))
	t.Execute(w, info)
}
