package mongo

import (
	"log"

	"github.com/jiffy-backend/helper"
	"github.com/jiffy-backend/config"
)

type DB struct{}

// registers a contract
func (DB *DB) RegisterContract(contract ContractObj) error {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	c := NewContractService(session.Copy(), config.DBNAME, config.ContractCollection)
	err = c.Register(contract)
	if err != nil {
		return err
	}
	helper.DBLogger.Info("Successfully added new contract", "Contract", contract.String())
	return nil
}

// register label per contract
func (DB *DB) RegisterLabel(name string) (err error) {

	return nil
}

// get all contracts
func (DB *DB) GetContracts() (contracts []ContractObj, err error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	c := NewContractService(session.Copy(), config.DBNAME, config.ContractCollection)
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
	c := NewContractService(session.Copy(), config.DBNAME, config.ContractCollection)
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
	c := NewContractService(session.Copy(), config.DBNAME, config.ContractCollection)
	err = c.GetContractByAddress(&contract, addr)
	if err != nil {
		helper.DBLogger.Error("Unable to get contract", "Address", addr, "Error", err)
		return contract, err
	}
	helper.DBLogger.Info("Fetched contract", "Contract", contract.String())
	return contract, nil
}

func (DB *DB) GetContractByIdentifier(hash [32]byte) (ContractObj,error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	c := NewContractService(session.Copy(), config.DBNAME, config.ContractCollection)
	var contract ContractObj
	err = c.GetContractByIdentifier(hash,&contract)
	if err != nil {
		helper.DBLogger.Error("Unable to get contract by hash","Error", err)
		return contract, err
	}

	return contract,nil
}
