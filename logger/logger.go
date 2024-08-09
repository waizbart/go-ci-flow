package logger

import (
    "io"
    "log"
)

type RealTimeLogger struct {
    writer io.Writer
}

func NewRealTimeLogger(w io.Writer) *RealTimeLogger {
    return &RealTimeLogger{writer: w}
}

func (l *RealTimeLogger) Write(p []byte) (n int, err error) {
    n, err = l.writer.Write(p)
    if err != nil {
        log.Printf("Erro ao escrever log: %v", err)
    }
    return n, err
}
