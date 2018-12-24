package src

import (
	tmlog "github.com/tendermint/tendermint/libs/log"
	"os"
)

// CheckpointLogger for checkpoint module logger
var Logger tmlog.Logger

func init() {
	Logger = tmlog.NewTMLogger(tmlog.NewSyncWriter(os.Stdout))
}
