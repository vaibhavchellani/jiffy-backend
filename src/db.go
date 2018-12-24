package src

import (
	"fmt"
	"context"
	"time"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"github.com/mongodb/mongo-go-driver/bson"
)

// SERVER the DB server
const SERVER = "localhost:27017"
// DBNAME the name of the DB instance
const DBNAME = "musicstore"
// DOCNAME the name of the document
const DOCNAME = "albums"

type DB struct {

}

func (db *DB) GetAllContracts() Contracts  {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
	if err != nil {
		fmt.Printf("Unable to connect to mongo , Error : %v", err)
	}

	pingCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		fmt.Printf("Unable to ping to mongo , Error : %v", err)
	} else {
		fmt.Printf("pinged ! ")
	}
	contractsInstance:=client.Database("jiffy").Collection("contracts")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := contractsInstance.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	id := res.InsertedID
	fmt.Println("Id generated %v",id)
	results := Contracts{}
	//if err := c.Find(nil).All(&results); err != nil {
	//	fmt.Println("Failed to write results:", err)
	//}
	return results
}