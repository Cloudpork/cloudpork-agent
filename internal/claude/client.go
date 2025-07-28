package claude

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/carsor007/cloudpork-agent/internal/types"
)

// Client handles interactions with Claude Code CLI
type Client struct {
	projectDir string
}

// New creates a new Claude Code client
func New(projectDir string) *Client {
	return &Client{
		projectDir: projectDir,
	}
}

// IsInstalled checks if Claude Code CLI is available
func IsInstalled() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}

// GetInstallInstructions returns installation instructions for Claude Code
func GetInstallInstructions() string {
	return `Claude Code CLI not found. Install it with:

📥 Installation options:
  • Web: https://claude.ai/cli
  • macOS: brew install claude-ai/tap/claude
  • Linux: curl -fsSL https://claude.ai/install.sh | sh
  • Windows: Download from https://claude.ai/cli/download

After installation, run: claude auth login`
}

// Analyze performs comprehensive code analysis using Claude Code
func (c *Client) Analyze(projectID string) (*types.CodeAnalysis, error) {
	if !IsInstalled() {
		return nil, fmt.Errorf("Claude Code CLI not installed")
	}
	
	analysis := &types.CodeAnalysis{
		ProjectID: projectID,
		Timestamp: time.Now(),
		Directory: c.projectDir,
	}
	
	// Run multiple analysis passes
	fmt.Print("🔍 Running code analysis")
	
	// 1. Basic project analysis
	fmt.Print(".")
	if err := c.analyzeBasicStructure(analysis); err != nil {
		return nil, fmt.Errorf("basic analysis failed: %v", err)
	}
	
	// 2. Database and API analysis
	fmt.Print(".")
	if err := c.analyzeDatabaseAndAPI(analysis); err != nil {
		return nil, fmt.Errorf("database/API analysis failed: %v", err)
	}
	
	// 3. Performance and scaling analysis
	fmt.Print(".")
	if err := c.analyzePerformanceAndScaling(analysis); err != nil {
		return nil, fmt.Errorf("performance analysis failed: %v", err)
	}
	
	// 4. Resource estimation
	fmt.Print(".")
	if err := c.estimateResources(analysis); err != nil {
		return nil, fmt.Errorf("resource estimation failed: %v", err)
	}
	
	fmt.Println(" ✅")
	
	return analysis, nil
}

// analyzeBasicStructure identifies language, framework, and dependencies
func (c *Client) analyzeBasicStructure(analysis *types.CodeAnalysis) error {
	prompt := `Analyze this codebase and identify:
1. Primary programming language
2. Web framework being used
3. Key dependencies and libraries
4. Number of API endpoints/routes
5. Background job processing (if any)
6. File upload capabilities

Respond in this JSON format:
{
  "language": "string",
  "framework": "string", 
  "dependencies": ["dep1", "dep2"],
  "api_endpoints": number,
  "background_jobs": ["job1", "job2"],
  "file_uploads": boolean
}`

	output, err := c.runClaudeCommand(prompt)
	if err != nil {
		return err
	}
	
	// Try to parse JSON response
	var result struct {
		Language       string   `json:"language"`
		Framework      string   `json:"framework"`
		Dependencies   []string `json:"dependencies"`
		ApiEndpoints   int      `json:"api_endpoints"`
		BackgroundJobs []string `json:"background_jobs"`
		FileUploads    bool     `json:"file_uploads"`
	}
	
	if err := json.Unmarshal([]byte(output), &result); err == nil {
		analysis.Language = result.Language
		analysis.Framework = result.Framework
		analysis.Dependencies = result.Dependencies
		analysis.ApiEndpoints = result.ApiEndpoints
		analysis.BackgroundJobs = result.BackgroundJobs
		analysis.FileUploads = result.FileUploads
	} else {
		// Fallback to heuristic parsing
		c.parseBasicStructureHeuristic(output, analysis)
	}
	
	return nil
}

