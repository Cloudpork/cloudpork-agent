package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configFileName = ".cloudpork"
	apiKeyEnvVar   = "CLOUDPORK_API_KEY"
	projectIDEnvVar = "CLOUDPORK_PROJECT_ID"
)

// GetAPIKey retrieves the API key from config or environment
func GetAPIKey() (string, error) {
	// Check environment variable first
	if key := os.Getenv(apiKeyEnvVar); key != "" {
		return key, nil
	}
	
	// Check config file
	if key := viper.GetString("api_key"); key != "" {
		return key, nil
	}
	
	return "", fmt.Errorf("no API key found. Run 'cloudpork auth login' to authenticate")
}

// SetAPIKey stores the API key in the config file
func SetAPIKey(apiKey string) error {
	viper.Set("api_key", apiKey)
	return saveConfig()
}

// GetProjectID retrieves the project ID from config or environment
func GetProjectID() (string, error) {
	// Check environment variable first
	if id := os.Getenv(projectIDEnvVar); id != "" {
		return id, nil
	}
	
	// Check config file
	if id := viper.GetString("project_id"); id != "" {
		return id, nil
	}
	
	return "", fmt.Errorf("no project ID found")
}

// SetProjectID stores the project ID in the config file
func SetProjectID(projectID string) error {
	viper.Set("project_id", projectID)
	return saveConfig()
}

// GenerateProjectID creates a new project ID
func GenerateProjectID() string {
	// Generate 8 random bytes
	bytes := make([]byte, 8)
	rand.Read(bytes)
	
	// Return as "proj_" + hex string
	return "proj_" + hex.EncodeToString(bytes)
}

// ClearCredentials removes stored credentials
func ClearCredentials() error {
	viper.Set("api_key", "")
	viper.Set("project_id", "")
	return saveConfig()
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName+".yaml"), nil
}

// saveConfig writes the current configuration to disk
func saveConfig() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %v", err)
	}
	
	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	
	// Write config file
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	
	// Set secure permissions on config file (readable only by owner)
	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("failed to set config file permissions: %v", err)
	}
	
	return nil
}

// LoadConfig loads configuration from file and environment
func LoadConfig() error {
	// Set config file name and paths
	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	
	// Add config paths
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home)
	}
	viper.AddConfigPath(".")
	
	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CLOUDPORK")
	
	// Read config file (ignore if not found)
	viper.ReadInConfig()
	
	return nil
}

// IsAuthenticated checks if user has valid credentials
func IsAuthenticated() bool {
	apiKey, err := GetAPIKey()
	return err == nil && apiKey != ""
}

// GetVerbose returns whether verbose mode is enabled
func GetVerbose() bool {
	return viper.GetBool("verbose")
}