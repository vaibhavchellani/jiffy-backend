package src

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

// SERVER the DB server
const SERVER = "localhost:27017"

// DBNAME the name of the DB instance
const DBNAME = "musicstore"

// DOCNAME the name of the document
const DOCNAME = "albums"

type DB struct{}

func (db *DB) RegisterContract(contract ContractObj) Contracts {

	connectionCtx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := mongo.Connect(connectionCtx, "mongodb://localhost:27017")
	if err != nil {
		fmt.Printf("Unable to connect to mongo , Error : %v", err)
	}

	pingCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		fmt.Printf("Unable to ping to mongo , Error : %v", err)
	}

	contractsInstance := client.Database("jiffy").Collection("contracts")
	insertContext, _ := context.WithTimeout(context.Background(), 1*time.Second)

	res, err := contractsInstance.InsertOne(insertContext, contract)
	if err!=nil{
		fmt.Println("ERROR is %v",err)
	}
	id := res.InsertedID
	fmt.Println("Id generated %v", id)
	results := Contracts{}
	//if err := c.Find(nil).All(&results); err != nil {
	//	fmt.Println("Failed to write results:", err)
	//}
	return results
}
