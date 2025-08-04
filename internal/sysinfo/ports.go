package sysinfo

import (
	"fmt"
	"sort"

	netutil "github.com/shirou/gopsutil/v3/net"
)

type PortInfo struct {
	Port     string
	Protocol string
	Process  string
}

func GetOpenPorts() ([]PortInfo, error) {
	var result []PortInfo

	conns, err := netutil.Connections("tcp")
	if err != nil {
		return nil, err
	}

	for _, conn := range conns {
		if conn.Status == "LISTEN" {
			result = append(result, PortInfo{
				Port:     fmt.Sprintf("%d", conn.Laddr.Port),
				Protocol: "tcp",
				Process:  conn.Status,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Port < result[j].Port
	})

	return result, nil
}
