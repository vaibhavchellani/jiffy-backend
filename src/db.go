package src

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"

	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"strings"
)

type DB struct{}

// establishes connection with mongodb
func (DB *DB) DBConnect() (client *mongo.Client, err error) {
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
	indexBuilder:=mongo.NewIndexOptionsBuilder()
	indexBuilder.Unique(true).Build()
	res, err := contractsInstance.InsertOne(insertContext, contract)
	if err != nil {
		return err
	}
	DBLogger.Info("Successfully added new contract", "Contract", contract.String(), "ID", res.InsertedID)
	return nil
}

func (DB *DB) GetContracts() (contracts []ContractObj, err error) {
	// connecting to DB
	client, err := DB.DBConnect()
	if err != nil {
		log.Fatal(err)
		return contracts, err
	}

	// creating contract instance
	contractsInstance := client.Database(DBNAME).Collection(ContractCollection)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// searching database
	cur, err := contractsInstance.Find(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return contracts, err
	}

	// iterating over cursor to append to slice of contracts
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var contract ContractObj
		err := cur.Decode(&contract)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
		contracts = append(contracts, contract)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return contracts, err
	}
	return contracts, nil
}

// get a contract by name
func (DB *DB) GetContractViaName(name string) (contract ContractObj, err error) {
	// connect to DB
	client, err := DB.DBConnect()
	if err != nil {
		DBLogger.Error("Unable to connect to DB", "Error", err)
		panic(err)
	}

	DBLogger.Debug("Searching contract", "name", name)

	// creating contract instance
	contractsInstance := client.Database(DBNAME).Collection(ContractCollection)
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	//filtering and getting from DB
	filter := bson.D{{"queryable_name", strings.ToLower(name)}}
	err = contractsInstance.FindOne(ctx, filter).Decode(&contract)
	if err != nil {
		DBLogger.Error("Unable to get contract", "Name", name, "Error", err)
		return contract, err
	}

	DBLogger.Info("Fetched contract", "Contract", contract.String())

	return contract, nil
}

// get a contract by address
func (DB *DB) GetContractViaAddr(addr string) (contract ContractObj, err error) {
	// connect to DB
	client, err := DB.DBConnect()
	if err != nil {
		DBLogger.Error("Unable to connect to DB", "Error", err)
		panic(err)
	}
	DBLogger.Debug("Searching contract", "Address", addr)

	// creating contract instance
	contractsInstance := client.Database(DBNAME).Collection(ContractCollection)
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	//filtering and getting from DB
	filter := bson.D{{"address", addr}}
	err = contractsInstance.FindOne(ctx, filter).Decode(&contract)
	if err != nil {
		DBLogger.Error("Unable to get contract", "Address", addr, "Error", err)
		return contract, err
	}

	DBLogger.Info("Fetched contract", "Contract", contract.String())

	return contract, nil
}
