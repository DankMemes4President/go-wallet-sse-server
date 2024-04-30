package main

import (
	"github.com/gorilla/mux"
	"go-wallet-sse-server/config"
	"go-wallet-sse-server/handlers"
	"net/http"
)

func routes(app *config.Application) http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/sse", handlers.HandleSSE(app))
	return mux
}
