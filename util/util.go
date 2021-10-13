package util

import (
	"os"
	"path/filepath"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func WriteToFile(filename, text string) {
	dir := filepath.Dir(filename)
	// use MkdirAll to return nil if dir already exists
	err := os.MkdirAll(dir, 0755)
	HandleErr(err)
	file, err := os.Create(filename)
	HandleErr(err)
	_, err = file.WriteString(text)
	HandleErr(err)
	file.Close()
}

func AppendToFile(filename, text string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	HandleErr(err)
	_, err = file.WriteString(text)
	HandleErr(err)
	file.Close()
}
