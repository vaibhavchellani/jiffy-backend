package src

import (
	logger "github.com/tendermint/tendermint/libs/log"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

var Logger logger.Logger

// CheckpointLogger for checkpoint module logger
var log tmlog.Logger

func init() {
	//log = Logger.With("module", "checkpoint")
}
