package controllers

import "github.com/trongtb88/urlservice/api/middlewares"

func (server *Server) InitializeRoutes() {
	// Health Check Route
	server.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(server.HealthCheck)).Methods("GET")
	server.Router.HandleFunc("/shorten", middlewares.SetMiddlewareJSON(server.Shorten)).Methods("POST")
	server.Router.HandleFunc("/{shortcode}", middlewares.SetMiddlewareJSON(server.RedirectURL)).Methods("GET")
	server.Router.HandleFunc("/{shortcode}/stats", middlewares.SetMiddlewareJSON(server.Stats)).Methods("GET")
}