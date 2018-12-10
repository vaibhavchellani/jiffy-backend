package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/abi", addAbi)
	http.ListenAndServe(":8080", nil)
}

func addAbi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hi")
}