// analyzeDatabaseAndAPI analyzes database usage and API patterns
func (c *Client) analyzeDatabaseAndAPI(analysis *types.CodeAnalysis) error {
	prompt := `Analyze database and API patterns in this codebase:
1. Count database queries/calls
2. Identify database connection patterns
3. Look for N+1 query problems
4. Find caching usage (Redis, Memcached, etc.)
5. Estimate complexity on a scale of 1-100

Focus on scalability concerns and potential bottlenecks.`

	output, err := c.runClaudeCommand(prompt)
	if err != nil {
		return err
	}
	
	// Parse response using heuristics
	analysis.DatabaseCalls = c.extractNumber(output, `(\d+).*(?:database|query)`)
	analysis.ComplexityScore = c.extractComplexity(output)
	analysis.CacheUsage = c.extractCacheUsage(output)
	
	// Check for N+1 queries
	analysis.Performance.HasNPlusOneQuery = strings.Contains(strings.ToLower(output), "n+1") || 
		strings.Contains(strings.ToLower(output), "n plus one")
	
	return nil
}

// analyzePerformanceAndScaling identifies scaling bottlenecks
func (c *Client) analyzePerformanceAndScaling(analysis *types.CodeAnalysis) error {
	prompt := `Identify scaling bottlenecks and performance issues:
1. Database connection limits
2. Memory-intensive operations  
3. CPU-heavy computations
4. Network bottlenecks
5. Synchronous operations that should be async
6. Large payload responses

For each issue, specify type (database/cpu/memory/network) and severity (low/medium/high/critical).`

	output, err := c.runClaudeCommand(prompt)
	if err != nil {
		return err
	}
	
	analysis.ScalingBottlenecks = c.extractBottlenecks(output)
	analysis.Performance.HasLargePayloads = strings.Contains(strings.ToLower(output), "large payload")
	
	return nil
}

// estimateResources calculates resource requirements
func (c *Client) estimateResources(analysis *types.CodeAnalysis) error {
	prompt := fmt.Sprintf(`Based on this %s/%s application with %d API endpoints and %d background jobs:

Estimate resource requirements for 1000 concurrent users:
1. Memory usage in MB
2. CPU cores needed  
3. Database connections required
4. Network bandwidth in Mbps
5. Storage requirements in GB

Consider the complexity score of %d and provide realistic estimates.`,
		analysis.Language, analysis.Framework, analysis.ApiEndpoints, 
		len(analysis.BackgroundJobs), analysis.ComplexityScore)

	output, err := c.runClaudeCommand(prompt)
	if err != nil {
		return err
	}
	
	// Parse resource estimates
	analysis.ResourceUsage.MemoryMB = c.extractNumber(output, `(\d+).*MB|(\d+).*memory`)
	analysis.ResourceUsage.CPUCores = c.extractFloat(output, `(\d+(?:\.\d+)?).*(?:core|cpu)`)
	analysis.ResourceUsage.DatabaseConns = c.extractNumber(output, `(\d+).*(?:connection|conn)`)
	analysis.ResourceUsage.NetworkMbps = c.extractNumber(output, `(\d+).*(?:Mbps|bandwidth)`)
	analysis.ResourceUsage.StorageGB = c.extractNumber(output, `(\d+).*(?:GB|storage)`)
	
	// Set defaults if parsing failed
	if analysis.ResourceUsage.MemoryMB == 0 {
		analysis.ResourceUsage.MemoryMB = c.estimateMemoryFromComplexity(analysis.ComplexityScore)
	}
	if analysis.ResourceUsage.CPUCores == 0 {
		analysis.ResourceUsage.CPUCores = c.estimateCPUFromEndpoints(analysis.ApiEndpoints)
	}
	
	return nil
}

// runClaudeCommand executes a Claude Code command with the given prompt
func (c *Client) runClaudeCommand(prompt string) (string, error) {
	cmd := exec.Command("claude", "code", "--prompt", prompt, "--directory", c.projectDir)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("claude command failed: %v\nOutput: %s", err, string(output))
	}
	
	return string(output), nil
}

