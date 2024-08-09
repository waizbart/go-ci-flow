package logger

import (
    "bufio"
    "io"
)

// StreamOutput lê a saída do comando em tempo real e envia para o canal
func StreamOutput(reader io.ReadCloser, logsChan chan string) {
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        logsChan <- scanner.Text() + "\n"
    }
    close(logsChan)
}
