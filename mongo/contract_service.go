package mongo

import (
	"log"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/jiffy-backend/config"
	"github.com/jiffy-backend/helper"
)

type IContractService interface {
	Register(contract ContractObj) error
	GetAllContracts(contracts *[]ContractObj) (err error)
	GetContractByName(contract *ContractObj, name string) (err error)
	GetContractByAddress(contract *ContractObj, address string) (err error)
	GetContractByIdentifier(hash string, contract *ContractObj) (err error)
	AddLabel(_id bson.ObjectId, contractID bson.ObjectId) (err error)
}

type ContractService struct {
	collection *mgo.Collection
}

func contractModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"queryable_name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

// factory method for generating new contract service
func NewContractService(session *Session, dbName string) *ContractService {
	collection := session.GetCollection(dbName, config.ContractCollection)
	collection.EnsureIndex(contractModelIndex())
	return &ContractService{collection}
}

// register new contract
func (c *ContractService) Register(contract ContractObj) error {
	objectID := bson.NewObjectId()
	contract.ID = objectID
	err := c.collection.Insert(contract)
	if err != nil {
		return err
	}
	return nil
}

// get all contracts deployed
func (c *ContractService) GetAllContracts(contracts *[]ContractObj) (err error) {
	err = c.collection.Find(nil).All(&contracts)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// get contract by contract name
func (c *ContractService) GetContractByName(contract *ContractObj, name string) (err error) {
	err = c.collection.Find(bson.D{{"queryable_name", strings.ToLower(name)}}).One(&contract)
	if err != nil {
		return err
	}
	return nil
}

// get contract by contract address
func (c *ContractService) GetContractByAddress(contract *ContractObj, address string) (err error) {
	err = c.collection.Find(bson.D{{"contract_address", strings.ToLower(address)}}).One(&contract)
	if err != nil {
		return err
	}
	return nil
}

// get contract by hash (address+network)
func (c *ContractService) GetContractByIdentifier(hash string, contract *ContractObj) (err error) {
	err = c.collection.Find(bson.D{{"contract_hash", hash}}).One(&contract)
	if err != nil {
		return err
	}
	return nil
}

// add label id to a contract
func (c *ContractService) AddLabel(_id bson.ObjectId, contractID bson.ObjectId) (err error) {
	selector := bson.M{"_id": contractID}
	changeInfo, err := c.collection.Upsert(selector, bson.M{"$push": bson.M{"label_id": _id}})
	if err != nil {
		return err
	}
	helper.DBLogger.Debug("Added label to contract", "ChangeInfo", changeInfo)
	return nil
}

// // get all labels by contract address
// func (c *ContractService) GetLabels(contractAddr string) (err error) {
// 	selector := bson.M{"_id": contractID}

// 	return nil
// }
