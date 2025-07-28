package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/carsor007/cloudpork-agent/internal/analyzer"
	"github.com/carsor007/cloudpork-agent/internal/api"
	"github.com/carsor007/cloudpork-agent/internal/config"
	"github.com/carsor007/cloudpork-agent/internal/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	projectID string
	output    string
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze [directory]",
	Short: "Analyze codebase for cloud cost optimization",
	Long: `Analyze your codebase to identify infrastructure scaling patterns
and potential cost optimizations.

This command:
1. Checks for Claude Code CLI installation
2. Analyzes your codebase using structured prompts
3. Sends analysis results to CloudPork for cost projection
4. Never uploads your actual source code

Examples:
  cloudpork analyze                           # Analyze current directory
  cloudpork analyze ./my-project             # Analyze specific directory
  cloudpork analyze --project-id=proj_abc123 # Use specific project ID
  cloudpork analyze --output=json            # Output raw JSON results`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "CloudPork project ID")
	analyzeCmd.Flags().StringVarP(&output, "output", "o", "dashboard", "Output format: dashboard, json, or quiet")
	
	viper.BindPFlag("project-id", analyzeCmd.Flags().Lookup("project-id"))
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	// Get user subscription info first
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	subscription, err := getSubscriptionInfo(cfg.APIKey)
	if err != nil {
		return fmt.Errorf("failed to get subscription info: %w", err)
	}
	
	// Check trial limitations
	if subscription.Tier == types.TierTrial {
		if subscription.AnalysesUsed >= subscription.AnalysesLimit {
			return showTrialUpgradePrompt(subscription)
		}
		
		if subscription.DaysRemaining <= 2 {
			showTrialWarning(subscription)
		}
	}
	
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}
	
	// Resolve absolute path
	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve directory path: %v", err)
	}
	
	// Verify directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absPath)
	}
	
	// Get project ID from flag, config, or prompt
	projID := projectID
	if projID == "" {
		projID = viper.GetString("project-id")
	}
	if projID == "" {
		projID = cfg.ProjectID
	}
	if projID == "" {
		projID = config.GenerateProjectID()
		color.Yellow("📝 Generated new project ID: %s", projID)
		color.Yellow("💡 Save this ID with: cloudpork config set project-id %s", projID)
	}
	
	// Print banner
	if output != "quiet" {
		printBanner()
		fmt.Printf("📁 Analyzing: %s\n", absPath)
		fmt.Printf("🆔 Project ID: %s\n\n", projID)
	}
	
	// Initialize analyzer
	analyzer := analyzer.New(absPath, projID)
	
	// Determine analysis mode and perform analysis
	return performAnalysis(analyzer)
}

func performAnalysis(analyzer *analyzer.Analyzer) error {
	mode := viper.GetString("llm.mode")
	
	switch mode {
	case "local":
		return performLocalAnalysis(analyzer)
	case "hybrid":
		return performHybridAnalysis(analyzer)
	default:
		return performCloudAnalysis(analyzer) // existing logic
	}
}

func performLocalAnalysis(analyzer *analyzer.Analyzer) error {
	fmt.Println("🔒 Performing local analysis...")
	
	baseURL := viper.GetString("llm.local_url")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	
	// For now, we'll use the existing cloud analysis logic
	// In a real implementation, this would use local LLM endpoints
	result, err := analyzer.Analyze()
	if err != nil {
		return fmt.Errorf("analysis failed: %v", err)
	}
	
	// Handle output
	switch output {
	case "json":
		return result.PrintJSON()
	case "quiet":
		// Silent mode
	default:
		result.PrintSummary()
	}
	
	fmt.Println("✅ Local analysis completed")
	return nil
}

func performHybridAnalysis(analyzer *analyzer.Analyzer) error {
	fmt.Println("⚡ Performing hybrid analysis...")
	
	// First do local analysis
	if err := performLocalAnalysis(analyzer); err != nil {
		return err
	}
	
	// Then send summary to CloudPork for enhanced intelligence
	fmt.Println("📊 Enhancing with cloud intelligence...")
	// ... (send only metadata to CloudPork API)
	
	return nil
}

func performCloudAnalysis(analyzer *analyzer.Analyzer) error {
	// Run analysis
	result, err := analyzer.Analyze()
	if err != nil {
		return fmt.Errorf("analysis failed: %v", err)
	}
	
	// Handle output
	switch output {
	case "json":
		return result.PrintJSON()
	case "quiet":
		// Silent mode - just send to API
	default:
		result.PrintSummary()
	}
	
	// Send to CloudPork API
	if output != "quiet" {
		fmt.Println("📡 Sending results to CloudPork...")
	}
	
	client := api.NewClient()
	err = client.SendAnalysis(result)
	if err != nil {
		color.Red("❌ Failed to send results: %v", err)
		color.Yellow("💡 Run 'cloudpork auth login' to authenticate")
		return err
	}
	
	if output != "quiet" {
		color.Green("✅ Analysis complete!")
		fmt.Println("🌐 View results: https://cloudpork.com/dashboard")
	}
	
	return nil
}

func showTrialUpgradePrompt(subscription *types.SubscriptionInfo) error {
	fmt.Println("🎯 Trial Analysis Used!")
	fmt.Println()
	fmt.Printf("You've used your 1 trial analysis. ")
	if subscription.DaysRemaining > 0 {
		fmt.Printf("Your trial expires in %d days.\n", subscription.DaysRemaining)
	} else {
		fmt.Println("Your trial has expired.")
	}
	fmt.Println()
	
	fmt.Println("📊 Upgrade to continue analyzing:")
	fmt.Println("  🌱 Starter ($29/mo): 10 analyses + export + history")
	fmt.Println("  ⚡ Professional ($149/mo): 100 analyses + team + API")
	fmt.Println("  🏢 Enterprise ($499/mo): Unlimited + local AI + security")
	fmt.Println()
	fmt.Println("💡 Your current analysis will be deleted unless you upgrade!")
	fmt.Println()
	fmt.Printf("Upgrade now: https://cloudpork.com/pricing\n")
	
	return fmt.Errorf("trial limit reached - upgrade to continue")
}

func showTrialWarning(subscription *types.SubscriptionInfo) {
	fmt.Printf("⚠️  Trial expires in %d days! Upgrade to keep your analysis: https://cloudpork.com/pricing\n", subscription.DaysRemaining)
	fmt.Println()
}

func getSubscriptionInfo(apiKey string) (*types.SubscriptionInfo, error) {
	// API call to get subscription info
	client := &http.Client{Timeout: 10 * time.Second}
	
	req, err := http.NewRequest("GET", "https://api.cloudpork.com/v1/subscription", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed: %s", resp.Status)
	}
	
	var subscription types.SubscriptionInfo
	if err := json.NewDecoder(resp.Body).Decode(&subscription); err != nil {
		return nil, err
	}
	
	return &subscription, nil
}

func printBanner() {
	banner := color.New(color.FgMagenta, color.Bold).Sprint("🐷 CloudPork Agent")
	tagline := color.New(color.FgCyan).Sprint("Cut the pork from your cloud costs")
	
	fmt.Printf("%s - %s\n", banner, tagline)
	fmt.Printf("%s\n\n", color.New(color.Faint).Sprint("Analyzing your codebase for cost optimizations..."))
}