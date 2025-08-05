package sysinfo

import (
	"fmt"
	"sort"

	netutil "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type PortInfo struct {
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Status   string `json:"status"`
	Process  string `json:"process"`
	PID      int32  `json:"pid"`
	Address  string `json:"address"`
}

func GetOpenPorts() ([]PortInfo, error) {
	var result []PortInfo

	// Get TCP connections
	tcpConns, err := netutil.Connections("tcp")
	if err == nil {
		for _, conn := range tcpConns {
			if conn.Status == "LISTEN" {
				processName := "unknown"
				if conn.Pid != 0 {
					if proc, err := process.NewProcess(conn.Pid); err == nil {
						if name, err := proc.Name(); err == nil {
							processName = name
						}
					}
				}

				result = append(result, PortInfo{
					Port:     fmt.Sprintf("%d", conn.Laddr.Port),
					Protocol: "TCP",
					Status:   conn.Status,
					Process:  processName,
					PID:      conn.Pid,
					Address:  conn.Laddr.IP,
				})
			}
		}
	}

	// Get UDP connections
	udpConns, err := netutil.Connections("udp")
	if err == nil {
		for _, conn := range udpConns {
			processName := "unknown"
			if conn.Pid != 0 {
				if proc, err := process.NewProcess(conn.Pid); err == nil {
					if name, err := proc.Name(); err == nil {
						processName = name
					}
				}
			}

			result = append(result, PortInfo{
				Port:     fmt.Sprintf("%d", conn.Laddr.Port),
				Protocol: "UDP",
				Status:   "ACTIVE",
				Process:  processName,
				PID:      conn.Pid,
				Address:  conn.Laddr.IP,
			})
		}
	}

	// Sort by port number
	sort.Slice(result, func(i, j int) bool {
		return result[i].Port < result[j].Port
	})

	return result, nil
}
