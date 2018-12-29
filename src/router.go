package src

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes defines the list of routes of our API
type Routes []Route

var controller = &Controller{DB: DB{}}
var routes = Routes{
	Route{
		"RegisterContract",
		"POST",
		"/register/{entity}",
		controller.Register,
	},
	Route{
		"GetContracts",
		"GET",
		"/contracts",
		controller.GetContracts,
	},
	Route{
		"GetContract",
		"GET",
		"/contract/{name}",
		controller.GetContract,
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
	}
	return router
}
