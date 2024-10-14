package logger

import (
    "bufio"
    "io"
)

func StreamOutput(reader io.ReadCloser, logsChan chan<- string) {
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        logsChan <- scanner.Text() + "\n"
    }
}
