package cmd

import (
	"fmt"
	"os"

	"github.com/Cloudpork/cloudpork-agent/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudpork",
	Short: "Cut the pork from your cloud costs 🐷",
	Long: `CloudPork analyzes your codebase to identify cloud cost optimization opportunities.

🎯 New to CloudPork? Start your free trial:
   cloudpork auth signup

🔍 Analyze your codebase:
   cloudpork analyze

📊 Check your subscription:
   cloudpork auth status

Cut the pork from your cloud costs with intelligent analysis!`,
	
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Check if user needs to sign up (except for auth commands)
		if cmd.Name() != "auth" && !isAuthenticated() {
			fmt.Println("👋 Welcome to CloudPork!")
			fmt.Println("Start your free trial: cloudpork auth signup")
			fmt.Println()
		}
	},
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

func isAuthenticated() bool {
	cfg, err := config.LoadConfig()
	return err == nil && cfg.APIKey != ""
}