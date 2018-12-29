package src

import (
	"os"

	tmlog "github.com/tendermint/tendermint/libs/log"
)

// CheckpointLogger for checkpoint module logger
var Logger tmlog.Logger
var DBLogger tmlog.Logger
var ControllerLogger tmlog.Logger
var MainLogger tmlog.Logger

func init() {
	Logger = tmlog.NewTMLogger(tmlog.NewSyncWriter(os.Stdout))
	DBLogger = Logger.With("module", "database")
	ControllerLogger = Logger.With("module", "controller")
	MainLogger = Logger.With("module","main")
}
