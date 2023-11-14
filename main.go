package main

import (
	"github.com/dyammarcano/eventReceiverQuel/cmd"
	"github.com/dyammarcano/eventReceiverQuel/internal/metadata"
)

var (
	Version    = "v0.0.1-manual-build"
	CommitHash string
	Date       string
)

func init() {
	metadata.Set(Version, CommitHash, Date)
}

func main() {
	cmd.Execute()
}
