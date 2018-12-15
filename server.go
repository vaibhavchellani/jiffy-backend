package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

func main() {
	http.HandleFunc("/abi", addAbi)
	http.ListenAndServe(":8080", nil)
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
	if err != nil {
		fmt.Printf("Unable to connect to mongo , Error : %v", err)
	}m
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Printf("Unable to ping to mongo , Error : %v", err)
	}
}

func addAbi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hi")
}
