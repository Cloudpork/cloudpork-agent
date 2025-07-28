package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudpork",
	Short: "Cut the pork from your cloud costs",
	Long: color.New(color.FgMagenta, color.Bold).Sprint("🐷 CloudPork Agent") + `

Analyze your codebase to identify wasteful cloud spending and 
optimize your infrastructure scaling.

The CloudPork agent uses Claude Code to analyze your project locally,
then sends only the analysis summary to generate cost projections.

Examples:
  cloudpork analyze                    # Analyze current directory
  cloudpork analyze --project-id=123  # Analyze with specific project ID
  cloudpork auth login                 # Authenticate with CloudPork
  cloudpork version                    # Show version information`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloudpork.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "enable verbose output")
	
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cloudpork")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && viper.GetBool("verbose") {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}