package cmd

import (
	"fmt"
	"github.com/dyammarcano/eventReceiverQuel/internal/metadata"
	"github.com/spf13/cobra"
	"os"
	"runtime/trace"
)

func PrintVersion(cmd *cobra.Command, _ []string) {
	defer trace.StartRegion(cmd.Context(), "version").End()
	// clean console
	fmt.Fprintf(cmd.OutOrStdout(), "\033[H\033[2J")

	// print version
	fmt.Fprintf(cmd.OutOrStdout(), metadata.String())
	os.Exit(0)
}
