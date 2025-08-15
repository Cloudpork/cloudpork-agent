package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Cloudpork/cloudpork-agent/internal/config"
	"github.com/Cloudpork/cloudpork-agent/internal/types"
)

const (
	defaultBaseURL = "https://api.cloudpork.com"
	defaultTimeout = 30 * time.Second
	userAgent      = "CloudPork-Agent/1.0"
)

// Client handles communication with CloudPork API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient() *Client {
	return &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// NewClientWithURL creates a new API client with custom base URL
func NewClientWithURL(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// SendAnalysis sends analysis results to CloudPork API
func (c *Client) SendAnalysis(analysis *types.CodeAnalysis) error {
	// Get API key from config
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return fmt.Errorf("no API key found: %v", err)
	}
	
	// Prepare request payload
	payload := struct {
		*types.CodeAnalysis
		AgentVersion string `json:"agent_version"`
		Platform     string `json:"platform"`
	}{
		CodeAnalysis: analysis,
		AgentVersion: "1.0.0", // TODO: Get from build info
		Platform:     "cli",
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal analysis: %v", err)
	}
	
	// Create HTTP request
	url := fmt.Sprintf("%s/v1/analysis", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("User-Agent", userAgent)
	
	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	
	// Handle response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	return nil
}

// ValidateAPIKey checks if the API key is valid
func (c *Client) ValidateAPIKey(apiKey string) error {
	url := fmt.Sprintf("%s/v1/auth/validate", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("User-Agent", userAgent)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid API key")
	}
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("validation failed with status %d", resp.StatusCode)
	}
	
	return nil
}

// GetProjectInfo retrieves project information
func (c *Client) GetProjectInfo(projectID string) (*ProjectInfo, error) {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return nil, fmt.Errorf("no API key found: %v", err)
	}
	
	url := fmt.Sprintf("%s/v1/projects/%s", c.baseURL, projectID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("User-Agent", userAgent)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("project not found")
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	
	var projectInfo ProjectInfo
	if err := json.NewDecoder(resp.Body).Decode(&projectInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	
	return &projectInfo, nil
}

// ProjectInfo represents project information from the API
type ProjectInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AnalysisCount int     `json:"analysis_count"`
	LastAnalysis  *time.Time `json:"last_analysis,omitempty"`
}