package src

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"

	"log"

	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/mongodb/mongo-go-driver/bson"
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
	res, err := contractsInstance.InsertOne(insertContext, contract)
	if err != nil {
		return err
	}
	DBLogger.Info("Successfully added new contract", "Contract",contract.String() , "ID", res.InsertedID)
	return nil
}

func (DB *DB) GetContracts() (contracts []ContractObj, err error) {
	client, err := DB.DBConnect()
	if err != nil {
		log.Fatal(err)
		return contracts, err
	}

	contractsInstance := client.Database(DBNAME).Collection(ContractCollection)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := contractsInstance.Find(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return contracts, err
	}

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
func (DB *DB) GetContract(name string) (contract ContractObj, err error) {
	client, err := DB.DBConnect()
	if err != nil {
		log.Fatal(err)
		return contract, err
	}

	DBLogger.Debug("Searching contract","name",name)
	contractsInstance := client.Database(DBNAME).Collection(ContractCollection)
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	var f bsonx.Doc
	f.Set("name",bsonx.String(name))
	filter := bson.D{{"name", name}}

	err = contractsInstance.FindOne(ctx,filter).Decode(&contract)
	if err!=nil{
		DBLogger.Error("Unable to get contract","Name",name,"Error",err)

	}

	DBLogger.Info("Fetched contract","Contract",contract.String())

	return contract, nil
}
