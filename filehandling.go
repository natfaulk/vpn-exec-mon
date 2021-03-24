package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const logName string = "log.txt"

var handle *os.File

func openLog() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	expath := filepath.Dir(ex)

	filename := filepath.Join(expath, logName)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	handle = f
	return nil
}

func closeLog() error {
	if handle != nil {
		return handle.Close()
	}
	return nil
}

func saveToLog(v string) error {
	timestamp := time.Now().Format(time.RFC3339)
	output := fmt.Sprintf("%s: %s\n", timestamp, v)
	_, err := handle.WriteString(output)
	return err
}
