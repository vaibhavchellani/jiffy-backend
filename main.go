package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterContractHandlerFn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hi")
}

func main() {
	// http.HandleFunc("/abi", addAbi)
	// http.ListenAndServe(":8080", nil)
	// ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	// client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
	// if err != nil {
	// 	fmt.Printf("Unable to connect to mongo , Error : %v", err)
	// }
	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	fmt.Printf("Unable to ping to mongo , Error : %v", err)
	// } else {
	// 	fmt.Printf("pinged ! ")
	// }
	r := mux.NewRouter()
	r.HandleFunc("/", RegisterContractHandlerFn).Methods("POST")
	http.Handle("/", r)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
