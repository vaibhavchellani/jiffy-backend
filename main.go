package main

import (
	"log"
	"net/http"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/jiffy-backend/src"
)

func main() {
	r := src.NewRouter()
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		fmt.Println("route %v", route)
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	})
	http.Handle("/", r)
	fmt.Printf("Http server started successfully ! Listening on port 8000 \n")
	// TODO pick port from config
	log.Fatal(http.ListenAndServe(":8000", r))
}
