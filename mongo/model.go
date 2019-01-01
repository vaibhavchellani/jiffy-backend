package mongo

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type ContractObj struct {
	Name      string `bson:"name"`
	Address   string `bson:"contract_address"`
	Network   string `bson:"network"`
	ABI       []byte `bson:"abi"`
	QueryName string `bson:"queryable_name"`
	Owner string `bson:"owner_address"`
}

func (c *ContractObj) String() string {
	result := fmt.Sprintf("Contract--> name: %v addr: %v chain: %v", c.Name, c.Address, c.Network)
	return result
}

type User struct {
	Address common.Address `json:"user_address"`
}
