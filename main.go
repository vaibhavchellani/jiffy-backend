package main

import (
	"log"
	"net/http"

	"github.com/jiffy-backend/helper"

	"github.com/jiffy-backend/server"
)

func main() {
	r := server.NewRouter()
	// TODO fix walker
	//fmt.Println("******Available Routes*******")
	//r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	//	t, err := route.GetPathTemplate()
	//	if err != nil {
	//		fmt.Printf("error getting path template for ROUTE:%v ERROR %v \n" ,route,err)
	//		return err
	//	}
	//	queries,err:=route.GetQueriesTemplates()
	//	if err!=nil{
	//	fmt.Println("err %v",err)
	//		fmt.Printf("error getting query templates for ROUTE:%v ERROR %v \n ",route,err)
	//		//return err
	//	}
	//	methods,err:=route.GetMethods()
	//	if err!=nil{
	//		fmt.Errorf("error getting methods for ROUTE:%v ERROR %v \n",route,err)
	//	}
	//
	//	fmt.Printf("avaiable endpoint: %v methods:%v  queries: %v \n",t,methods,queries)
	//
	//	return nil
	//})
	http.Handle("/", r)
	helper.MainLogger.Info("Http server started successfully ! Listening on port 8000")

	// TODO pick port from config
	log.Fatal(http.ListenAndServe(":8000", r))

}
