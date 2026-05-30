// code to start the raft service via nodes

package internal

import "context"

type Config struct {
	ID       string
	HTTPAddr string
	RaftAddr string
	JoinAddr string
	RaftDir  string
	Inmem    bool
}

func Run(ctx context.Context, cfg Config) error {
	//TODO
	// mkdir raft dir
	// create raft store
	// open raft
	// start the HTTP service
	// join the leader if cfg.JoinAddr != ""
	// wait for the the context cancellation
}
