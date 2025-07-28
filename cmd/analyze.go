package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/carsor007/cloudpork-agent/internal/analyzer"
	"github.com/carsor007/cloudpork-agent/internal/api"
	"github.com/carsor007/cloudpork-agent/internal/config"
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

func printBanner() {
	banner := color.New(color.FgMagenta, color.Bold).Sprint("🐷 CloudPork Agent")
	tagline := color.New(color.FgCyan).Sprint("Cut the pork from your cloud costs")
	
	fmt.Printf("%s - %s\n", banner, tagline)
	fmt.Printf("%s\n\n", color.New(color.Faint).Sprint("Analyzing your codebase for cost optimizations..."))
}