package src

import (
	"github.com/ethereum/go-ethereum/common"
	"fmt"
)

type ContractObj struct {
	Name    string `bson:"name"`
	Address string `bson:"contract_address"`
	Network string `bson:"network"`
	ABI     []byte `bson:"abi"`
}

func (c *ContractObj) String() string {
	result:= fmt.Sprintf("Contract--> name: %v addr: %v chain: %v",c.Name,c.Address,c.Network)
	return result
}

type Contracts []ContractObj

type User struct {
	Address common.Address `json:"user_address"`
}
