package hardware

import (
	"fmt"
	"runtime"
)

// SystemInfo represents system hardware information
type SystemInfo struct {
	OS             string  `json:"os"`
	Architecture   string  `json:"architecture"`
	CPUCores       int     `json:"cpu_cores"`
	TotalRAMGB     float64 `json:"total_ram_gb"`
	AvailableRAMGB float64 `json:"available_ram_gb"`
	DiskSpaceGB    float64 `json:"disk_space_gb"`
	GPUName        string  `json:"gpu_name,omitempty"`
	GPUMemoryGB    float64 `json:"gpu_memory_gb,omitempty"`
}

// Validator provides hardware validation functionality
type Validator struct{}

// NewValidator creates a new hardware validator
func NewValidator() *Validator {
	return &Validator{}
}

// GetSystemInfo returns basic system information
// This is a simplified implementation for now
func (v *Validator) GetSystemInfo() (*SystemInfo, error) {
	return &SystemInfo{
		OS:             runtime.GOOS,
		Architecture:   runtime.GOARCH,
		CPUCores:       runtime.NumCPU(),
		TotalRAMGB:     16.0,  // Placeholder - would need platform-specific code
		AvailableRAMGB: 12.0,  // Placeholder - would need platform-specific code
		DiskSpaceGB:    50.0,  // Placeholder - would need platform-specific code
		GPUName:        "",    // Placeholder - would need GPU detection
		GPUMemoryGB:    0.0,   // Placeholder - would need GPU detection
	}, nil
}

// ValidateMinimumRequirements checks if system meets minimum requirements
func (v *Validator) ValidateMinimumRequirements() error {
	systemInfo, err := v.GetSystemInfo()
	if err != nil {
		return err
	}

	// Basic validation - in a real implementation, you'd check actual system resources
	if systemInfo.CPUCores < 2 {
		return fmt.Errorf("insufficient CPU cores: %d (minimum 2 required)", systemInfo.CPUCores)
	}

	return nil
}