package llm

import (
	"encoding/json"
	"net/http"
	"time"
)

// IsOllamaHealthy checks if Ollama service is running and responding
func IsOllamaHealthy(baseURL string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get(baseURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == 200
}

// IsModelAvailable checks if a specific model is available in Ollama
func IsModelAvailable(baseURL, modelName string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get(baseURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	var response struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false
	}
	
	for _, model := range response.Models {
		if model.Name == modelName {
			return true
		}
	}
	
	return false
}