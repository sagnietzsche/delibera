package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "delibera",
	Short: "Consensus-driven multi-agent reasoning",
	Long:  "Delibera runs local Raft nodes for consensus-driven multi-agent reasoning.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
