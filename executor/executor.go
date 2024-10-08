package executor

import (
	"bytes"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Step struct {
    Name    string `yaml:"name"`
    Command string `yaml:"command"`
}

type Workflow struct {
    Name  string `yaml:"name"`
    Steps []Step `yaml:"steps"`
}

func ExecuteWorkflow(templateName string, variables map[string]string, logWriter io.Writer) error {
    // Carregar o template do arquivo
    workflow, err := loadTemplate(templateName)
    if err != nil {
        return err
    }

    // Substituir variáveis no template
    tmpl, err := template.New("workflow").Parse(workflow)
    if err != nil {
        return err
    }

    var buffer bytes.Buffer
    err = tmpl.Execute(&buffer, variables)
    if err != nil {
        return err
    }

    // Converter o resultado para um objeto Workflow
    var parsedWorkflow Workflow
    err = yaml.Unmarshal(buffer.Bytes(), &parsedWorkflow)
    if err != nil {
        return err
    }

    // Executar cada etapa do workflow
    for _, step := range parsedWorkflow.Steps {
        logWriter.Write([]byte("Executando: " + step.Name + "\n"))
        
        // Executar o comando e capturar os logs
        var cmd *exec.Cmd
        if runtime.GOOS == "windows" {
            cmd = exec.Command("cmd.exe", "/c", step.Command)
        } else {
            cmd = exec.Command("sh", "-c", step.Command)
        }
        cmd.Stdout = logWriter
        cmd.Stderr = logWriter

        err = cmd.Run()
        if err != nil {
            logWriter.Write([]byte("Erro na etapa: " + step.Name + "\n"))
            return err
        }

        logWriter.Write([]byte("Concluído: " + step.Name + "\n"))
    }

    return nil
}

func loadTemplate(templateName string) (string, error) {
    filePath := "templates/" + templateName + ".yaml"
    data, err := os.ReadFile(filePath)
    if err != nil {
        return "", err
    }
    return string(data), nil
}