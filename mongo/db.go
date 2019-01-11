package mongo

import (
	"log"

	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/jiffy-backend/config"
	"github.com/jiffy-backend/helper"
)

// ALL pre prossing before insert/update/get/delete to be applied here
// then forwarded to services

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

// get contract by hash
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
	// make contract name lowercase for query
	label.ContractName = strings.ToLower(label.ContractName)

	// TODO make this operation atomic
	// first add lable to contract then create entity

	err = l.Register(label)
	if err != nil {
		helper.DBLogger.Error("Unable to register label", "Error", err, "Label", label.String())
		return err
	}

	err = c.AddLabel(labelID, label.ContractID)
	if err != nil {
		helper.DBLogger.Error("Unable to add label to contract", "Error", err, "Label", label.String(), "Contract", label.ContractName)
		return err
	}

	return nil
}

// get all labels created by an address aka creator
func (DB *DB) GetLabelViaCreator(creatorAddr string) (labels []Label, err error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	l := NewLabelService(session.Copy(), config.DBNAME)
	l.GetLabelByCreator(creatorAddr, &labels)
	if err != nil {
		helper.DBLogger.Error("Cannot fetch all contracts", "Error", err)
		return labels, err
	}
	return labels, nil
}

// Gets labels by contract address
func (DB *DB) GetLabelViaContractAddr(contractAddr string) (labels []Label, err error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()

	// get contract by contract address
	contract, err := DB.GetContractByAddr(contractAddr)
	if err != nil {
		helper.DBLogger.Debug("Contract not found", "Address", contractAddr)
		return
	}

	// get labels using label IDS
	labels, err = DB.GetLabelByIDs(contract.Labels)
	if err != nil {
		helper.DBLogger.Debug("Lables not found for provided IDs", "IDs", contract.Labels)
		return nil, err
	}
	helper.DBLogger.Info("Successfully fetched all labels for contract", "ContractAddr", contractAddr, "Labels", labels, "Contract", contract.String())
	return labels, nil
}

// get labels by contract name
func (DB *DB) GetLabelViaContractName(name string) (labels []Label, err error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	l := NewLabelService(session.Copy(), config.DBNAME)
	l.GetLabelByCreator(name, &labels)
	if err != nil {
		helper.DBLogger.Error("Cannot fetch all contracts", "Error", err)
		return labels, err
	}
	return labels, nil
}

// get label by ID
func (DB *DB) GetLabelByID(labelID bson.ObjectId) (label Label, err error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	l := NewLabelService(session.Copy(), config.DBNAME)
	helper.DBLogger.Info("Fetching Label", "LabelID", labelID)
	err = l.GetLabelByID(labelID, &label)
	if err != nil {
		helper.DBLogger.Error("Unable to fetch label from DB", "ID", labelID)
		return label, err
	}
	return label, nil
}

// get all labels for given slice of IDs
func (DB *DB) GetLabelByIDs(labelIDs []bson.ObjectId) ([]Label, error) {
	session, err := NewSession(config.SERVER)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer session.Close()
	l := NewLabelService(session.Copy(), config.DBNAME)
	helper.DBLogger.Info("Fetching Labels", "LabelIDs", labelIDs)
	labels, err := l.GetLabelByIDS(labelIDs)
	if err != nil {
		helper.DBLogger.Error("Unable to fetch label from DB", "ID", labelIDs)
		return labels, err
	}
	return labels, nil
}

// -------
