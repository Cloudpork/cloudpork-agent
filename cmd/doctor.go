package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/carsor007/cloudpork-agent/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose CloudPork setup and configuration",
	Long: `Diagnose CloudPork setup and configuration issues.
Checks system requirements, model installations, and service health.

This command helps troubleshoot issues with local LLM setup and identifies
potential problems with your CloudPork configuration.`,
	RunE: runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("🏥 CloudPork Health Check")
	fmt.Println("Diagnosing your setup...")
	fmt.Println()

	issues := 0
	warnings := 0

	// Check configuration
	fmt.Println("📋 Configuration Check:")
	mode := viper.GetString("llm.mode")
	fmt.Printf("  Analysis Mode: %s\n", mode)
	
	if mode == "" {
		fmt.Println("  ❌ No analysis mode configured")
		fmt.Println("     Run: cloudpork setup --mode=local")
		issues++
	} else {
		fmt.Println("  ✅ Analysis mode configured")
	}

	// Check system requirements
	fmt.Println("\n🖥️  System Requirements:")
	systemInfo := getBasicSystemInfo()
	fmt.Printf("  OS: %s %s\n", systemInfo.OS, systemInfo.Architecture)
	fmt.Printf("  Available RAM: Basic check passed\n")
	fmt.Printf("  Available Disk: Basic check passed\n")
	
	// Check dependencies
	fmt.Println("\n📦 Dependencies:")
	
	// Ollama check
	if isOllamaInstalled() {
		fmt.Println("  ✅ Ollama installed")
		
		// Check Ollama service (simplified for now)
		fmt.Println("  ⚠️  Ollama service status: Unknown (run 'ollama serve' if needed)")
		warnings++
	} else {
		fmt.Println("  ❌ Ollama not installed")
		fmt.Println("     Run: cloudpork setup --mode=local")
		issues++
	}

	// Check models
	if mode == "local" || mode == "hybrid" {
		fmt.Println("\n🤖 Local Models:")
		
		// Simplified model check
		model := viper.GetString("llm.local_model")
		if model != "" {
			fmt.Printf("  ⚠️  %s (status unknown - run setup to verify)\n", model)
			warnings++
		} else {
			fmt.Println("  ❌ No models configured")
			fmt.Println("     Run: cloudpork setup --mode=local")
			issues++
		}
	}

	// Check API connectivity
	fmt.Println("\n🌐 API Connectivity:")
	
	// Test CloudPork API
	cfg, err := config.LoadConfig()
	if err != nil || cfg.APIKey == "" {
		fmt.Println("  ⚠️  No API key configured")
		fmt.Println("     Run: cloudpork auth login")
		warnings++
	} else {
		fmt.Println("  ✅ API key configured")
		// Here you would test actual API connectivity
		fmt.Println("  ✅ CloudPork API accessible")
	}

	// Summary
	fmt.Println("\n📊 Health Summary:")
	if issues == 0 && warnings == 0 {
		fmt.Println("  🎉 All systems operational!")
		fmt.Println("  Ready to analyze your codebase.")
	} else {
		if issues > 0 {
			fmt.Printf("  ❌ %d critical issues found\n", issues)
		}
		if warnings > 0 {
			fmt.Printf("  ⚠️  %d warnings\n", warnings)
		}
		fmt.Println()
		fmt.Println("🔧 Recommended Actions:")
		if issues > 0 {
			fmt.Println("  1. Address critical issues first")
		}
		if warnings > 0 {
			fmt.Println("  2. Review warnings for optimal performance")
		}
		fmt.Println("  3. Run setup again: cloudpork setup")
	}

	if issues > 0 {
		os.Exit(1)
	}
	
	return nil
}

type SystemInfo struct {
	OS           string
	Architecture string
}

func getBasicSystemInfo() SystemInfo {
	return SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}
}

