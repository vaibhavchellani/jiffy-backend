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
	Queries     []string
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
		[]string{},
	},
	Route{
		"RegisterContract",
		"POST",
		"/update/{entity}",
		controller.Update,
		[]string{},
	},
	Route{
		"GetContracts",
		"GET",
		"/contracts",
		controller.GetContracts,
		[]string{},
	},
	Route{
		"GetContract",
		"GET",
		"/contract",
		controller.GetContract,
		[]string{"filter"},
	},
	// get label by contract address/name
	Route{
		"GetContract",
		"GET",
		"/label",
		controller.GetLabelsByContract,
		[]string{"contract"},
	},
	// get label by creator address
	Route{
		"GetContract",
		"GET",
		"/label",
		controller.GetContract,
		[]string{"creator"},
	},
	// get label by ID
	Route{
		"GetContract",
		"GET",
		"/label",
		controller.GetContract,
		[]string{"id"},
	},
	Route{
		"CheckExistence",
		"GET",
		"/exists",
		controller.CheckExistence,
		[]string{"address", "network"},
	},
	Route{
		"GetDapp",
		"GET",
		"/{dapp_name}",
		controller.GetDapp,
		[]string{},
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
		if len(route.Queries) != 0 {
			for _, query := range route.Queries {
				router.Queries(query, "{"+query+"}")
			}
		}
	}
	return router
}
