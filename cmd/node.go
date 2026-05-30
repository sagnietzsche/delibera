package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sagnikc395/delibera/internal"
	"github.com/sagnikc395/delibera/internal/node"
	"github.com/spf13/cobra"
)

var nodeCfg node.Config

var nodeCmd = &cobra.Command{
	Use:   "node <raft-data-path>",
	Short: "Start one Delibera raft node",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nodeCfg.RaftDir = args[0]

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		return node.Run(ctx, nodeCfg)
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)

	nodeCmd.Flags().BoolVar(&nodeCfg.Inmem, "inmem", false, "Use in-memory Raft storage")
	nodeCmd.Flags().StringVar(&nodeCfg.HTTPAddr, "haddr", internal.DefaultHTTPAddr, "HTTP bind address")
	nodeCmd.Flags().StringVar(&nodeCfg.RaftAddr, "raddr", internal.DefaultRaftAddr, "Raft bind address")
	nodeCmd.Flags().StringVar(&nodeCfg.JoinAddr, "join", "", "Leader HTTP address to join")
	nodeCmd.Flags().StringVar(&nodeCfg.ID, "id", "", "Node ID")
}
