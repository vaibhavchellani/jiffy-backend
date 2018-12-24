package src

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

type DB struct{}

// establishes connection with mongodb
func EstablishConnection() (client *mongo.Client, err error) {
	var dblogger = Logger.With("module", "database")

	connectionCtx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err = mongo.Connect(connectionCtx, SERVER)
	if err != nil {
		dblogger.Error("Unable to connect to mongo", "Error", err)
		return client, err
	}

	pingCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		dblogger.Error("Unable to ping to mongo", "Error", err)
		return client, err
	}

	return client, nil
}

// registers a contract
func (db *DB) RegisterContract(contract ContractObj) error {
	var dblogger = Logger.With("module", "database")

	client, err := EstablishConnection()
	if err != nil {
		return err
	}

	contractsInstance := client.Database(DBNAME).Collection(ContractCollection)
	insertContext, _ := context.WithTimeout(context.Background(), 1*time.Second)
	res, err := contractsInstance.InsertOne(insertContext, contract)
	if err != nil {
		return err
	}
	dblogger.Info("Successfully added new contract", "Address", contract.Address, "Name", contract.Name, "ID", res.InsertedID)
	return nil
}
