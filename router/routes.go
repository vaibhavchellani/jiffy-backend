package router

import "net/http"

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes defines the list of routes of our API
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		controller.Index,
	},
	//Route{
	//	"AddAlbum",
	//	"POST",
	//	"/",
	//	controller.AddAlbum,
	//},
	//Route{
	//	"UpdateAlbum",
	//	"PUT",
	//	"/",
	//	controller.UpdateAlbum,
	//},
	//Route{
	//	"DeleteAlbum",
	//	"DELETE",
	//	"/",
	//	controller.DeleteAlbum,
	//},
}
