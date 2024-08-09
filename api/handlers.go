package api

import (
    "encoding/json"
    "net/http"
    "log"
    "gociflow/executor"
)

// Estrutura para receber variáveis
type WorkflowRequest struct {
    TemplateName string            `json:"template_name"`
    Variables    map[string]string `json:"variables"`
}

// Handler para listar templates de workflows
func GetWorkflowsHandler(w http.ResponseWriter, r *http.Request) {
    workflows := []string{"build-and-deploy", "run-tests"} // Exemplos de templates
    json.NewEncoder(w).Encode(workflows)
}

// Handler para executar um workflow
func ExecuteWorkflowHandler(w http.ResponseWriter, r *http.Request) {
    var request WorkflowRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Solicitação inválida", http.StatusBadRequest)
        return
    }

    // Configurar cabeçalho para streaming de logs
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    // Flushing habilita a escrita dos dados enquanto o workflow executa
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming não suportado", http.StatusInternalServerError)
        return
    }

    // Criar um canal para capturar logs em tempo real
    logsChan := make(chan string)

    // Goroutine para executar o workflow
    go func() {
        defer close(logsChan)
        err := executor.ExecuteWorkflow(request.TemplateName, request.Variables, w)
        if err != nil {
            logsChan <- "Erro ao executar o workflow: " + err.Error()
        }
    }()

    // Enviar logs para o cliente em tempo real
    for logLine := range logsChan {
        _, err := w.Write([]byte("data: " + logLine + "\n\n"))
        if err != nil {
            log.Printf("Erro ao enviar logs: %v", err)
            break
        }
        flusher.Flush()
    }
}