// Heuristic parsing functions
func (c *Client) parseBasicStructureHeuristic(output string, analysis *types.CodeAnalysis) {
	lower := strings.ToLower(output)
	
	// Language detection
	languages := []string{"javascript", "python", "go", "java", "php", "ruby", "typescript"}
	for _, lang := range languages {
		if strings.Contains(lower, lang) {
			analysis.Language = strings.Title(lang)
			break
		}
	}
	
	// Framework detection  
	frameworks := []string{"react", "vue", "angular", "express", "fastapi", "django", "gin", "echo"}
	for _, framework := range frameworks {
		if strings.Contains(lower, framework) {
			analysis.Framework = strings.Title(framework)
			break
		}
	}
	
	// Endpoint counting
	analysis.ApiEndpoints = c.extractNumber(output, `(\d+).*(?:endpoint|route|api)`)
	if analysis.ApiEndpoints == 0 {
		analysis.ApiEndpoints = 5 // Default estimate
	}
}

func (c *Client) extractNumber(text, pattern string) int {
	re := regexp.MustCompile(`(?i)` + pattern)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return num
		}
	}
	return 0
}

func (c *Client) extractFloat(text, pattern string) float64 {
	re := regexp.MustCompile(`(?i)` + pattern)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		if num, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return num
		}
	}
	return 0.0
}

func (c *Client) extractComplexity(text string) int {
	// Look for complexity scores
	complexity := c.extractNumber(text, `(?:complexity|score).*?(\d+)`)
	if complexity == 0 {
		complexity = c.extractNumber(text, `(\d+).*(?:complexity|score)`)
	}
	if complexity == 0 || complexity > 100 {
		return 50 // Default medium complexity
	}
	return complexity
}

func (c *Client) extractCacheUsage(text string) []string {
	lower := strings.ToLower(text)
	var caches []string
	
	cacheTypes := []string{"redis", "memcached", "memory cache", "cdn", "browser cache"}
	for _, cache := range cacheTypes {
		if strings.Contains(lower, cache) {
			caches = append(caches, cache)
		}
	}
	
	return caches
}

func (c *Client) extractBottlenecks(text string) []types.Bottleneck {
	var bottlenecks []types.Bottleneck
	
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))
		if lower == "" {
			continue
		}
		
		// Look for bottleneck indicators
		if strings.Contains(lower, "database") && (strings.Contains(lower, "slow") || 
			strings.Contains(lower, "bottleneck") || strings.Contains(lower, "limit")) {
			bottlenecks = append(bottlenecks, types.Bottleneck{
				Type:        "database",
				Description: line,
				Severity:    c.extractSeverity(line),
				Impact:      "May cause slow response times under load",
			})
		}
		
		if strings.Contains(lower, "memory") && strings.Contains(lower, "leak") {
			bottlenecks = append(bottlenecks, types.Bottleneck{
				Type:        "memory",
				Description: line,
				Severity:    "high",
				Impact:      "Could cause application crashes",
			})
		}
	}
	
	return bottlenecks
}

func (c *Client) extractSeverity(text string) string {
	lower := strings.ToLower(text)
	if strings.Contains(lower, "critical") {
		return "critical"
	}
	if strings.Contains(lower, "high") {
		return "high"
	}
	if strings.Contains(lower, "medium") {
		return "medium"
	}
	return "low"
}

func (c *Client) estimateMemoryFromComplexity(complexity int) int {
	// Estimate memory based on complexity (MB)
	switch {
	case complexity >= 80:
		return 2048 // 2GB for very complex apps
	case complexity >= 60:
		return 1024 // 1GB for complex apps
	case complexity >= 40:
		return 512  // 512MB for medium apps
	default:
		return 256  // 256MB for simple apps
	}
}

func (c *Client) estimateCPUFromEndpoints(endpoints int) float64 {
	// Estimate CPU cores based on endpoints
	switch {
	case endpoints >= 50:
		return 4.0
	case endpoints >= 20:
		return 2.0
	case endpoints >= 10:
		return 1.0
	default:
		return 0.5
	}
}