package main

import (
    "log"
    "net/http"
    "gociflow/api"
)

func main() {
    // Criação do servidor HTTP
    router := api.NewRouter()

    log.Println("Iniciando o servidor na porta 8080...")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatalf("Erro ao iniciar o servidor: %v", err)
    }
}