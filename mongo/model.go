package mongo

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"encoding/hex"
)

type ContractObj struct {
	Name      string `bson:"name"`
	Address   string `bson:"contract_address"`
	Network   string `bson:"network"`
	ABI       []byte `bson:"abi"`
	QueryName string `bson:"queryable_name"`
	Owner     string `bson:"owner_address"`
	Identifier [32]byte `bson:contract_hash`
}

func (c *ContractObj) String() string {
	result := fmt.Sprintf("Contract--> name: %v addr: %v chain: %v", c.Name, c.Address, c.Network)
	return result
}

type contractObjJson struct{
	Name      string `json:"name"`
	Address   string `json:"contract_address"`
	Network   string `json:"network"`
	ABI       []byte `json:"abi"`
	Owner     string `json:"owner_address"`
	Identifier string `json:contract_hash`
}

func (c *ContractObj) Json() (contractObjJson) {
	identifier:=hex.EncodeToString(c.Identifier[:])
	contract:=contractObjJson{
		Name:c.Name,
		Address:c.Address,
		Network:c.Network,
		ABI:c.ABI,
		Owner:c.Owner,
		Identifier:identifier,
	}
	return contract
}

type User struct {
	Address common.Address `json:"user_address"`
}
