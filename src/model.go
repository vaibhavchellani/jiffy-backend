package src

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
)
type ContractObj struct {
	Name string `json:"name"`
	Address common.Address `json:"contract_address`
	ABI abi.ABI
}

type Contracts []ContractObj

type User struct {
	Address common.Address `json:"user_address"`

}

