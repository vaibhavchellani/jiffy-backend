package src

import (

	"log"

	"github.com/globalsign/mgo"
	"strings"
	"github.com/globalsign/mgo/bson"
)

type DB struct{}

// establishes connection with mongodb
func (DB *DB) DBConnect() (session *mgo.Session, err error) {
	session,err =mgo.Dial(SERVER)
	if err != nil {
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
	c:=session.DB(DBNAME).C(ContractCollection)
	err =c.Insert(contract)
	if err!=nil{
		return err
	}
	//DBLogger.Info("Successfully added new contract", "Contract", contract.String(), "ID", res.InsertedID)
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
	c:=session.DB(DBNAME).C(ContractCollection)
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
	c:=session.DB(DBNAME).C(ContractCollection)

	DBLogger.Debug("Searching contract", "name", name)

	err=c.Find(bson.D{{"queryable_name", strings.ToLower(name)}}).One(&contract)
	if err != nil {
			DBLogger.Error("Unable to get contract", "Name", name, "Error", err)
			return contract, err
	}

	DBLogger.Info("Fetched contract", "Contract", contract.String())

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
	c:=session.DB(DBNAME).C(ContractCollection)
	DBLogger.Debug("Searching contract", "Address", addr)
	err=c.Find(bson.D{{"contract_address", addr}}).One(&contract)
	if err != nil {
		DBLogger.Error("Unable to get contract", "Address", addr, "Error", err)
		return contract, err
	}
	DBLogger.Info("Fetched contract", "Contract", contract.String())
	return contract, nil
}
