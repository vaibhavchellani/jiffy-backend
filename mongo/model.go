package mongo

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
)

// Make sure to update JSON when doing this
type ContractObj struct {
	Name        string        `bson:"name"`             // unique name of contract
	Address     string        `bson:"contract_address"` // contract address
	NetworkName string        `bson:"network_name"`     // network name where contract is deployed
	ABI         string        `bson:"abi"`              // string representation of contract
	QueryName   string        `bson:"queryable_name"`   // lower case name
	Owner       string        `bson:"owner_address"`    // address of contract registrar on jiffy
	Identifier  string        `bson:"contract_hash"`    // hash(network+address)
	NetworkURL  string        `bson:"network_url"`      // network URL
	Cloned      string        `bson:"cloned_from"`      // name of previous dapp --> identified by identifier
	ID          bson.ObjectId `bson:"_id"`              // id for contract
	Labels      []Label       `bson:"labels"`
}

// generate string representation for contract
func (c *ContractObj) String() string {
	result := fmt.Sprintf("Contract--> name: %v addr: %v chain: %v owner:%v", c.Name, c.Address, c.NetworkName, c.Owner)
	return result
}

// prettify the contract obj for the response
type ContractObjJson struct {
	Name        string `json:"name"`
	Address     string `json:"contract_address"`
	NetworkName string `bson:"network_name"`
	ABI         string `json:"abi"`
	Owner       string `json:"owner_address"`
	Identifier  string `json:"contract_hash"`
	NetworkURL  string `json:"network_url"`
}

func (c *ContractObj) Json() ContractObjJson {
	contract := ContractObjJson{
		Name:        c.Name,
		Address:     c.Address,
		NetworkName: c.NetworkName,
		ABI:         c.ABI,
		Owner:       c.Owner,
		Identifier:  c.Identifier,
		NetworkURL:  c.NetworkURL,
	}
	return contract
}

func (c *ContractObj) ValidateBasic() error {

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

// Label Model
//--------------------

type Label struct {
	ContractName string        `bson:"contract_name"`
	ContractID   bson.ObjectId `bson:"contract_id"`
	CreatorAddr  string        `bson:"creator_addr"`
	Functions    []Function    `bson:"functions"`
	ID           bson.ObjectId `bson:"_id"`
	Name         string        `bson:"name"`
	Description  string        `bson:"description"`
}

func (l *Label) String() string {
	return fmt.Sprintf("Label ---> ContractName:%v LabelName:%v ", l.ContractName, l.Name)
}

//--------------------

type Function struct {
	MethodSig   []byte `bson:"method_sig"`
	Skippable   bool   `bson:"skippable"`
	Usage       string `bson:"usage"` // transaction or call
	Description string `bson:"description"`
}

//--------------------

type User struct {
	Address common.Address `bson:"user_address"`
}
