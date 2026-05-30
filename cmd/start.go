package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/sagnikc395/delibera/internal"
	"github.com/spf13/cobra"
)

var (
	startDataDir string
	startNodes   int
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a local Delibera cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		if startNodes != 1 && startNodes != 3 {
			return fmt.Errorf("--nodes must be 1 or 3")
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		if err := os.MkdirAll(startDataDir, 0o700); err != nil {
			return fmt.Errorf("create data dir: %w", err)
		}

		exits := make(chan nodeExit, 3)
		var waits sync.WaitGroup

		stopStartedNodes := func(err error) error {
			stop()
			waits.Wait()
			return err
		}

		startNode := func(id string, httpPort, raftPort int, join string) error {
			nodeArgs := []string{
				"node",
				"--id", id,
				"--haddr", fmt.Sprintf("localhost:%d", httpPort),
				"--raddr", fmt.Sprintf("localhost:%d", raftPort),
			}

			if join != "" {
				nodeArgs = append(nodeArgs, "--join", join)
			}

			nodeArgs = append(nodeArgs, filepath.Join(startDataDir, id))

			proc := exec.CommandContext(ctx, os.Args[0], nodeArgs...)
			proc.Stdout = os.Stdout
			proc.Stderr = os.Stderr

			if err := proc.Start(); err != nil {
				return fmt.Errorf("start %s: %w", id, err)
			}

			waits.Add(1)
			go func() {
				defer waits.Done()
				exits <- nodeExit{id: id, err: proc.Wait()}
			}()

			return nil
		}

		if err := startNode("node1", 11000, 12000, ""); err != nil {
			return stopStartedNodes(err)
		}

		if startNodes == 3 {
			time.Sleep(time.Second)

			if err := startNode("node2", 11001, 12001, "localhost:11000"); err != nil {
				return stopStartedNodes(err)
			}
			if err := startNode("node3", 11002, 12002, "localhost:11000"); err != nil {
				return stopStartedNodes(err)
			}
		}

		fmt.Println("Cluster up:")
		fmt.Println("  node1  http://localhost:11000  raft: localhost:12000")
		if startNodes == 3 {
			fmt.Println("  node2  http://localhost:11001  raft: localhost:12001")
			fmt.Println("  node3  http://localhost:11002  raft: localhost:12002")
		}
		fmt.Println("Press Ctrl+C to stop.")

		select {
		case <-ctx.Done():
		case exit := <-exits:
			stopStartedNodes(nil)
			if exit.err != nil {
				return fmt.Errorf("%s exited: %w", exit.id, exit.err)
			}
			return nil
		}

		waits.Wait()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVar(&startDataDir, "data-dir", internal.DefaultDataDir, "Cluster data directory")
	startCmd.Flags().IntVar(&startNodes, "nodes", 1, "Number of nodes to start")
}

type nodeExit struct {
	id  string
	err error
}
