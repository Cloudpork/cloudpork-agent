package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/carsor007/cloudpork-agent/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage CloudPork authentication",
	Long: `Manage your CloudPork authentication credentials.

Use these commands to login, logout, and check your authentication status.`,
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

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Long:  `Check if you're currently authenticated with CloudPork.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
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

func runStatus(cmd *cobra.Command, args []string) error {
	apiKey, err := config.GetAPIKey()
	if err != nil || apiKey == "" {
		color.Red("❌ Not authenticated")
		fmt.Println("Run 'cloudpork auth login' to authenticate")
		return nil
	}
	
	// Mask API key for display
	maskedKey := maskAPIKey(apiKey)
	
	projectID, _ := config.GetProjectID()
	
	color.Green("✅ Authenticated with CloudPork")
	fmt.Printf("🔑 API Key: %s\n", maskedKey)
	if projectID != "" {
		fmt.Printf("📝 Project ID: %s\n", projectID)
	}
	fmt.Println("🌐 Dashboard: https://cloudpork.com/dashboard")
	
	return nil
}

func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "cp_***"
	}
	return apiKey[:3] + "***" + apiKey[len(apiKey)-4:]
}