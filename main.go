package main

import (
	"log"
	"net/http"

	"fmt"
	"github.com/jiffy-backend/src"
)

func main() {
	r := src.NewRouter()
	http.Handle("/", r)
	fmt.Printf("Http server started successfully ! Listening on port 8000 \n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
