package mongo

import (
	"log"

	"github.com/globalsign/mgo/bson"
	"github.com/jiffy-backend/config"
	"github.com/jiffy-backend/helper"
)

type DB struct{}

// registers a contract
func (DB *DB) RegisterContract(contract ContractObj) error {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	c := NewContractService(session.Copy(), config.DBNAME)
	err = c.Register(contract)
	if err != nil {
		return err
	}
	helper.DBLogger.Info("Successfully added new contract", "Contract", contract.String())
	return nil
}

// get all contracts
func (DB *DB) GetContracts() (contracts []ContractObj, err error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	c := NewContractService(session.Copy(), config.DBNAME)
	err = c.GetAllContracts(&contracts)
	if err != nil {
		helper.DBLogger.Error("Cannot fetch all contracts", "Error", err)
		return contracts, err
	}
	return contracts, nil
}

// get a contract by name
func (DB *DB) GetContractByName(name string) (contract ContractObj, err error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	c := NewContractService(session.Copy(), config.DBNAME)
	err = c.GetContractByName(&contract, name)
	if err != nil {
		helper.DBLogger.Error("Unable to get contract", "Name", name, "Error", err)
		return contract, err
	}
	helper.DBLogger.Info("Fetched contract", "Contract", contract.String(), "Name", name)
	return contract, nil
}

// get a contract by address
func (DB *DB) GetContractByAddr(addr string) (contract ContractObj, err error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	c := NewContractService(session.Copy(), config.DBNAME)
	err = c.GetContractByAddress(&contract, addr)
	if err != nil {
		helper.DBLogger.Error("Unable to get contract", "Address", addr, "Error", err)
		return contract, err
	}
	helper.DBLogger.Info("Fetched contract", "Contract", contract.String())
	return contract, nil
}

func (DB *DB) GetContractByIdentifier(hash string) (ContractObj, error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	c := NewContractService(session.Copy(), config.DBNAME)
	var contract ContractObj
	err = c.GetContractByIdentifier(hash, &contract)
	if err != nil {
		helper.DBLogger.Error("Unable to get contract by hash", "Error", err, "hash", hash)
		return contract, err
	}
	return contract, nil
}

// ---- Label related ops

// register label
func (DB *DB) RegisterLabel(label Label) (err error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	l := NewLabelService(session.Copy(), config.DBNAME)
	c := NewContractService(session.Copy(), config.DBNAME)
	labelID := bson.NewObjectId()
	label.ID = labelID

	// TODO make this operation atomic
	// first add lable to contract then create entity

	err = l.Register(label)
	if err != nil {
		helper.DBLogger.Error("Unable to register label", "Error", err, "Label", label.String())
		return err
	}else{
		err = c.AddLabel(labelID,label.ContractID)
		if err!=nil{
			helper.DBLogger.Error("Unable to add label to contract", "Error", err, "Label", label.String(),"Contract",label.ContractName)
			return err
		}
	}

	return nil
}

// -------
