package analyzer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/carsor007/cloudpork-agent/internal/claude"
	"github.com/carsor007/cloudpork-agent/internal/types"
	"github.com/fatih/color"
)

// Analyzer performs code analysis using Claude Code
type Analyzer struct {
	projectDir string
	projectID  string
	claude     *claude.Client
}

// New creates a new analyzer instance
func New(projectDir, projectID string) *Analyzer {
	return &Analyzer{
		projectDir: projectDir,
		projectID:  projectID,
		claude:     claude.New(projectDir),
	}
}

// Analyze performs comprehensive code analysis
func (a *Analyzer) Analyze() (*types.CodeAnalysis, error) {
	// Pre-flight checks
	if err := a.preflightChecks(); err != nil {
		return nil, err
	}
	
	// Run Claude Code analysis
	result, err := a.claude.Analyze(a.projectID)
	if err != nil {
		return nil, fmt.Errorf("claude analysis failed: %v", err)
	}
	
	// Post-process results
	a.postProcess(result)
	
	return result, nil
}

// preflightChecks validates prerequisites
func (a *Analyzer) preflightChecks() error {
	// Check if directory exists and is accessible
	if _, err := os.Stat(a.projectDir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", a.projectDir)
	}
	
	// Check if it looks like a code project
	if !a.isCodeProject() {
		color.Yellow("⚠️  Directory doesn't appear to contain a typical code project")
		color.Yellow("   Continuing anyway, but results may be limited")
	}
	
	// Check Claude Code installation
	if !claude.IsInstalled() {
		color.Red("❌ Claude Code CLI not found")
		fmt.Println()
		fmt.Println(claude.GetInstallInstructions())
		return fmt.Errorf("Claude Code CLI not installed")
	}
	
	return nil
}

// isCodeProject checks if directory contains typical code project files
func (a *Analyzer) isCodeProject() bool {
	// Common project indicators
	indicators := []string{
		"package.json",     // Node.js
		"requirements.txt", // Python
		"go.mod",          // Go
		"Cargo.toml",      // Rust
		"composer.json",   // PHP
		"Gemfile",         // Ruby
		"pom.xml",         // Java (Maven)
		"build.gradle",    // Java (Gradle)
		"Dockerfile",      // Docker
		".git",           // Git repository
	}
	
	for _, indicator := range indicators {
		if _, err := os.Stat(filepath.Join(a.projectDir, indicator)); err == nil {
			return true
		}
	}
	
	// Check for source code directories
	sourceDirs := []string{"src", "lib", "app", "components", "pages"}
	for _, dir := range sourceDirs {
		if info, err := os.Stat(filepath.Join(a.projectDir, dir)); err == nil && info.IsDir() {
			return true
		}
	}
	
	return false
}

// postProcess enhances analysis results
func (a *Analyzer) postProcess(analysis *types.CodeAnalysis) {
	// Set estimated users based on complexity and endpoints
	analysis.EstimatedUsers = a.estimateCurrentUsers(analysis)
	
	// Add missing default values
	if analysis.Language == "" {
		analysis.Language = "Unknown"
	}
	if analysis.Framework == "" {
		analysis.Framework = "Unknown"
	}
	if analysis.ComplexityScore == 0 {
		analysis.ComplexityScore = 50 // Default medium complexity
	}
	
	// Validate resource estimates
	a.validateResourceEstimates(analysis)
}

// estimateCurrentUsers estimates current user base from codebase complexity
func (a *Analyzer) estimateCurrentUsers(analysis *types.CodeAnalysis) int {
	// Base estimate on complexity and infrastructure patterns
	baseUsers := 1000 // Default assumption
	
	// Adjust based on complexity
	complexityMultiplier := float64(analysis.ComplexityScore) / 50.0
	
	// Adjust based on endpoints (more endpoints = more features = more users)
	endpointMultiplier := float64(analysis.ApiEndpoints) / 10.0
	if endpointMultiplier < 0.5 {
		endpointMultiplier = 0.5
	}
	
	// Adjust based on background jobs (indicates more complex operations)
	jobMultiplier := 1.0
	if len(analysis.BackgroundJobs) > 0 {
		jobMultiplier = 1.5
	}
	
	estimated := int(float64(baseUsers) * complexityMultiplier * endpointMultiplier * jobMultiplier)
	
	// Reasonable bounds
	if estimated < 100 {
		estimated = 100
	}
	if estimated > 1000000 {
		estimated = 1000000
	}
	
	return estimated
}

// validateResourceEstimates ensures resource estimates are reasonable
func (a *Analyzer) validateResourceEstimates(analysis *types.CodeAnalysis) {
	resources := &analysis.ResourceUsage
	
	// Memory validation (minimum 128MB, maximum 16GB)
	if resources.MemoryMB < 128 {
		resources.MemoryMB = 128
	}
	if resources.MemoryMB > 16384 {
		resources.MemoryMB = 16384
	}
	
	// CPU validation (minimum 0.1 cores, maximum 32 cores)
	if resources.CPUCores < 0.1 {
		resources.CPUCores = 0.1
	}
	if resources.CPUCores > 32.0 {
		resources.CPUCores = 32.0
	}
	
	// Database connections (minimum 1, maximum 1000)
	if resources.DatabaseConns < 1 {
		resources.DatabaseConns = 1
	}
	if resources.DatabaseConns > 1000 {
		resources.DatabaseConns = 1000
	}
	
	// Network bandwidth (minimum 1Mbps, maximum 10Gbps)
	if resources.NetworkMbps < 1 {
		resources.NetworkMbps = 1  
	}
	if resources.NetworkMbps > 10000 {
		resources.NetworkMbps = 10000
	}
	
	// Storage (minimum 1GB, maximum 10TB)
	if resources.StorageGB < 1 {
		resources.StorageGB = 1
	}
	if resources.StorageGB > 10000 {
		resources.StorageGB = 10000
	}
}