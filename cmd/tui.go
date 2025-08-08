package cmd

import (
	"collectify/internal/cli"

	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the text user interface of Collectify.",
	Long: `
	Launches the Text User Interface (TUI) of Collectify.
	This interface allows you to interact with your collections directly from the terminal using keyboard inputs,
	providing an efficient way to manage your collections without needing a graphical environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoTui()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
