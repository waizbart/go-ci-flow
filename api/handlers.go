package api

import (
	"encoding/json"
	"fmt"
	"gociflow/executor"
	"net/http"
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

    // Executar o workflow usando o executor
    err := executor.ExecuteWorkflow(request.TemplateName, request.Variables, w)
    if err != nil {
        http.Error(w, "Erro ao executar o workflow", http.StatusInternalServerError)
		fmt.Println("Erro ao executar o workflow:", err)
        return
    }
}
