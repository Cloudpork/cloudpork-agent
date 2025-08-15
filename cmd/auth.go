package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/Cloudpork/cloudpork-agent/internal/config"
	"github.com/Cloudpork/cloudpork-agent/internal/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  "Manage CloudPork authentication and trial signup",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to CloudPork",
	Long: `Login to your CloudPork account using your API key.

You can get your API key from: https://cloudpork.com/settings/api-keys

The API key will be stored securely in your system keychain.`,
	RunE: runLogin,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from CloudPork",
	Long:  `Remove your stored CloudPork credentials.`,
	RunE:  runLogout,
}

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Start your 7-day free trial",
	Long: `Start your 7-day free trial with CloudPork.
Get 1 analysis to see how much you can save on cloud costs.

After signup, you'll receive:
• 1 codebase analysis (choose your most important project!)
• Basic cost projections and optimization recommendations
• 7 days to see the value before deciding to upgrade

No credit card required for trial.`,
	RunE: runSignup,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check your subscription status",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(signupCmd)
	authCmd.AddCommand(statusCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	fmt.Println("🔐 CloudPork Authentication")
	fmt.Println()
	
	// Prompt for API key
	fmt.Print("Enter your CloudPork API key: ")
	byteAPIKey, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read API key: %v", err)
	}
	fmt.Println() // New line after password input
	
	apiKey := strings.TrimSpace(string(byteAPIKey))
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	
	// Validate API key format
	if !strings.HasPrefix(apiKey, "cp_") {
		color.Yellow("⚠️  API key should start with 'cp_'")
	}
	
	// Store API key
	err = config.SetAPIKey(apiKey)
	if err != nil {
		return fmt.Errorf("failed to store API key: %v", err)
	}
	
	// Optional: Prompt for project ID
	fmt.Print("Enter your default project ID (optional): ")
	reader := bufio.NewReader(os.Stdin)
	projectID, _ := reader.ReadString('\n')
	projectID = strings.TrimSpace(projectID)
	
	if projectID != "" {
		err = config.SetProjectID(projectID)
		if err != nil {
			color.Yellow("⚠️  Failed to store project ID: %v", err)
		}
	}
	
	color.Green("✅ Successfully authenticated with CloudPork!")
	if projectID != "" {
		color.Green("📝 Default project ID set to: %s", projectID)
	}
	fmt.Println("🌐 Dashboard: https://cloudpork.com/dashboard")
	
	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	err := config.ClearCredentials()
	if err != nil {
		return fmt.Errorf("failed to clear credentials: %v", err)
	}
	
	color.Green("✅ Successfully logged out from CloudPork")
	return nil
}

func runSignup(cmd *cobra.Command, args []string) error {
	fmt.Println("🐷 Welcome to CloudPork!")
	fmt.Println("Let's start your 7-day free trial...")
	fmt.Println()
	
	// Get user info
	email, err := promptForInput("Work email: ")
	if err != nil {
		return err
	}
	
	name, err := promptForInput("Your name: ")
	if err != nil {
		return err
	}
	
	company, err := promptForInput("Company (optional): ")
	if err != nil {
		return err
	}
	
	// Create trial account
	trialInfo, err := createTrialAccount(email, name, company)
	if err != nil {
		return fmt.Errorf("failed to create trial: %w", err)
	}
	
	// Save config
	err = config.SetAPIKey(trialInfo.APIKey)
	if err != nil {
		return fmt.Errorf("failed to save API key: %w", err)
	}
	
	err = config.SetProjectID(trialInfo.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to save project ID: %w", err)
	}
	
	fmt.Println("🎉 Trial activated!")
	fmt.Printf("   • Trial ends: %s\n", trialInfo.TrialEndsAt.Format("January 2, 2006"))
	fmt.Printf("   • Analyses remaining: %d\n", trialInfo.AnalysesRemaining)
	fmt.Printf("   • Project ID: %s\n", trialInfo.ProjectID)
	fmt.Println()
	fmt.Println("🚀 Ready to analyze! Run: cloudpork analyze")
	fmt.Println("💡 Choose your most important project - you get 1 analysis!")
	
	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("not logged in. Run: cloudpork auth login")
	}
	
	subscription, err := getSubscriptionInfo(cfg.APIKey)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}
	
	fmt.Println("📊 CloudPork Status")
	fmt.Println("==================")
	fmt.Printf("Plan: %s\n", strings.Title(subscription.Tier))
	fmt.Printf("Status: %s\n", strings.Title(subscription.Status))
	
	if subscription.IsTrialing {
		fmt.Printf("Trial ends: %s\n", subscription.TrialEndsAt.Format("January 2, 2006"))
		fmt.Printf("Days remaining: %d\n", subscription.DaysRemaining)
	}
	
	if subscription.AnalysesLimit == -1 {
		fmt.Printf("Analyses used: %d (unlimited)\n", subscription.AnalysesUsed)
	} else {
		fmt.Printf("Analyses used: %d/%d\n", subscription.AnalysesUsed, subscription.AnalysesLimit)
	}
	
	fmt.Printf("Project ID: %s\n", cfg.ProjectID)
	
	if subscription.Tier == types.TierTrial {
		fmt.Println()
		fmt.Println("🎯 Upgrade Options:")
		fmt.Println("   🌱 Starter ($29/mo): 10 analyses/month")
		fmt.Println("   ⚡ Professional ($149/mo): 100 analyses/month + team features")
		fmt.Println("   🏢 Enterprise ($499/mo): Unlimited + local AI + security")
		fmt.Println()
		fmt.Println("   Upgrade: https://cloudpork.com/pricing")
	}
	
	return nil
}

type TrialInfo struct {
	APIKey             string    `json:"api_key"`
	ProjectID          string    `json:"project_id"`
	TrialEndsAt        time.Time `json:"trial_ends_at"`
	AnalysesRemaining  int       `json:"analyses_remaining"`
}

func createTrialAccount(email, name, company string) (*TrialInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	payload := map[string]string{
		"email":   email,
		"name":    name,
		"company": company,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	
	resp, err := client.Post(
		"https://api.cloudpork.com/v1/auth/trial",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("signup failed: %s", resp.Status)
	}
	
	var trialInfo TrialInfo
	if err := json.NewDecoder(resp.Body).Decode(&trialInfo); err != nil {
		return nil, err
	}
	
	return &trialInfo, nil
}

func promptForInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "cp_***"
	}
	return apiKey[:3] + "***" + apiKey[len(apiKey)-4:]
}