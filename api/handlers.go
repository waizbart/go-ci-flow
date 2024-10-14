package api

import (
	"encoding/json"
	"gociflow/executor"
	"net/http"
)

type WorkflowRequest struct {
    TemplateName string            `json:"template_name"`
    Variables    map[string]string `json:"variables"`
}

func GetWorkflowsHandler(w http.ResponseWriter, r *http.Request) {
    workflows := []string{"build-and-deploy", "run-tests"} 
    json.NewEncoder(w).Encode(workflows)
}

func ExecuteWorkflowHandler(w http.ResponseWriter, r *http.Request) {
    var request WorkflowRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Solicitação inválida", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming não suportado", http.StatusInternalServerError)
        return
    }

    logsChan := make(chan string)

    go func() {
        if err := executor.ExecuteWorkflow(request.TemplateName, request.Variables, logsChan); err != nil {
            logsChan <- "Erro ao executar o workflow: " + err.Error()
        }
        close(logsChan) 
    }()

    for logLine := range logsChan {
        _, err := w.Write([]byte("data: " + logLine + "\n\n"))
        if err != nil {
            break
        }
        flusher.Flush()
    }
}
