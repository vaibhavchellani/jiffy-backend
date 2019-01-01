package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jiffy-backend/mongo"
)

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Queries     string
}

// Routes defines the list of routes of our API
type Routes []Route

var controller = &mongo.Controller{DB: mongo.DB{}}
var routes = Routes{
	Route{
		"RegisterContract",
		"POST",
		"/register/{entity}",
		controller.Register,
		"",
	},
	Route{
		"GetContracts",
		"GET",
		"/contracts",
		controller.GetContracts,
		"",
	},
	Route{
		"GetContract",
		"GET",
		"/contract",
		controller.GetContract,
		"filter",
	},
}

//NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
		if route.Queries != "" {
			router.Queries(route.Queries, "{"+route.Queries+"}")
		}
	}
	return router
}
