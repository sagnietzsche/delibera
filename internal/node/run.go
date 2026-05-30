package node

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	httpd "github.com/otoolep/hraftd/http"
	"github.com/otoolep/hraftd/store"
)

// Config contains the runtime settings for one Raft-backed Delibera node.
type Config struct {
	ID       string
	HTTPAddr string
	RaftAddr string
	JoinAddr string
	RaftDir  string
	Inmem    bool
}

// Run starts one node and blocks until ctx is canceled.
func Run(ctx context.Context, cfg Config) error {
	if cfg.RaftDir == "" {
		return fmt.Errorf("raft data path is required")
	}

	if cfg.ID == "" {
		cfg.ID = cfg.RaftAddr
	}

	if err := os.MkdirAll(cfg.RaftDir, 0o700); err != nil {
		return fmt.Errorf("create raft dir: %w", err)
	}

	s := store.New(cfg.Inmem)
	s.RaftDir = cfg.RaftDir
	s.RaftBind = cfg.RaftAddr

	if err := s.Open(cfg.JoinAddr == "", cfg.ID); err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	h := httpd.New(cfg.HTTPAddr, s)
	if err := h.Start(); err != nil {
		return fmt.Errorf("start http service: %w", err)
	}

	if cfg.JoinAddr != "" {
		if err := join(cfg.JoinAddr, cfg.RaftAddr, cfg.ID); err != nil {
			return fmt.Errorf("join node at %s: %w", cfg.JoinAddr, err)
		}
	}

	log.Printf("node %s listening on http://%s, raft %s", cfg.ID, cfg.HTTPAddr, cfg.RaftAddr)

	<-ctx.Done()
	return nil
}

func join(joinAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(map[string]string{
		"addr": raftAddr,
		"id":   nodeID,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/join", joinAddr),
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("join returned %s", resp.Status)
	}

	return nil
}
