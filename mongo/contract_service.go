package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/jiffy-backend/config"
	"log"
	"strings"
	"github.com/jiffy-backend/helper"
)

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

func NewContractService(session *Session, dbName string) *ContractService {
	collection := session.GetCollection(dbName, config.ContractCollection)
	collection.EnsureIndex(contractModelIndex())
	return &ContractService{collection}
}

func (c *ContractService) Register(contract ContractObj) error {
	objectID := bson.NewObjectId()
	contract.ID = objectID
	err := c.collection.Insert(contract)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContractService) GetAllContracts(contracts *[]ContractObj) (err error) {
	err = c.collection.Find(nil).All(&contracts)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (c *ContractService) GetContractByName(contract *ContractObj, name string) (err error) {
	err = c.collection.Find(bson.D{{"queryable_name", strings.ToLower(name)}}).One(&contract)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContractService) GetContractByAddress(contract *ContractObj, address string) (err error) {
	err = c.collection.Find(bson.D{{"contract_address", strings.ToLower(address)}}).One(&contract)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContractService) GetContractByIdentifier(hash string, contract *ContractObj) (err error) {
	err = c.collection.Find(bson.D{{"contract_hash", hash}}).One(&contract)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContractService) AddLabel(_id bson.ObjectId,contractID bson.ObjectId) (err error) {
	selector:= bson.M{"_id":contractID}
	// TODO add label id to contract obj
	changeInfo,err :=c.collection.Upsert(selector,bson.M{"$push":bson.M{"label_id":_id}})
	if err!=nil{
		return err
	}
	helper.DBLogger.Debug("Added label to contract","ChangeInfo",changeInfo)
	return nil
}
