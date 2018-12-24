package src

import (
	"github.com/ethereum/go-ethereum/common"
)

type ContractObj struct {
	Name    string `bson:"name"`
	Address string `bson:"address`
	Network string `bson:"network"`
	ABI     []byte `bson:"abi"`
}

type Contracts []ContractObj

type User struct {
	Address common.Address `json:"user_address"`
}
