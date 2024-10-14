package executor

import (
    "bytes"
    "html/template"
    "os"
    "os/exec"
    "gociflow/logger"
    "gopkg.in/yaml.v2"
)

type Step struct {
    Name    string `yaml:"name"`
    Command string `yaml:"command"`
}

type Workflow struct {
    Name       string `yaml:"name"`
    WorkingDir string `yaml:"working_dir"`
    Steps      []Step `yaml:"steps"`
}

func ExecuteWorkflow(templateName string, variables map[string]string, logsChan chan<- string) error {
    workflowContent, err := loadTemplate(templateName)
    if err != nil {
        return err
    }

    tmpl, err := template.New("workflow").Parse(workflowContent)
    if err != nil {
        return err
    }

    var buffer bytes.Buffer
    err = tmpl.Execute(&buffer, variables)
    if err != nil {
        return err
    }

    var workflow Workflow
    err = yaml.Unmarshal(buffer.Bytes(), &workflow)
    if err != nil {
        return err
    }

    workingDir := workflow.WorkingDir
    if _, err := os.Stat(workingDir); os.IsNotExist(err) {
        if err := os.MkdirAll(workingDir, os.ModePerm); err != nil {
            return err
        }
    }

    for _, step := range workflow.Steps {
        logsChan <- "Executando: " + step.Name + "\n"

        cmdTemplate, err := template.New("command").Parse(step.Command)
        if err != nil {
            return err
        }
        var cmdBuffer bytes.Buffer
        err = cmdTemplate.Execute(&cmdBuffer, variables)
        if err != nil {
            return err
        }
        commandStr := cmdBuffer.String()

        cmd := exec.Command("sh", "-c", commandStr)
        cmd.Dir = workingDir 

        stdoutPipe, err := cmd.StdoutPipe()
        if err != nil {
            return err
        }

        stderrPipe, err := cmd.StderrPipe()
        if err != nil {
            return err
        }

        if err := cmd.Start(); err != nil {
            return err
        }

        go logger.StreamOutput(stdoutPipe, logsChan)
        go logger.StreamOutput(stderrPipe, logsChan)

        if err := cmd.Wait(); err != nil {
            logsChan <- "Erro na etapa: " + step.Name + "\n"
            return err
        }

        logsChan <- "Concluído: " + step.Name + "\n"
    }

    logsChan <- "Workflow concluído com sucesso\n"

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
