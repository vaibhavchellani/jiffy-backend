package src

import (
	tmlog "github.com/tendermint/tendermint/libs/log"
	"os"
)

// CheckpointLogger for checkpoint module logger
var Logger tmlog.Logger
var DBLogger tmlog.Logger
var ControllerLogger tmlog.Logger
func init() {
	Logger = tmlog.NewTMLogger(tmlog.NewSyncWriter(os.Stdout))
	DBLogger = Logger.With("module", "database")
	ControllerLogger = Logger.With("module", "controller")

}
