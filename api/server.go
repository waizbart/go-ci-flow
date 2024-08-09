package api

import (
    "github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)

    // Definindo rotas
    router.HandleFunc("/workflows", GetWorkflowsHandler).Methods("GET")
    router.HandleFunc("/workflows", ExecuteWorkflowHandler).Methods("POST")
    
    return router
}
