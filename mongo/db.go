package mongo

import (

	"log"

	"github.com/globalsign/mgo"
	"strings"
	"github.com/globalsign/mgo/bson"
	"github.com/jiffy-backend/helper"
)

type DB struct{}

// establishes connection with mongodb
func (DB *DB) DBConnect() (session *mgo.Session, err error) {
	session,err =mgo.Dial(helper.SERVER)
	if err != nil {
		return session,err
	}
	c:=session.DB(helper.DBNAME).C(helper.ContractCollection)
	index := mgo.Index{
		Key: []string{"$text:name"},
		Unique:true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		helper.DBLogger.Error("index error")
		return session,err
	}

	return session, nil
}

// registers a contract
func (DB *DB) RegisterContract(contract ContractObj) error {
	session, err := DB.DBConnect()
	if err != nil {
		return err
	}
	defer session.Close()
	c:=session.DB(helper.DBNAME).C(helper.ContractCollection)
	//index := mgo.Index{
	//	Key: []string{"$text:name"},
	//	Unique:true,
	//}
	//err = c.EnsureIndex(index)
	//if err != nil {
	//	DBLogger.Error("index error")
	//	return err
	//}
	err =c.Insert(contract)
	if err!=nil{
		return err
	}
	helper.DBLogger.Info("Successfully added new contract", "Contract", contract.String())
	return nil
}

func (DB *DB) GetContracts() (contracts []ContractObj, err error) {
	// connecting to DB
	session, err := DB.DBConnect()
	if err != nil {
		log.Fatal(err)
		return contracts, err
	}
	defer session.Close()
	c:=session.DB(helper.DBNAME).C(helper.ContractCollection)
	err =c.Find(nil).All(&contracts)
	if err != nil {
		log.Fatal(err)
		return contracts, err
	}
	return contracts, nil
}

// get a contract by name
func (DB *DB) GetContractViaName(name string) (contract ContractObj, err error) {
	// connecting to DB
	session, err := DB.DBConnect()
	if err != nil {
		log.Fatal(err)
		return contract, err
	}
	defer session.Close()
	c:=session.DB(helper.DBNAME).C(helper.ContractCollection)

	helper.DBLogger.Debug("Searching contract", "name", name)

	err=c.Find(bson.D{{"queryable_name", strings.ToLower(name)}}).One(&contract)
	if err != nil {
			helper.DBLogger.Error("Unable to get contract", "Name", name, "Error", err)
			return contract, err
	}

	helper.DBLogger.Info("Fetched contract", "Contract", contract.String())

	return contract, nil
}

// get a contract by address
func (DB *DB) GetContractViaAddr(addr string) (contract ContractObj, err error) {
	// connecting to DB
	session, err := DB.DBConnect()
	if err != nil {
		log.Fatal(err)
		return contract, err
	}
	defer session.Close()
	c:=session.DB(helper.DBNAME).C(helper.ContractCollection)
	helper.DBLogger.Debug("Searching contract", "Address", addr)
	err=c.Find(bson.D{{"contract_address", addr}}).One(&contract)
	if err != nil {
		helper.DBLogger.Error("Unable to get contract", "Address", addr, "Error", err)
		return contract, err
	}
	helper.DBLogger.Info("Fetched contract", "Contract", contract.String())
	return contract, nil
}
