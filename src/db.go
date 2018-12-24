package src

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"log"
)

type DB struct{}

// establishes connection with mongodb
func (DB *DB)DBConnect() (client *mongo.Client, err error) {
	connectionCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	client, err = mongo.Connect(connectionCtx, SERVER)
	if err != nil {
		DBLogger.Error("Unable to connect to mongo", "Error", err)
		return client, err
	}


	err = client.Ping(connectionCtx, readpref.Primary())
	if err != nil {
		DBLogger.Error("Unable to ping to mongo", "Error", err)
		return client, err
	}

	return client, nil
}

// registers a contract
func (DB *DB) RegisterContract(contract ContractObj) error {
	client, err := DB.DBConnect()
	if err != nil {
		return err
	}
	contractsInstance := client.Database(DBNAME).Collection(ContractCollection)
	insertContext, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := contractsInstance.InsertOne(insertContext, contract)
	if err != nil {
		return err
	}
	DBLogger.Info("Successfully added new contract", "Address", contract.Address, "Name", contract.Name, "ID", res.InsertedID)
	return nil
}

func (DB *DB) GetContracts() (contracts []ContractObj,err error){
	client, err := DB.DBConnect()
	if err != nil {
		log.Fatal(err)
		return contracts,err

	}
	contractsInstance := client.Database(DBNAME).Collection(ContractCollection)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := contractsInstance.Find(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return contracts,err

	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var contract ContractObj
		err := cur.Decode(&contract)
		if err != nil { log.Fatal(err) }
		// do something with result....
		contracts = append(contracts, contract)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return contracts,err
	}
	return contracts,nil
}
func (DB *DB) GetContract() (contract []ContractObj,err error) {
	return contract,nil
}
