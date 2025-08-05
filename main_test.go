package main

import (
	"testing"
	"time"

	"github.com/junler/sysinfo/internal/sysinfo"
)

func TestGetSystemInfo(t *testing.T) {
	start := time.Now()
	info, err := sysinfo.GetSystemInfo()
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("GetSystemInfo failed: %v", err)
	}

	if info == nil {
		t.Fatal("GetSystemInfo returned nil")
	}

	// Check that basic fields are populated
	if info.OS == "" {
		t.Error("OS field is empty")
	}

	if info.Hostname == "" {
		t.Error("Hostname field is empty")
	}

	if info.CPU.ModelName == "" {
		t.Error("CPU ModelName is empty")
	}

	if info.Memory.Total == 0 {
		t.Error("Memory Total is zero")
	}

	t.Logf("GetSystemInfo completed in %v", duration)

	// Performance check - should complete within 5 seconds
	if duration > 5*time.Second {
		t.Errorf("GetSystemInfo took too long: %v", duration)
	}
}

func TestGetOpenPorts(t *testing.T) {
	start := time.Now()
	ports, err := sysinfo.GetOpenPorts()
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("GetOpenPorts failed: %v", err)
	}

	if ports == nil {
		t.Fatal("GetOpenPorts returned nil")
	}

	t.Logf("GetOpenPorts found %d ports in %v", len(ports), duration)

	// Performance check - should complete within 10 seconds
	if duration > 10*time.Second {
		t.Errorf("GetOpenPorts took too long: %v", duration)
	}

	// Check port data structure
	for i, port := range ports {
		if port.Port == "" {
			t.Errorf("Port %d has empty Port field", i)
		}
		if port.Protocol == "" {
			t.Errorf("Port %d has empty Protocol field", i)
		}
		if i >= 5 { // Only check first 5 ports
			break
		}
	}
}

func BenchmarkGetSystemInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := sysinfo.GetSystemInfo()
		if err != nil {
			b.Fatalf("GetSystemInfo failed: %v", err)
		}
	}
}

func BenchmarkGetOpenPorts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := sysinfo.GetOpenPorts()
		if err != nil {
			b.Fatalf("GetOpenPorts failed: %v", err)
		}
	}
}
