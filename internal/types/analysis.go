package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Tier constants for subscription tiers
const (
	TierTrial        = "trial"        // 7-day trial
	TierStarter      = "starter"
	TierProfessional = "professional"
	TierEnterprise   = "enterprise"
)

// SubscriptionInfo represents user subscription details
type SubscriptionInfo struct {
	Tier           string     `json:"tier"`
	Status         string     `json:"status"`
	AnalysesUsed   int        `json:"analyses_used"`
	AnalysesLimit  int        `json:"analyses_limit"`
	TrialEndsAt    *time.Time `json:"trial_ends_at,omitempty"`
	IsTrialing     bool       `json:"is_trialing"`
	DaysRemaining  int        `json:"days_remaining,omitempty"`
}

// CodeAnalysis represents the complete analysis result
type CodeAnalysis struct {
	ProjectID        string          `json:"project_id"`
	Timestamp        time.Time       `json:"timestamp"`
	Directory        string          `json:"directory"`
	Language         string          `json:"language"`
	Framework        string          `json:"framework"`
	Dependencies     []string        `json:"dependencies"`
	DatabaseCalls    int             `json:"database_calls"`
	ApiEndpoints     int             `json:"api_endpoints"`
	StatelessFuncs   int             `json:"stateless_functions"`
	BackgroundJobs   []string        `json:"background_jobs"`
	CacheUsage       []string        `json:"cache_usage"`
	FileUploads      bool            `json:"file_uploads"`
	ComplexityScore  int             `json:"complexity_score"`
	ScalingBottlenecks []Bottleneck  `json:"scaling_bottlenecks"`
	ResourceUsage    ResourceMetrics `json:"resource_usage"`
	EstimatedUsers   int             `json:"estimated_users"`
	SecurityIssues   []SecurityIssue `json:"security_issues"`
	Performance      PerformanceMetrics `json:"performance"`
}

// ResourceMetrics represents estimated resource requirements
type ResourceMetrics struct {
	MemoryMB       int     `json:"memory_mb"`
	CPUCores       float64 `json:"cpu_cores"`
	DatabaseConns  int     `json:"database_connections"`
	NetworkMbps    int     `json:"network_mbps"`
	StorageGB      int     `json:"storage_gb"`
}

// Bottleneck represents a scaling bottleneck
type Bottleneck struct {
	Type        string `json:"type"`        // "database", "cpu", "memory", "network"
	Description string `json:"description"`
	Severity    string `json:"severity"`    // "low", "medium", "high", "critical"
	Impact      string `json:"impact"`      // Description of impact
}

// SecurityIssue represents a security concern
type SecurityIssue struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Line        int    `json:"line,omitempty"`
	File        string `json:"file,omitempty"`
}

// PerformanceMetrics represents performance characteristics
type PerformanceMetrics struct {
	AvgResponseTime  int  `json:"avg_response_time_ms"`
	DatabaseQueries  int  `json:"database_queries_per_request"`
	CacheHitRate     int  `json:"cache_hit_rate_percent"`
	HasNPlusOneQuery bool `json:"has_n_plus_one_query"`
	HasLargePayloads bool `json:"has_large_payloads"`
}

// PrintSummary prints a formatted summary of the analysis
func (ca *CodeAnalysis) PrintSummary() {
	fmt.Printf("%s\n", color.New(color.FgCyan, color.Bold).Sprint("📊 Analysis Summary"))
	fmt.Printf("%s\n\n", color.New(color.Faint).Sprint(strings.Repeat("=", 50)))
	
	// Basic Info
	fmt.Printf("🏗️  %s: %s\n", color.New(color.Bold).Sprint("Framework"), ca.Framework)
	fmt.Printf("💾  %s: %s\n", color.New(color.Bold).Sprint("Language"), ca.Language)
	fmt.Printf("📦  %s: %d\n", color.New(color.Bold).Sprint("Dependencies"), len(ca.Dependencies))
	fmt.Printf("🔌  %s: %d\n", color.New(color.Bold).Sprint("API Endpoints"), ca.ApiEndpoints)
	fmt.Printf("⚡  %s: %d\n", color.New(color.Bold).Sprint("Background Jobs"), len(ca.BackgroundJobs))
	fmt.Println()
	
	// Resource Usage
	fmt.Printf("%s\n", color.New(color.FgGreen, color.Bold).Sprint("💻 Resource Requirements"))
	fmt.Printf("  Memory: %d MB\n", ca.ResourceUsage.MemoryMB)
	fmt.Printf("  CPU: %.1f cores\n", ca.ResourceUsage.CPUCores)
	fmt.Printf("  DB Connections: %d\n", ca.ResourceUsage.DatabaseConns)
	fmt.Printf("  Network: %d Mbps\n", ca.ResourceUsage.NetworkMbps)
	fmt.Println()
	
	// Bottlenecks
	if len(ca.ScalingBottlenecks) > 0 {
		fmt.Printf("%s\n", color.New(color.FgYellow, color.Bold).Sprint("⚠️  Scaling Bottlenecks"))
		for _, bottleneck := range ca.ScalingBottlenecks {
			severity := getSeverityColor(bottleneck.Severity)
			fmt.Printf("  %s %s: %s\n", 
				severity.Sprint("●"), 
				color.New(color.Bold).Sprint(bottleneck.Type), 
				bottleneck.Description)
		}
		fmt.Println()
	}
	
	// Performance Issues
	if ca.Performance.HasNPlusOneQuery || ca.Performance.HasLargePayloads {
		fmt.Printf("%s\n", color.New(color.FgRed, color.Bold).Sprint("🐌 Performance Issues"))
		if ca.Performance.HasNPlusOneQuery {
			fmt.Println("  • N+1 query patterns detected")
		}
		if ca.Performance.HasLargePayloads {
			fmt.Println("  • Large payload responses found")
		}
		fmt.Println()
	}
	
	// Complexity Score
	complexityColor := getComplexityColor(ca.ComplexityScore)
	fmt.Printf("🎯 %s: %s\n", 
		color.New(color.Bold).Sprint("Complexity Score"), 
		complexityColor.Sprintf("%d/100", ca.ComplexityScore))
}

// PrintJSON prints the analysis as JSON
func (ca *CodeAnalysis) PrintJSON() error {
	jsonData, err := json.MarshalIndent(ca, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

func getSeverityColor(severity string) *color.Color {
	switch severity {
	case "critical":
		return color.New(color.FgRed, color.Bold)
	case "high":
		return color.New(color.FgRed)
	case "medium":
		return color.New(color.FgYellow)
	case "low":
		return color.New(color.FgGreen)
	default:
		return color.New(color.FgWhite)
	}
}

func getComplexityColor(score int) *color.Color {
	switch {
	case score >= 80:
		return color.New(color.FgRed, color.Bold)
	case score >= 60:
		return color.New(color.FgYellow, color.Bold)
	case score >= 40:
		return color.New(color.FgCyan, color.Bold)
	default:
		return color.New(color.FgGreen, color.Bold)
	}
}