package cmd

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// These variables are set by goreleaser during build
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CloudPork agent version",
	Long:  `Display version information for the CloudPork agent.`,
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("%s %s\n", 
		color.New(color.FgMagenta, color.Bold).Sprint("CloudPork Agent"),
		color.New(color.FgCyan).Sprint("v"+version))
	
	if version != "dev" {
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built: %s\n", date)
	}
	
	fmt.Printf("Go: %s\n", runtime.Version())
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}