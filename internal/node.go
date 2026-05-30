// code to start the raft service via nodes

package internal

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

type Config struct {
	ID       string
	HTTPAddr string
	RaftAddr string
	JoinAddr string
	RaftDir  string
	Inmem    bool
}

func Run(ctx context.Context, cfg Config) error {
	//set default id if no id is present
	if cfg.ID == "" {
		cfg.ID = cfg.RaftAddr
	}

	if err := os.MkdirAll(cfg.RaftAddr, 0o700); err != nil {
		return fmt.Errorf("create raft dir: %w", err)
	}

	s := store.New(cfg.Inmem)
	//set the current directory to the raft directory
	// and also the bind port
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

	defer h.Close()
	if cfg.JoinAddr != "" {
		if err := join(cfg.JoinAddr, cfg.RaftAddr, cfg.ID); err != nil {
			return fmt.Errorf("join node at %s: %w", cfg.JoinAddr, err)
		}
	}

	log.Printf("node %s listening on http://%s, raft %s", cfg.ID, cfg.HTTPAddr, cfg.RaftAddr)

	<-ctx.Done()

	return nil

}

func join(joinAddr, raftAddr, nodeId string) error {
	b, err := json.Marshal(map[string]string{
		"addr": raftAddr,
		"id":   nodeId,
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

	return nil
}
