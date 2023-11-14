package cmd

import (
	"github.com/dyammarcano/template-go/internal/cmd"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Execute API commands",
	Long: `The "api" command executes various API operations using defined subcommands.
For example, to get the information about a user with "api get-user",

You can further specify options with each subcommand to tailor the requests.`,
	RunE: cmd.CallExternalAPI,
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
