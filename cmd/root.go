package cmd

import (
	"collectify/internal/cli"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "collectify",
	Short: "Collectify is a versatile application for managing your collections.",
	Long: `
	Collectify is a powerful and flexible application designed to help you manage various types of collections.
	From web-based interfaces to text user interfaces (TUI) and graphical user interfaces (GUI), 
	Collectify provides multiple ways to interact with your data.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default to Web
		cli.DoWeb()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
