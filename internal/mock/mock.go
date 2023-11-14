package mock

import (
	"context"
	"github.com/spf13/cobra"
)

func CobraMockCommand() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	return cmd
}
