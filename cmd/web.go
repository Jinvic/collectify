package cmd

import (
	"collectify/internal/cli"

	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web interface of Collectify.",
	Long: `
	Starts the web-based interface for Collectify, allowing you to manage your collections through a browser.
	This command boots up a web server that serves the interactive web pages for your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoWeb()
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}
