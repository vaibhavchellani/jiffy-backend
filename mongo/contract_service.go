package mongo

import (
	"github.com/globalsign/mgo"
	"log"
	"strings"
	"github.com/globalsign/mgo/bson"
)

type ContractService struct {
	collection *mgo.Collection
}

func contractModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func NewContractService(session *Session, dbName string, collectionName string) *ContractService {
	collection := session.GetCollection(dbName, collectionName)
	collection.EnsureIndex(contractModelIndex())
	return &ContractService {collection}
}

func (c *ContractService) Register(contract ContractObj) error {
	err :=c.collection.Insert(contract)
	if err!=nil{
		return err
	}
	return nil
}

func (c *ContractService) GetAllContracts(contracts *[]ContractObj) (err error)  {
	err = c.collection.Find(nil).All(&contracts)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (c *ContractService) GetContractByName(contract *ContractObj, name string)(err error)  {
	err=c.collection.Find(bson.D{{"queryable_name", strings.ToLower(name)}}).One(&contract)
	if err != nil {
		return err
	}
	return nil
}
func (c *ContractService) GetContractByAddress(contract *ContractObj, address string)(err error)  {
	err=c.collection.Find(bson.D{{"contract_address", strings.ToLower(address)}}).One(&contract)
	if err != nil {
		return err
	}
	return nil
}

