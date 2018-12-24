package src



import (
	tmlog "github.com/tendermint/tendermint/libs/log"
	logger "github.com/tendermint/tendermint/libs/log"

)
var Logger logger.Logger


// CheckpointLogger for checkpoint module logger
var log tmlog.Logger

func init() {
	//log = Logger.With("module", "checkpoint")
}
