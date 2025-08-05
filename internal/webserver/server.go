package webserver

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junler/sysinfo/internal/sysinfo"
)

//go:embed web/*
var WebFiles embed.FS

type WebServer struct {
	router *gin.Engine
	port   string
}

func NewWebServer(port string) *WebServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	ws := &WebServer{
		router: router,
		port:   port,
	}

	ws.setupRoutes()
	return ws
}

func (ws *WebServer) setupRoutes() {
	// Serve static files from embedded filesystem
	ws.router.StaticFS("/static", http.FS(WebFiles))

	// Serve the main page
	ws.router.GET("/", func(c *gin.Context) {
		data, err := WebFiles.ReadFile("web/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error loading page")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// API endpoints
	api := ws.router.Group("/api")
	{
		api.GET("/info", ws.getSystemInfo)
		api.GET("/ports", ws.getPorts)
		api.GET("/health", ws.healthCheck)
		api.GET("/processes", ws.getTopProcesses)
		api.GET("/monitoring", ws.getMonitoringData)
		api.GET("/temperature", ws.getTemperature)
		api.GET("/iostats", ws.getIOStats)
		api.GET("/users", ws.getUsers)
		api.GET("/services", ws.getServices)
	}
}

func (ws *WebServer) getSystemInfo(c *gin.Context) {
	info, err := sysinfo.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (ws *WebServer) getPorts(c *gin.Context) {
	ports, err := sysinfo.GetOpenPorts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ports)
}

func (ws *WebServer) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (ws *WebServer) getTopProcesses(c *gin.Context) {
	info, err := sysinfo.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"processes": info.TopProcesses})
}

func (ws *WebServer) getMonitoringData(c *gin.Context) {
	info, err := sysinfo.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a subset of monitoring data for dashboard
	monitoringData := gin.H{
		"load_average": info.LoadAverage,
		"memory":       info.Memory,
		"swap":         info.Swap,
		"cpu":          info.CPU,
		"io_stats":     info.IOStats,
		"network": gin.H{
			"bytes_sent":   info.Network.BytesSent,
			"bytes_recv":   info.Network.BytesRecv,
			"packets_sent": info.Network.PacketsSent,
			"packets_recv": info.Network.PacketsRecv,
			"errors_in":    info.Network.ErrorsIn,
			"errors_out":   info.Network.ErrorsOut,
		},
		"top_processes": info.TopProcesses[:min(10, len(info.TopProcesses))],
		"disk":          info.Disk,
	}

	c.JSON(http.StatusOK, monitoringData)
}

func (ws *WebServer) getTemperature(c *gin.Context) {
	info, err := sysinfo.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info.Temperature)
}

func (ws *WebServer) getIOStats(c *gin.Context) {
	info, err := sysinfo.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info.IOStats)
}

func (ws *WebServer) getUsers(c *gin.Context) {
	info, err := sysinfo.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": info.Users})
}

func (ws *WebServer) getServices(c *gin.Context) {
	info, err := sysinfo.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"services": info.SystemServices})
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (ws *WebServer) Start() error {
	return ws.router.Run(":" + ws.port)
}
