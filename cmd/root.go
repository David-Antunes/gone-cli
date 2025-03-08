package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gone-cli",
	Short: "CLI tool to interact with GONE",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().Bool("json", false, "Prints json request")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Does not make request. Useful for generating json request or curl request.")
	rootCmd.PersistentFlags().Bool("curl", false, "Prints curl request")
	rootCmd.PersistentFlags().Bool("quiet", false, "Supresses the sucessfull Response")
}
