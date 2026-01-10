package shared

import (
	"bytes"
	"log"
	"os"
)

type Log struct {
	Info  *log.Logger
	Error *log.Logger
}

func NewLogger() *Log {
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	return &Log{
		Info:  infoLog,
		Error: errorLog,
	}
}

func NewSpyLogger() *Log {
	infoBuffer := bytes.Buffer{}
	errorBuffer := bytes.Buffer{}

	infoLog := log.New(&infoBuffer, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(&errorBuffer, "[ERROR]\t", log.Ldate|log.Ltime)

	return &Log{
		Info:  infoLog,
		Error: errorLog,
	}
}
