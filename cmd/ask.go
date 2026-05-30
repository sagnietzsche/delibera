package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sagnikc395/delibera/internal"
	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask <message>",
	Short: "Ask the configured LLM provider",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		provider, err := internal.NewInferenceProviderFromEnv()
		if err != nil {
			return err
		}

		response, err := provider.Prompt(ctx, strings.Join(args, " "))
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), response)
		return err
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
}
