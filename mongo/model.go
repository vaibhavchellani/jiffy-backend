package mongo

import (
	"encoding/hex"
	"fmt"
	abi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jiffy-backend/helper"
)

// Make sure to update JSON when doing this
type ContractObj struct {
	Name        string   `bson:"name"`
	Address     string   `bson:"contract_address"`
	NetworkName string   `bson:"network_name"`
	ABI         []byte   `bson:"abi"`
	QueryName   string   `bson:"queryable_name"`
	Owner       string   `bson:"owner_address"`
	Identifier  [32]byte `bson:"contract_hash"`
	NetworkURL  string   `bson:"network_url"`
}

// generate string representation for contract
func (c *ContractObj) String() string {
	result := fmt.Sprintf("Contract--> name: %v addr: %v chain: %v owner:%v", c.Name, c.Address, c.NetworkName,c.Owner)
	return result
}

// prettify the contract obj for the response
type ContractObjJson struct {
	Name       string  `json:"name"`
	Address    string  `json:"contract_address"`
	NetworkName string   `bson:"network_name"`
	ABI        abi.ABI `json:"abi"`
	Owner      string  `json:"owner_address"`
	Identifier string  `json:"contract_hash"`
	NetworkURL  string   `json:"network_url"`
}
func (c *ContractObj) Json() ContractObjJson {
	identifier := hex.EncodeToString(c.Identifier[:])
	var abi abi.ABI
	helper.UnMarshallABI(c.ABI, &abi)
	contract := ContractObjJson{
		Name:       c.Name,
		Address:    c.Address,
		NetworkName:    c.NetworkName,
		ABI:        abi,
		Owner:      c.Owner,
		Identifier: identifier,
		NetworkURL:c.NetworkURL,
	}
	return contract
}

func (c *ContractObj) ValidateBasic() (error){

	// TODO add other input checks

	//if helper.IsValidAddress(m.Address) {
	//	err := errors.New("Contract address is not valid ethereum address")
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//if helper.IsValidAddress(m.Owner) {
	//	err := errors.New("Owner address is not valid ethereum address")
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Write([]byte(err.Error()))
	//	return
	//}

	return nil
}


type User struct {
	Address common.Address `json:"user_address"`
}
