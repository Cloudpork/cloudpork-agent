package models

import (
	"fmt"
	"os/exec"
	"strings"
)

// ModelInfo represents information about a model
type ModelInfo struct {
	Name  string `json:"name"`
	Size  string `json:"size"`
	Speed string `json:"speed"`
}

// ModelStatus represents the status of a model
type ModelStatus struct {
	Name    string `json:"name"`
	Running bool   `json:"running"`
	Error   string `json:"error,omitempty"`
}

// Manager handles model operations
type Manager struct {
	modelsDir string
}

// NewManager creates a new model manager
func NewManager(modelsDir string) *Manager {
	return &Manager{
		modelsDir: modelsDir,
	}
}

// GetRecommendedModel returns the recommended model for the current system
func (m *Manager) GetRecommendedModel() (*ModelInfo, error) {
	// Simplified recommendation logic
	return &ModelInfo{
		Name:  "codellama:7b",
		Size:  "3.8GB",
		Speed: "fast",
	}, nil
}

// ListInstalledModels returns a list of installed models
func (m *Manager) ListInstalledModels() ([]*ModelInfo, error) {
	// This would normally query Ollama for installed models
	// For now, return a simple list
	models := []*ModelInfo{
		{
			Name:  "codellama:7b",
			Size:  "3.8GB",
			Speed: "fast",
		},
	}
	
	return models, nil
}

// InstallModel installs a model via Ollama
func (m *Manager) InstallModel(modelName string) error {
	cmd := exec.Command("ollama", "pull", modelName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install model %s: %v\nOutput: %s", modelName, err, string(output))
	}
	
	return nil
}

// GetModelStatus returns the status of a specific model
func (m *Manager) GetModelStatus(modelName string) (*ModelStatus, error) {
	// Check if model is running by listing running models
	cmd := exec.Command("ollama", "ps")
	output, err := cmd.Output()
	if err != nil {
		return &ModelStatus{
			Name:    modelName,
			Running: false,
			Error:   fmt.Sprintf("failed to check status: %v", err),
		}, nil
	}
	
	// Simple check if model name appears in output
	running := strings.Contains(string(output), modelName)
	
	return &ModelStatus{
		Name:    modelName,
		Running: running,
	}, nil
}

// IsModelInstalled checks if a model is installed
func (m *Manager) IsModelInstalled(modelName string) bool {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	return strings.Contains(string(output), modelName)
}