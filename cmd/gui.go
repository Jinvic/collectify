package cmd

import (
	"collectify/internal/cli"

	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Open the graphical user interface of Collectify.",
	Long: `
	Opens the Graphical User Interface (GUI) of Collectify.
	With this interface, you can visually manage your collections using a mouse and graphical controls.
	It's designed to be intuitive and user-friendly.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoGui()
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
}
